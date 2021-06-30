package acs

import (
	"acc/config"
	"acc/idchannel"
	"acc/pb"
	"github.com/pkg/errors"
)

type ACSTYPE int

const (
	BENOR ACSTYPE = 1
	V     ACSTYPE = 2
)

// ACSDecided :if typ is BENOR, q is omit
// else q is used to judge if transaction is valid
func ACSDecided(typ ACSTYPE, value *pb.Message, s config.Start,
	q func(message *pb.Message) bool) ([]*pb.Message, error) {
	pig := s.Getpig()
	pid := pig.GetRootPID(value.Id)

	acs, err := NewACS(typ, pid, s, q)
	if err != nil {
		return nil, errors.Wrap(err, "create ACS fail")
	}

	d, err := acs.Decided(value)
	if err != nil {
		return nil, errors.Wrap(err, "decide ACS fail")
	}

	return d, nil
}

func NewACS(typ ACSTYPE, pid *idchannel.PrimitiveID, s config.Start,
	q func(message *pb.Message) bool) (ACS, error) {
	if typ == BENOR {
		return NewBenorACS(pid, s)
	} else if typ == V {
		return NewVACS(pid, s, q)
	} else {
		return nil, errors.New("no such version acs protocol")
	}
}

type ACS interface {
	Decided(message *pb.Message) ([]*pb.Message, error)
}
