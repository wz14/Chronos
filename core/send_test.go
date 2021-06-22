package core

import (
	"acc/core/mock"
	"acc/idchannel"
	"acc/pb"
	"bytes"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	s := mock.NewStart()

	s.Run(func(s mock.Start) {
		lis, err := net.Listen("tcp", ":"+strconv.Itoa(s.C.PortList[s.C.MyID]))
		if err != nil {
			t.Fatalf("tcp port open fail: %s in %d server", err, s.C.MyID)
		}
		server := grpc.NewServer()

		s.Pig, err = idchannel.NewPIDGroup(&s.C)
		if err != nil {
			log.Fatalf("primitive group create fail: %s", err.Error())
		}

		pb.RegisterNodeConServer(server, idchannel.NewClassifier(&s.C, s.C.MyID, s.Pig))
		go server.Serve(lis)

		// wait time
		time.Sleep(5 * time.Second)

		// create id group
		s.Nig, err = idchannel.NewIDGroup(&s.C)
		if err != nil {
			t.Fatalf("group create fail: %s", err.Error())
		}

		if s.C.MyID == 1 {
			err := Send(pb.Message{
				Id:       "Send_1",
				Sender:   1,
				Receiver: 2,
				Data:     []byte("what are you doing"),
			}, s.Pig.GetRootPID("Send_1"), s.Nig)
			if err != nil {
				t.Fatalf("send error: %s", err.Error())
			}
		}

		if s.C.MyID == 2 {
			m, err := Receive(s.Pig.GetRootPID("Send_1"))
			if err != nil {
				t.Fatalf("receive error: %s", err.Error())
			}
			if m.Id != "Send_1" {
				t.Error("id is not Send_1 in m")
			}
			if bytes.Equal(m.Data, []byte("what are you doing")) {
				t.Error("data is not 'what are you doing' in m")
			}
			if m.Sender != 1 {
				t.Error("sender is not 1 in m")
			}
			if m.Receiver != 2 {
				t.Error("receiver is not 2 in m")
			}
		}

		s.Wg.Wait()
	})
}

func TestReceive(t *testing.T) {

}
