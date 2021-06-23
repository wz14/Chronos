package RBC

import (
	"acc/config"
	"acc/core"
	"acc/crypto"
	"acc/idchannel"
	"acc/logger"
	"acc/pb"
	proto "github.com/golang/protobuf/proto"
	"sync"
)

func NewBrachaRBC(sender *idchannel.NodeID, pid *idchannel.PrimitiveID,
	s config.Start) *BrachaRBC {
	nig := s.Getnig()
	pig := s.Getpig()
	id := s.GetConfig()
	return &BrachaRBC{
		sender:      sender,
		pid:         pid,
		nig:         nig,
		pig:         pig,
		s:           s,
		c:           id,
		l:           logger.NewLoggerWithID("brachaRBC", id.MyID),
		v:           make(chan []byte, 1),
		isReady:     false,
		isReadyLock: sync.Mutex{},
	}
}

type BrachaRBC struct {
	c      *config.Config
	sender *idchannel.NodeID
	pid    *idchannel.PrimitiveID
	nig    *idchannel.NodeIDGroup
	pig    *idchannel.PIDGroup
	s      config.Start
	l      *logger.Logger
	// wirte once, no need to lock
	v           chan []byte
	isReady     bool
	isReadyLock sync.Mutex
}

/*

Bracha's Reliable Broadcast

Upon receive broadcast v:
	sender s broadcast <init, v>

Upon receive <init, v>:
	node i broadcast <echo, i, v>

Upon receive 2f+1 <echo, j, v>:
	node i broadcast <ready, i, hash(v)>

Upon receive f+1 <ready, j, hash(v)>:
	node i broadcast <ready, i, hash(v)>

Upon receive 2f+1 <ready, j, hash(v)>:
	node i deliver v

*/

// non-blocking in broadcast
func (b *BrachaRBC) Broadcast(value *pb.Message) error {
	go b.init(value)
	return nil
}

// blocking in deliver
func (b *BrachaRBC) Deliver() (*pb.Message, error) {
	ch := make(chan *pb.Message, 0)
	go b.echo()
	go b.ready()
	go b.final(ch)
	v := <-ch
	return v, nil
}

func (b *BrachaRBC) init(value *pb.Message) {
	// build <init, v> message
	childp := b.pig.GetChildPID("Init", b.pid)
	byt, err := proto.Marshal(value)
	if err != nil {
		b.l.Errorf("marshal value fail in %s", childp.Id)
	}
	b.l.Infof("broadcast init message")
	core.BroadCast(&pb.Message{
		Id:       childp.Id,
		Sender:   uint32(b.c.MyID),
		Receiver: 0,
		Data:     byt,
	}, b.s)
}

func (b *BrachaRBC) echo() {
	m, err := core.Receive(b.pig.GetChildPID("Init", b.pid).Id, b.s)
	if err != nil {
		b.l.Errorf("init message receive fail")
		return
	}
	b.l.Infof("receive init message from %d", m.Sender)
	b.l.Infof("broadcast echo message")
	childp := b.pig.GetChildPID("Echo", b.pid)
	core.BroadCast(&pb.Message{
		Id:       childp.Id,
		Sender:   uint32(b.c.MyID),
		Receiver: 0,
		Data:     m.Data,
	}, b.s)
}

func (b *BrachaRBC) ready() {
	EchoCount := 0
	var mes []byte
	for {
		m, err := core.Receive(b.pig.GetChildPID("Echo", b.pid).Id, b.s)
		if err != nil {
			b.l.Errorf("echo message receive fail")
			continue
		}
		b.l.Infof("receive echo message from %d", m.Sender)
		// TODO:check m is valid and first come and same byt
		EchoCount += 1
		if EchoCount >= 2*b.c.F+1 {
			b.l.Infof("collect 2f+1 echos")
			mes = m.Data
			break
		}
	}
	b.l.Infof("is it run here?")
	b.v <- mes
	b.isReadyLock.Lock()
	if b.isReady {
		b.isReadyLock.Unlock()
		b.l.Infof("already broadcast ready, stop ready routine")
		return
	}
	b.isReady = true
	b.isReadyLock.Unlock()
	b.l.Infof("broadcast ready message")
	childp := b.pig.GetChildPID("Ready", b.pid)
	core.BroadCast(&pb.Message{
		Id:       childp.Id,
		Sender:   uint32(b.c.MyID),
		Receiver: 0,
		Data:     crypto.Hash(mes), /* marshal(m.Data)*/
	}, b.s)
}

func (b *BrachaRBC) final(ch chan *pb.Message) {
	//if deliver, then send mesage to ch
	var hash []byte
	ReadyCount := 0
	for {
		m, err := core.Receive(b.pig.GetChildPID("Ready", b.pid).Id, b.s)
		if err != nil {
			b.l.Errorf("ready message receive fail")
			continue
		}
		b.l.Infof("receive ready message from %d", m.Sender)
		// TODO: check m's data
		ReadyCount += 1
		if ReadyCount >= b.c.F+1 {
			b.isReadyLock.Lock()
			if !b.isReady {
				b.isReady = true
				b.isReadyLock.Unlock()
				b.l.Infof("receive f+1 ready message and not yet broadcast ready")
				childp := b.pig.GetChildPID("Ready", b.pid)
				b.l.Infof("broadcast ready message")
				core.BroadCast(&pb.Message{
					Id:       childp.Id,
					Sender:   uint32(b.c.MyID),
					Receiver: 0,
					Data:     m.Data, /* marshal(m.Data)*/
				}, b.s)
			} else {
				b.isReadyLock.Unlock()
			}
		}

		if ReadyCount >= b.c.F*2+1 {
			b.l.Infof("receive 2f+1 ready message")
			hash = m.Data
			message := <-b.v
			if crypto.HashVerify(message, hash) {
				mes := pb.Message{}
				err := proto.Unmarshal(message, &mes)
				if err != nil {
					b.l.Errorf("unmarshal message fail in %s", b.pid.Id)
				}
				b.l.Infof("return message to deliver channel")
				ch <- &mes
				return
			} else {
				b.l.Errorf("corrupted seems bigger than N/3 in %s", b.pid.Id)
				return
			}
		}
	}

}
