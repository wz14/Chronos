package consensus

import (
	"acc/acs"
	"acc/config"
	"acc/crypto"
	"acc/idchannel"
	"acc/logger"
	"acc/pb"
	"acc/rbc"
	"bytes"
	"github.com/golang/protobuf/proto"
	"strconv"
	"time"
)

func NewCHOS(pid *idchannel.PrimitiveID, s config.Start) (Consen, error) {
	conf := s.GetConfig()
	c := &CHOS{
		rootpid: pid,
		e:       s.GetEConfig(),
		pig:     s.Getpig(),
		c:       conf,
		s:       s,
		l:       logger.NewLoggerWithID("CHOS", conf.MyID),

		rbcDeliver:   make([]chan *pb.Message, conf.N),
		finalDeliver: make([]chan *pb.Message, conf.N),
		tsDeliver:    make([]chan *pb.TS, conf.N),
		controls:     make([]chan interface{}, conf.N),

		rbc1Notify: make(chan int, conf.N),
		rbc0Notify: make(chan int, conf.N),
		acs1Notify: make(chan int, conf.N),
		acs0Notify: make(chan int, conf.N),
	}

	for i := 0; i < conf.N; i++ {
		c.rbcDeliver[i] = make(chan *pb.Message, 1)
		c.finalDeliver[i] = make(chan *pb.Message, 1)
		c.tsDeliver[i] = make(chan *pb.TS, 1)
		c.controls[i] = make(chan interface{}, 0)
	}

	return c, nil
}

type CHOS struct {
	rootpid *idchannel.PrimitiveID
	e       *crypto.TPKE
	pig     *idchannel.PIDGroup
	c       *config.Config
	s       config.Start
	l       *logger.Logger

	rbcDeliver   []chan *pb.Message // [N] len = 1
	finalDeliver []chan *pb.Message
	tsDeliver    []chan *pb.TS

	controls []chan interface{}

	rbc1Notify chan int // len = 0
	rbc0Notify chan int // len = 0
	acs1Notify chan int // len = 0
	acs0Notify chan int // len = 0
}

func (h *CHOS) Propose(value *pb.Message) ([]*pb.Message, []*pb.TS, error) {

	go h.rbcb(value)
	go h.rbcd()
	go h.sacs()

	l := h.rollback()
	h.l.Infof("rollback index set: %v", l)
	output := []*pb.Message{}
	ts := []*pb.TS{}
	for index, d := range l {
		if d == 1 {
			v := <-h.finalDeliver[index]
			h.l.Infof("get message from %d RBC", index)
			// unmarshal v
			d := pb.Message{}
			err := proto.Unmarshal(v.Data, &d)
			if err != nil {
				h.l.Errorf("unmarshal fail for %d", index)
			}
			output = append(output, &d)
			t := <-h.tsDeliver[index]
			ts = append(ts, t)
		}
	}
	h.l.Infof("output %d things", len(output))
	return output, ts, nil
}

func (h *CHOS) rbcb(message *pb.Message) {
	// marshal value
	byt, err := proto.Marshal(message)
	if err != nil {
		h.l.Error("marshal fail")
	}

	cpid := h.pig.GetChildPID("AVIDRBC"+strconv.Itoa(h.c.MyID), h.rootpid)

	rbc.RBCBroadcast(rbc.AVID, &pb.Message{
		Id:       cpid.Id,
		Sender:   uint32(h.c.MyID),
		Receiver: 0,
		Data:     byt,
	}, h.s)
}

func (h *CHOS) rbcd() {
	for i := 0; i < h.c.N; i++ {
		go func(j int) {
			rbcpid := h.pig.GetChildPID("AVIDRBC"+strconv.Itoa(j), h.rootpid)
			mes, err := rbc.RBCDeliver(rbc.AVID, rbcpid.Id, h.s)
			if err != nil {
				h.l.Errorf("receive from AVID RBC fail: %s", err.Error())
			}
			h.l.Infof("AVID RBC receive from %d", j)
			h.rbc1Notify <- j
			h.rbcDeliver[j] <- mes
			h.finalDeliver[j] <- mes
		}(i)
	}
}

func (h *CHOS) sacs() {
	already := []bool{}
	for w := 0; w < h.c.N; w++ {
		already = append(already, false)
	}
	for {
		var i int
		var mes []byte
		var k string
		var q func(message *pb.Message) bool
		select {
		case i = <-h.rbc1Notify:
			m := <-h.rbcDeliver[i]
			mes = h.buildTS(m)
			q = h.buildPredicate(mes)
			k = "timestamp"
		case i = <-h.rbc0Notify:
			mes = h.buildEmptyTS()
			q = h.buildTruePredicate()
			k = "dummy"
		}
		if already[i] {
			continue
		} else {
			already[i] = true
		}
		h.l.Info("open %d sacs with %s", i, k)
		go func() {
			acspid := h.pig.GetChildPID("SACS"+strconv.Itoa(i), h.rootpid)
			output, err := acs.ACSDecided(acs.S, &pb.Message{
				Id:       acspid.Id,
				Sender:   uint32(h.c.MyID),
				Receiver: 0,
				Data:     mes,
			}, h.s, q, h.controls[i]) // build predicates and controls
			if err != nil {
				h.l.Errorf("decided fail from acs: %s", err.Error())
			}
			midts, ok := h.checkOutput(output)
			if ok {
				h.acs1Notify <- i
				h.tsDeliver[i] <- midts
				// send output to main routine
			} else {
				h.acs0Notify <- i
			}
		}()
	}
}

