package rbc

import (
	"acc/config"
	"acc/idchannel"
	"acc/pb"
	"github.com/pkg/errors"
)

var RBCSError = "node %d RBC send error in %s"
var RBCDError = "node deliver error in %s"

func RBCBroadcast(value *pb.Message, s config.Start) error {
	nig := s.Getnig()
	pig := s.Getpig()
	pid := pig.GetRootPID(value.Id)
	senderID, err := nig.GetID(int(value.Sender))
	if err != nil {
		return errors.Wrapf(err, RBCSError, value.Sender, pid.Id)
	}

	rbc, err := NewRBC(senderID, pid, s)
	if err != nil {
		return errors.Wrapf(err, RBCSError, value.Sender, pid.Id)
	}

	err = rbc.Broadcast(value)
	if err != nil {
		return errors.Wrapf(err, RBCSError, value.Sender, pid.Id)
	}

	return nil
}

func RBCDeliver(id string, s config.Start) (*pb.Message, error) {
	pig := s.Getpig()
	pid := pig.GetRootPID(id)
	rbc, err := NewRBC(nil, pid, s)
	if err != nil {
		return nil, errors.Wrapf(err, RBCDError, pid.Id)
	}

	m, err := rbc.Deliver()
	if err != nil {
		return nil, errors.Wrapf(err, RBCDError, pid.Id)
	}

	return m, nil
}

func NewRBC(sender *idchannel.NodeID, pid *idchannel.PrimitiveID,
	s config.Start) (RBC, error) {
	// may read from config
	version := "bracha"
	if version == "bracha" {
		return NewBrachaRBC(sender, pid, s), nil
	} else {
		return nil, nil
	}
}

type RBC interface {
	Broadcast(*pb.Message) error
	Deliver() (*pb.Message, error)
}
