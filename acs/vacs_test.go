package acs

import (
	"acc/config"
	"acc/pb"
	"bytes"
	"strconv"
	"testing"
)

// plain test
func TestNewVACS(t *testing.T) {
	nls := config.NewLocalStartWithReadLocalConfig(func(s config.Start) {
		q := func(message *pb.Message) bool {
			return bytes.Contains(message.Data, []byte("mockdata"))
		}
		conf := s.GetConfig()
		ml, err := ACSDecided(V, &pb.Message{
			Id:       "ACS1",
			Sender:   uint32(conf.MyID),
			Receiver: 0,
			Data:     []byte("mockdata" + strconv.Itoa(conf.MyID)),
		}, s, q, nil)
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

// TestNewVACS1 is a error test
func TestNewVACS1(t *testing.T) {
	nls := config.NewLocalStartWithReadLocalConfig(func(s config.Start) {
		q := func(message *pb.Message) bool {
			return bytes.Contains(message.Data, []byte("gooddata"))
		}
		conf := s.GetConfig()
		var ml []*pb.Message
		var err error
		if conf.MyID == 0 {
			ml, err = ACSDecided(V, &pb.Message{
				Id:       "ACS1",
				Sender:   uint32(conf.MyID),
				Receiver: 0,
				Data:     []byte("baddata" + strconv.Itoa(conf.MyID)),
			}, s, q, nil)
		} else {
			ml, err = ACSDecided(V, &pb.Message{
				Id:       "ACS1",
				Sender:   uint32(conf.MyID),
				Receiver: 0,
				Data:     []byte("gooddata" + strconv.Itoa(conf.MyID)),
			}, s, q, nil)
		}

		if err != nil {
			t.Error("decide fail")
		}
		for i, m := range ml {
			if !bytes.Equal(m.Data, []byte("gooddata"+
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
