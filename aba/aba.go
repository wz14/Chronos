package aba

import (
	"acc/config"
	"acc/idchannel"
	"acc/pb"
	"github.com/pkg/errors"
)

var zero = []byte("0")
var one = []byte("1")
var ABAError = "node %d ABA send error in %s"

func ABADecided(value *pb.Message, s config.Start) (*pb.Message, error) {
	pig := s.Getpig()
	pid := pig.GetRootPID(value.Id)

	aba, err := NewABA(pid, s)
	if err != nil {
		return nil, errors.Wrapf(err, ABAError, s.GetConfig().MyID, pid.Id)
	}

	d, err := aba.Decided(value)
	if err != nil {
		return nil, errors.Wrapf(err, ABAError, s.GetConfig().MyID, pid.Id)
	}
	return d, nil
}

func NewABA(pid *idchannel.PrimitiveID, s config.Start) (ABA, error) {
	// may read from config
	version := "mmr"
	if version == "mmr" {
		return NewMMRABA(pid, s), nil
	} else {
		return nil, nil
	}
}

type ABA interface {
	Decided(*pb.Message) (*pb.Message, error)
}
