package core

import (
	"acc/idchannel"
	"acc/pb"
	"context"
	"github.com/pkg/errors"
)

var BroadcastError = "node %d send error in %s" // sender, receiver, pid

func BroadCast(value pb.Message, pid *idchannel.PrimitiveID, nig *idchannel.NodeIDGroup) error {
	senderID, err := nig.GetID(int(value.Sender))
	if err != nil {
		return errors.Wrapf(err, BroadcastError, value.Sender, pid)
	}

	bc := BroadCaster{
		pid:    pid,
		sender: senderID,
		nig:    nig,
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
	nig    *idchannel.NodeIDGroup
}

func (s *BroadCaster) Broadcast(value pb.Message) error {
	for i := 0; i < s.nig.N; i++ {
		id, err := s.nig.GetID(i)
		if err != nil {
			return errors.Wrapf(err, "get id %d fail", i)
		}
		c := pb.NewNodeConClient(id.Connect)
		value.Receiver = uint32(i)
		_, err = c.SendMessage(context.Background(), &value)
		if err != nil {
			return errors.Wrapf(err, "send message fail to %d", i)
		}
	}
	return nil
}
