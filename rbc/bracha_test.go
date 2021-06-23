package rbc

import (
	"acc/pb"
	"acc/rbc/mock"
	"bytes"
	"testing"
)

func TestBrachaRBC(t *testing.T) {
	s := mock.NewLocalStart()
	s.MockRun(func(s *mock.LocalStart) {
		pig := s.Getpig()
		nig := s.Getnig()
		config := s.GetConfig()
		nid, err := nig.GetID(1)
		if err != nil {
			t.Errorf("bad get id for 1 in %d node", config.MyID)
		}
		pid := pig.GetRootPID("RBC1")
		bracha := NewBrachaRBC(nid, pid, s)
		if config.MyID == 0 {
			go bracha.Broadcast(&pb.Message{
				Id:       "pickpickpick",
				Sender:   0,
				Receiver: 0,
				Data:     []byte("you are welcome"),
			})
		}
		mes, err := bracha.Deliver()
		t.Logf("%d node receive message: %s", config.MyID, string(mes.Data))
		if err != nil {
			t.Errorf("bracha deliver fail for %d node in %s bracha",
				config.MyID, bracha.s)
		}
		if mes.Id != "pickpickpick" {
			t.Errorf("mes.Id is not pickpick... , mes.Id is %s", mes.Id)
		}
		if mes.Sender != 0 {
			t.Errorf("mes.sender should be 0, but is %d", mes.Sender)
		}
		if !bytes.Equal([]byte("you are welcome"), mes.Data) {
			t.Errorf("mes.Data is %s", mes.Data)
		}
	})
}
