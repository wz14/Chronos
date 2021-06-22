package core

import (
	"acc/idchannel"
	"acc/pb"
	"context"
	"github.com/pkg/errors"
)

var SendError = "send error %d --> %d in %s" // sender, receiver, pid

func Send(value pb.Message, pid *idchannel.PrimitiveID, nig *idchannel.NodeIDGroup) error {
	senderID, err := nig.GetID(int(value.Sender))
	if err != nil {
		return errors.Wrapf(err, SendError, value.Sender, value.Receiver, pid)
	}
	receiverID, err := nig.GetID(int(value.Receiver))
	if err != nil {
		return errors.Wrapf(err, SendError, value.Sender, value.Receiver, pid)
	}
	s := Sender{
		pid:      pid,
		sender:   senderID,
		receiver: receiverID,
	}
	err = s.Send(value)
	if err != nil {
		return errors.Wrapf(err, SendError, value.Sender, value.Receiver, pid)
	}
	return nil
}

func Receive(pid *idchannel.PrimitiveID) (pb.Message, error) {
	s := Sender{
		pid: pid,
	}
	return s.Receive()
}

type Sender struct {
	pid      *idchannel.PrimitiveID
	sender   *idchannel.NodeID
	receiver *idchannel.NodeID
}

func (s *Sender) Send(value pb.Message) error {
	c := pb.NewNodeConClient(s.receiver.Connect)
	_, err := c.SendMessage(context.Background(), &value)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sender) Receive() (pb.Message, error) {
	m := <-s.pid.C
	return m, nil
}
