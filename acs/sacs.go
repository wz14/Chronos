package acs

import (
	"acc/aba"
	"acc/config"
	"acc/idchannel"
	"acc/logger"
	"acc/pb"
	"acc/rbc"
	"bytes"
	"github.com/golang/protobuf/proto"
	"strconv"
)

func NewSACS(pid *idchannel.PrimitiveID, s config.Start,
	q func(message *pb.Message) bool,
	control chan interface{}) (ACS, error) {
	conf := s.GetConfig()
	b := &SACS{
		rootpid:    pid,
		nig:        s.Getnig(),
		pig:        s.Getpig(),
		l:          logger.NewLoggerWithID(pid.Id, conf.MyID),
		s:          s,
		c:          conf,
		rbcDeliver: make([]chan *pb.Message, conf.N),
		rbc1Notify: make(chan int, conf.N),
		rbc0Notify: make(chan int, conf.N),
		aba1Notify: make(chan int, conf.N),
		aba0Notify: make(chan int, conf.N),
		q:          q,
		control:    control,
	}

	for i := 0; i < conf.N; i++ {
		b.rbcDeliver[i] = make(chan *pb.Message, 1)
	}

	return b, nil
}

type SACS struct {
	rootpid *idchannel.PrimitiveID
	nig     *idchannel.NodeIDGroup
	pig     *idchannel.PIDGroup
	l       *logger.Logger
	s       config.Start
	c       *config.Config

	rbcDeliver []chan *pb.Message // [N] len = 1

	rbc1Notify chan int // len = 0
	rbc0Notify chan int // len = 0
	aba1Notify chan int // len = 0
	aba0Notify chan int // len = 0

	//Predicate Q(v) -> 0/1
	q       func(message *pb.Message) bool
	control chan interface{}
}

// return a set of pb.Message
func (b *SACS) Decided(message *pb.Message) ([]*pb.Message, error) {
	go b.rbcb(message)
	go b.rbcd()
	go b.aba()
	l := b.rollback()
	b.l.Infof("rollback index set: %v", l)
	output := []*pb.Message{}
	for index, d := range l {
		if d == 1 {
			v := <-b.rbcDeliver[index]
			b.l.Infof("get message from %d RBC", index)
			// unmarshal v
			d := pb.Message{}
			err := proto.Unmarshal(v.Data, &d)
			if err != nil {
				b.l.Errorf("unmarshal fail")
			}
			output = append(output, &d)
		}
	}
	b.l.Infof("output %d things", len(output))
	return output, nil
}

func (b *SACS) rbcb(message *pb.Message) {
	byt, err := proto.Marshal(message)
	if err != nil {
		b.l.Errorf("marshal value fail in %s", b.rootpid.Id)
	}
	rbcpid := b.pig.GetChildPID("RBC"+strconv.Itoa(b.c.MyID), b.rootpid)
	b.l.Infof("RBC broadcast")
	rbc.RBCBroadcast(rbc.BRACHA, &pb.Message{
		Id:       rbcpid.Id,
		Sender:   uint32(b.c.MyID),
		Receiver: 0,
		Data:     byt,
	}, b.s)
}

func (b *SACS) rbcd() {
	for i := 0; i < b.c.N; i++ {
		go func(j int) {
			rbcpid := b.pig.GetChildPID("RBC"+strconv.Itoa(j), b.rootpid)
			mes, err := rbc.RBCDeliver(rbc.BRACHA, rbcpid.Id, b.s)
			if err != nil {
				b.l.Errorf("receive from RBC fail: %s", err.Error())
			}
			b.l.Infof("RBC receive from %d", j)
			tmp := &pb.Message{}
			err = proto.Unmarshal(mes.Data, tmp)
			if err != nil {
				b.l.Errorf("unmarshal fail in tmp")
			}
			b.l.Infof("tmpmessage: %s", string(tmp.Data))
			if b.q(tmp) {
				b.l.Infof("satisfy sacs predicate")
				b.rbc1Notify <- j
				b.rbcDeliver[j] <- mes
			} else {
				<-b.control
				b.rbc1Notify <- j
				b.rbcDeliver[j] <- mes
				b.l.Infof("receive sacs signal")
			}
		}(i)
	}
}

func (b *SACS) aba() {
	already := []bool{}
	for w := 0; w < b.c.N; w++ {
		already = append(already, false)
	}
	for {
		var est []byte
		var i int
		select {
		case i = <-b.rbc1Notify:
			est = aba.One
		case i = <-b.rbc0Notify:
			est = aba.Zero
		}
		if already[i] {
			continue
		} else {
			already[i] = true
		}
		b.l.Infof("open %d aba with %s", i, string(est))
		go func() {
			// mark i already output
			abapid := b.pig.GetChildPID("ABA"+strconv.Itoa(i), b.rootpid)
			output, err := aba.ABADecided(&pb.Message{
				Id:       abapid.Id,
				Sender:   uint32(b.c.MyID),
				Receiver: 0,
				Data:     est,
			}, b.s)
			if err != nil {
				b.l.Errorf("decide fail from aba: %s", err.Error())
			}
			b.l.Infof("decide %s from %d aba", string(output.Data), i)
			if bytes.Equal(output.Data, aba.One) {
				b.aba1Notify <- i
			} else if bytes.Equal(output.Data, aba.Zero) {
				b.aba0Notify <- i
			} else {
				b.l.Errorf("decide weired value from aba: %s", string(output.Data))
			}
		}()
	}
}

// blocking
func (b *SACS) rollback() []int {
	// 0 init
	// 1 receive 1 from aba
	// 2 receive 0 from aba
	l := []int{}
	for i := 0; i < b.c.N; i++ {
		l = append(l, 0)
	}
	aba1count := 0
	abacount := 0
	for {
		select {
		case k := <-b.aba0Notify:
			l[k] = 2
		case k := <-b.aba1Notify:
			l[k] = 1
			aba1count += 1
		}
		abacount += 1
		if aba1count >= b.c.F*2+1 {
			b.l.Infof("collect %d decide 1 aba, send 0 to others aba", aba1count)
			for i := 0; i < b.c.N; i++ {
				if l[i] == 0 {
					b.rbc0Notify <- i
				}
			}
		}

		if abacount >= b.c.N {
			b.l.Info("collect all aba")
			break
		}
	}
	return l
}
