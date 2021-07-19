package rbc

import (
	"acc/config"
	"acc/pb"
	"bytes"
	"testing"
)

func TestAVIDRBC(t *testing.T) {
	s := config.NewLocalStartWithReadLocalConfig(func(s config.Start) {
		pig := s.Getpig()
		nig := s.Getnig()
		config := s.GetConfig()
		nid, err := nig.GetID(1)
		if err != nil {
			t.Errorf("bad get id for 1 in %d node", config.MyID)
		}
		pid := pig.GetRootPID("RBC1")
		avidrbc := NewAVIDRBC(nid, pid, s)
		if config.MyID == 0 {
			go avidrbc.Broadcast(&pb.Message{
				Id:       "pickpickpick",
				Sender:   0,
				Receiver: 0,
				Data:     []byte("you are welcome"),
			})
		}
		mes, err := avidrbc.Deliver()
		t.Logf("%d node receive message: %s", config.MyID, string(mes.Data))
		if err != nil {
			t.Errorf("avidrbc deliver fail for %d node in %s avidrbc",
				config.MyID, avidrbc.s)
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
	}, "./mock/config1.yaml")
	s.Run()
}