func (h *CHOS) buildPredicate(byt []byte) func(message *pb.Message) bool {
	// if m is empty message,
	// striq always return false || loseq always return true but unmarshal fail.
	// else
	// striq if hash()==tx return true || loseq always return true but unmarshal fail

	// 1. unmarshal ts
	myts := pb.TS{}
	err := proto.Unmarshal(byt, &myts)
	if err != nil {
		h.l.Errorf("unmarshal ts fail: %s", err.Error())
	}

	q := func(message *pb.Message) bool {
		byt := message.Data
		ts := pb.TS{}
		err := proto.Unmarshal(byt, &ts)
		if err != nil {
			h.l.Errorf("unmarshal ts fail: %s", err.Error())
		}

		if ts.Dummy {
			return false
		}

		if !bytes.Equal(ts.Hash, myts.Hash) {
			return false
		}

		if ts.Sender != myts.Sender {
			return false
		}

		if ts.Num != ts.Num {
			return false
		}

		if len(ts.TS) != len(ts.TS) {
			return false
		}

		return true
	}
	return q
}

func (h *CHOS) buildTruePredicate() func(message *pb.Message) bool {
	q := func(message *pb.Message) bool {
		return true
	}
	return q
}

/*
TS{
	dummy bool
	m hash
	sender uint32
	txnum int64
	TSs []int54
}
*/
func (h *CHOS) buildTS(m *pb.Message) []byte {
	// unmarshal m.Data to realm
	// for all tx in real.m, assign timestamp
	// construct TS and marshal TS to byt bytes.
	// return a pb.Message
	inputv := &pb.Message{}
	err := proto.Unmarshal(m.Data, inputv)
	if err != nil {
		h.l.Errorf("build ts fail: %s", err.Error())
	}
	l := len(inputv.Data) / 250
	// length of one tx is 250
	times := make([]uint64, l)
	for i := 0; i < l; i++ {
		times[i] = uint64(time.Now().UnixNano())
	}
	ts := &pb.TS{
		Dummy:  false,
		Hash:   crypto.Hash(inputv.Data),
		Sender: inputv.Sender,
		Num:    uint64(l),
		TS:     times,
	}
	byt, err := proto.Marshal(ts)
	if err != nil {
		h.l.Errorf("TS marshal fail: %s", err.Error())
	}
	return byt
}

func (h *CHOS) buildEmptyTS() []byte {
	ts := pb.TS{
		Dummy:  true,
		Hash:   nil,
		Sender: 0,
		Num:    0,
		TS:     nil,
	}
	byt, err := proto.Marshal(&ts)
	if err != nil {
		h.l.Errorf("marshal empty ts fail: %s", err.Error())
	}
	return byt
}

func (h *CHOS) checkOutput(m []*pb.Message) (*pb.TS, bool) {
	// 1. hash must be same
	// 2. length must be same
	// 3. sender must be same
	// 4. no dummy message
	// 5. generate middle timestamp
	p := h.buildPredicate(m[0].Data)

	for i := 1; i < len(m); i++ {
		if !p(m[i]) {
			return nil, false
		}
	}

	tss := make([]*pb.TS, 2*h.c.F+1)

	for i := 0; i < 2*h.c.F+1; i++ {
		ts := &pb.TS{}
		err := proto.Unmarshal(m[i].Data, ts)
		if err != nil {
			h.l.Errorf("checkouput fail: %s", err.Error())
		}
		tss[i] = ts
	}

	var midts []uint64 = make([]uint64, tss[0].Num)
	for i := 0; i < int(tss[0].Num); i++ {
		for j := 0; j < 2*h.c.F+1; j++ {
			midts[i] += tss[j].TS[i] / uint64(2*h.c.F+1)
		}
	}

	result := pb.TS{
		Dummy:  false,
		Hash:   tss[0].Hash,
		Sender: tss[0].Sender,
		Num:    tss[0].Num,
		TS:     midts,
	}

	return &result, true
}

func (h *CHOS) rollback() []int {
	// 0 init
	// 1 receive 1 from acs
	// 2 receive 0 from acs
	l := []int{}
	for i := 0; i < h.c.N; i++ {
		l = append(l, 0)
	}
	acs1count := 0
	acscount := 0
	for {
		select {
		case k := <-h.acs1Notify:
			l[k] = 2
		case k := <-h.acs0Notify:
			l[k] = 1
			acs1count += 1
		}
		acscount += 1
		if acs1count >= h.c.F*2+1 {
			h.l.Infof("collect %d decide 1 aba, send 0 to others aba", acs1count)
			for i := 0; i < h.c.N; i++ {
				if l[i] == 0 {
					h.rbc0Notify <- i
					// control send signal
					go func() {
						for {
							h.controls[i] <- 0
						}
					}()
				}
			}
		}

		if acscount >= h.c.N {
			h.l.Info("collect all aba")
			break
		}
	}
	return l
}
