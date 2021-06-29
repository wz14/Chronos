package acs

import (
	"acc/config"
	"acc/idchannel"
	"acc/pb"
	"github.com/pkg/errors"
)

func ACSDecided(value *pb.Message, s config.Start) ([]*pb.Message, error) {
	pig := s.Getpig()
	pid := pig.GetRootPID(value.Id)

	acs, err := NewACS(pid, s)
	if err != nil {
		return nil, errors.Wrap(err, "create ACS fail")
	}

	d, err := acs.Decided(value)
	if err != nil {
		return nil, errors.Wrap(err, "decide ACS fail")
	}

	return d, nil
}

func NewACS(pid *idchannel.PrimitiveID, s config.Start) (ACS, error) {
	// version := "benor"
	version := "benor"
	if version == "benor" {
		return NewBenorACS(pid, s)
	} else {
		return nil, errors.New("no such version acs protocol")
	}
	return nil, nil
}

type ACS interface {
	Decided(message *pb.Message) ([]*pb.Message, error)
}
