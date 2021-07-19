package core

import (
	"acc/config"
	"acc/pb"
	"bytes"
	"testing"
)

func TestBroadCast(t *testing.T) {
	config.NewLocalStartWithReadLocalConfig(func(s config.Start) {
		conf := s.GetConfig()
		t.Log("enter sender routine")
		if conf.MyID == 1 {
			err := BroadCast(&pb.Message{
				Id:       "BC_1",
				Sender:   1,
				Receiver: 0,
				Data:     []byte("what are you doing"),
			}, s)
			if err != nil {
				t.Fatalf("send error: %s", err.Error())
			}
		}

		if true {
			t.Logf("enter receiver routine %d", conf.MyID)
			m, err := Receive("BC_1", s)
			if err != nil {
				t.Fatalf("receive error: %s", err.Error())
			}
			if m.Id != "BC_1" {
				t.Error("id is not Send_1 in m")
			}
			if !bytes.Equal(m.Data, []byte("what are you doing")) {
				t.Error("data is not 'what are you doing' in m")
			}
			if m.Sender != 1 {
				t.Error("sender is not 1 in m")
			}
			if m.Receiver != uint32(conf.MyID) {
				t.Errorf("receiver is %d in m but I'm %d", m.Receiver, conf.MyID)
			}
		}
		t.Log("done with ", conf.MyID)
	}, "./mock/config1.yaml")
}
