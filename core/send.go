package core

import (
	"acc/config"
	"acc/idchannel"
	"acc/pb"
	"context"
	"github.com/pkg/errors"
)

var SendError = "send error %d --> %d in %s" // sender, receiver, pid

func Send(value *pb.Message, s config.Start) error {
	nig := s.Getnig()
	pig := s.Getpig()
	pid := pig.GetRootPID(value.Id)
	senderID, err := nig.GetID(int(value.Sender))
	if err != nil {
		return errors.Wrapf(err, SendError, value.Sender, value.Receiver, pid)
	}
	receiverID, err := nig.GetID(int(value.Receiver))
	if err != nil {
		return errors.Wrapf(err, SendError, value.Sender, value.Receiver, pid)
	}
	sender := Sender{
		pid:      pid,
		sender:   senderID,
		receiver: receiverID,
		s:        s,
	}
	err = sender.Send(value)
	if err != nil {
		return errors.Wrapf(err, SendError, value.Sender, value.Receiver, pid)
	}
	return nil
}

func Receive(pid string, s config.Start) (*pb.Message, error) {
	p := s.Getpig().GetRootPID(pid)
	sender := Sender{
		pid: p,
	}
	return sender.Receive()
}

type Sender struct {
	pid      *idchannel.PrimitiveID
	sender   *idchannel.NodeID
	receiver *idchannel.NodeID
	s        config.Start
}

func (s *Sender) Send(value *pb.Message) error {
	c := pb.NewNodeConClient(s.receiver.Connect)
	_, err := c.SendMessage(context.Background(), value)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sender) Receive() (*pb.Message, error) {
	m := <-s.pid.C
	return m, nil
}
