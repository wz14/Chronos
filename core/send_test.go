package core

import (
	"acc/config"
	"acc/pb"
	"bytes"
	"testing"
)

func TestSend(t *testing.T) {
	nls := config.NewLocalStart(func(s config.Start) {
		conf := s.GetConfig()
		if conf.MyID == 1 {
			err := Send(&pb.Message{
				Id:       "Send_1",
				Sender:   1,
				Receiver: 2,
				Data:     []byte("what are you doing"),
			}, s)
			if err != nil {
				t.Fatalf("send error: %s", err.Error())
			}
		}

		if conf.MyID == 2 {
			t.Log("run this")
			m, err := Receive("Send_1", s)
			if err != nil {
				t.Fatalf("receive error: %s", err.Error())
			}
			if m.Id != "Send_1" {
				t.Error("id is not Send_1 in m")
			}
			if !bytes.Equal(m.Data, []byte("what are you doing")) {
				t.Error("data is not 'what are you doing' in m")
			}
			if m.Sender != 1 {
				t.Error("sender is not 1 in m")
			}
			if m.Receiver != 2 {
				t.Error("receiver is not 2 in m")
			}
		}

	}, "./mock/config1.yaml")
	nls.Run()
}

func TestReceive(t *testing.T) {

}
