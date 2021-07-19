package aba

import (
	"acc/config"
	"acc/pb"
	"bytes"
	"testing"
)

// Test the condition that all node propose 1
func TestNewMMRABA(t *testing.T) {
	//mmr := NewMMRABA()
	nls := config.NewLocalStartWithReadLocalConfig(func(s config.Start) {
		conf := s.GetConfig()
		m, err := ABADecided(&pb.Message{
			Id:       "testMMR1",
			Sender:   uint32(conf.MyID),
			Receiver: 0,
			Data:     One,
		}, s)
		if err != nil {
			t.Errorf("some error in decided one: %s", err.Error())
			return
		}
		if !bytes.Equal(m.Data, One) {
			t.Error("all nodes input 1 but output 0")
		}
	}, "./mock/config1.yaml")
	nls.Run()
}

// Test the condition that all node propose 0
func TestNewMMRABA2(t *testing.T) {
	nls := config.NewLocalStartWithReadLocalConfig(func(s config.Start) {
		conf := s.GetConfig()
		m, err := ABADecided(&pb.Message{
			Id:       "testMMR2",
			Sender:   uint32(conf.MyID),
			Receiver: 0,
			Data:     Zero,
		}, s)
		if err != nil {
			t.Errorf("some error in decided zero: %s", err.Error())
			return
		}
		if !bytes.Equal(m.Data, Zero) {
			t.Error("all nodes input 0 but output 1")
		}

	}, "./mock/config1.yaml")
	nls.Run()
}

// Test the condition that all honest node propose 1 and dishonest node propose 0
func TestNewMMRABA3(t *testing.T) {
	nls := config.NewLocalStartWithReadLocalConfig(func(s config.Start) {
		conf := s.GetConfig()
		var m *pb.Message
		var err error
		if conf.MyID <= 2 { // 0, 1, 2
			m, err = ABADecided(&pb.Message{
				Id:       "testMMR3",
				Sender:   uint32(conf.MyID),
				Receiver: 0,
				Data:     One,
			}, s)
		} else { // 3
			m, err = ABADecided(&pb.Message{
				Id:       "testMMR3",
				Sender:   uint32(conf.MyID),
				Receiver: 0,
				Data:     Zero,
			}, s)
		}
		if err != nil {
			t.Errorf("some error in decided zero: %s", err.Error())
			return
		}
		if !bytes.Equal(m.Data, One) {
			t.Error("all nodes input 0 but output 1")
		}

	}, "./mock/config1.yaml")
	nls.Run()
}

func TestNewMMRABA4(t *testing.T) {
	nls := config.NewLocalStartWithReadLocalConfig(func(s config.Start) {
		conf := s.GetConfig()
		var m *pb.Message
		var err error
		if conf.MyID <= 1 { // 0, 1
			m, err = ABADecided(&pb.Message{
				Id:       "testMMR4",
				Sender:   uint32(conf.MyID),
				Receiver: 0,
				Data:     One,
			}, s)
		} else { // 2, 3
			m, err = ABADecided(&pb.Message{
				Id:       "testMMR4",
				Sender:   uint32(conf.MyID),
				Receiver: 0,
				Data:     Zero,
			}, s)
		}
		if err != nil {
			t.Errorf("some error in decided: %s", err.Error())
			return
		}
		if !(bytes.Equal(m.Data, One) || bytes.Equal(m.Data, Zero)) {
			t.Error("nodes output weired values ")
		}

	}, "./mock/config1.yaml")
	nls.Run()
}
