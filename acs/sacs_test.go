package acs

import (
	"acc/config"
	"acc/pb"
	"bytes"
	"strconv"
	"testing"
)

// TestNewSACS :all inputs of nodes satisfy strict predicate && no control message be sent
// expected :all nodes return output satisfy strict  predicate
func TestNewSACS(t *testing.T) {
	nls := config.NewLocalStart(func(s config.Start) {
		q := func(message *pb.Message) bool {
			return bytes.Contains(message.Data, []byte("mockdata"))
		}
		control := make(chan interface{}, 0)
		conf := s.GetConfig()
		ml, err := ACSDecided(S, &pb.Message{
			Id:       "SACS1",
			Sender:   uint32(conf.MyID),
			Receiver: 0,
			Data:     []byte("mockdata" + strconv.Itoa(conf.MyID)),
		}, s, q, control)
		if err != nil {
			t.Error("decide fail")
		}
		for i, m := range ml {
			if !bytes.Equal(m.Data, []byte("mockdata"+
				strconv.FormatUint(uint64(m.Sender), 10))) {
				t.Error("acs decided bad value")
			}
			t.Logf("%d values: %s", i, string(m.Data))
		}
		if len(ml) < 2*conf.F+1 {
			t.Error("too few values in common set")
		}
	}, "./mock/config1.yaml")
	nls.Run()
}

// TestNewSACS1 is a error test for finish property
// if all inputs of honest nodes satisfy loose q ,
// and all honest nodes receive control signal. Then,
// all honest node will output eventually (may not satisfy strict q).
func TestNewSACS1(t *testing.T) {
	nls := config.NewLocalStart(func(s config.Start) {
		q := func(message *pb.Message) bool {
			return bytes.Contains(message.Data, []byte("strictdata"))
		}
		control := make(chan interface{}, 0)

		conf := s.GetConfig()
		var ml []*pb.Message
		var err error
		go func() {
			for {
				control <- 0
			}
		}()
		ml, err = ACSDecided(S, &pb.Message{
			Id:       "SACS2",
			Sender:   uint32(conf.MyID),
			Receiver: 0,
			Data:     []byte("loosedata" + strconv.Itoa(conf.MyID)),
		}, s, q, control)

		if err != nil {
			t.Error("decide fail")
		}
		for i, m := range ml {
			if !bytes.Equal(m.Data, []byte("loosedata"+
				strconv.FormatUint(uint64(m.Sender), 10))) {
				t.Error("acs decided bad value")
			}
			t.Logf("%d values: %s", i, string(m.Data))
		}
		if len(ml) < 2*conf.F+1 {
			t.Error("too few values in common set")
		}
	}, "./mock/config1.yaml")
	nls.Run()
}
