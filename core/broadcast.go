package core

import (
	"acc/config"
	"acc/idchannel"
	"acc/pb"
	"context"
	"github.com/pkg/errors"
)

var BroadcastError = "node %d send error in %s" // sender, receiver, pid

func BroadCast(value *pb.Message, s config.Start) error {
	pig := s.Getpig()
	pid := pig.GetRootPID(value.Id)

	nig := s.Getnig()
	senderID, err := nig.GetID(int(value.Sender))
	if err != nil {
		return errors.Wrapf(err, BroadcastError, value.Sender, pid)
	}
	bc := BroadCaster{
		pid:    pid,
		sender: senderID,
		s:      s,
	}

	err = bc.Broadcast(value)
	if err != nil {
		return errors.Wrapf(err, BroadcastError, value.Sender, pid)
	}

	return nil
}

type BroadCaster struct {
	pid    *idchannel.PrimitiveID
	sender *idchannel.NodeID
	s      config.Start
}

func (b *BroadCaster) Broadcast(value *pb.Message) error {
	conf := b.s.GetConfig()
	nig := b.s.Getnig()
	for i := 0; i < conf.N; i++ {
		id, err := nig.GetID(i)
		if err != nil {
			return errors.Wrapf(err, "get id %d fail", i)
		}
		i := i
		go func() {
			c := pb.NewNodeConClient(id.Connect)
			newValue := pb.Message{
				Id:       value.Id,
				Sender:   value.Sender,
				Receiver: uint32(i),
				Data:     value.Data,
			}
			_, _ = c.SendMessage(context.Background(), &newValue)
		}()
	}
	return nil
}
