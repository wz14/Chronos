package consensus

import (
	"acc/config"
	"acc/idchannel"
	"acc/pb"
	"github.com/pkg/errors"
)

func Consensus(value *pb.Message, s config.Start) ([]*pb.Message, error) {
	pig := s.Getpig()
	pid := pig.GetRootPID(value.Id)

	acs, err := NewConsen(pid, s)
	if err != nil {
		return nil, errors.Errorf("create Consensus fail: %s", err.Error())
	}
	return acs.Propose(value)
}

func NewConsen(pid *idchannel.PrimitiveID, s config.Start) (Consen, error) {
	version := "HB"
	if version == "HB" {
		return NewHB(pid, s)
	} else {
		return nil, errors.New("no such consensus type")
	}
}

type Consen interface {
	Propose(message *pb.Message) ([]*pb.Message, error)
}
