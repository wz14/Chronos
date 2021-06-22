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
	"sync"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	t.Log("is it run?")
	s := mock.NewStart()
	s.Run(func(smock mock.Start, wg *sync.WaitGroup) {
		t.Log("is it run?")
		lis, err := net.Listen("tcp", ":"+strconv.Itoa(smock.C.PortList[smock.C.MyID]))
		if err != nil {
			t.Fatalf("tcp port open fail: %s in %d server", err, smock.C.MyID)
		}
		defer lis.Close()
		server := grpc.NewServer()

		smock.Pig, err = idchannel.NewPIDGroup(&smock.C)
		if err != nil {
			log.Fatalf("primitive group create fail: %s", err.Error())
		}

		pb.RegisterNodeConServer(server, idchannel.NewClassifier(&smock.C, smock.C.MyID, smock.Pig))
		go server.Serve(lis)

		// wait time
		time.Sleep(1 * time.Second)

		// create id group
		smock.Nig, err = idchannel.NewIDGroup(&smock.C)
		if err != nil {
			t.Fatalf("group create fail: %s", err.Error())
		}
		t.Log("is it run?")
		if smock.C.MyID == 1 {
			err := Send(pb.Message{
				Id:       "Send_1",
				Sender:   1,
				Receiver: 2,
				Data:     []byte("what are you doing"),
			}, smock.Pig.GetRootPID("Send_1"), smock.Nig)
			if err != nil {
				t.Fatalf("send error: %s", err.Error())
			}
		}

		if smock.C.MyID == 2 {
			t.Log("run this")
			m, err := Receive(smock.Pig.GetRootPID("Send_1"))
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

		wg.Done()
	})
}

func TestReceive(t *testing.T) {

}
