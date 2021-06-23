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

func TestBroadCast(t *testing.T) {
	s := mock.NewStart()

	s.MockRun(func(s mock.Start, wg *sync.WaitGroup) {
		lis, err := net.Listen("tcp", ":"+strconv.Itoa(s.C.PortList[s.C.MyID]))
		if err != nil {
			t.Fatalf("tcp port open fail: %s in %d server", err, s.C.MyID)
		}
		defer lis.Close()
		server := grpc.NewServer()

		s.Pig, err = idchannel.NewPIDGroup(&s.C)
		if err != nil {
			log.Fatalf("primitive group create fail: %s", err.Error())
		}

		pb.RegisterNodeConServer(server, idchannel.NewClassifier(&s.C, s.C.MyID, s.Pig))
		go server.Serve(lis)

		// wait time
		time.Sleep(1 * time.Second)

		// create id group
		s.Nig, err = idchannel.NewIDGroup(&s.C)
		if err != nil {
			t.Fatalf("group create fail: %s", err.Error())
		}

		go func() {
			t.Log("enter sender routine")
			if s.C.MyID == 1 {
				err := BroadCast(&pb.Message{
					Id:       "BC_1",
					Sender:   1,
					Receiver: 0,
					Data:     []byte("what are you doing"),
				}, &s)
				if err != nil {
					t.Fatalf("send error: %s", err.Error())
				}
			}
		}()

		if true {
			t.Logf("enter receiver routine %d", s.C.MyID)
			m, err := Receive("BC_1", &s)
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
			if m.Receiver != uint32(s.C.MyID) {
				t.Errorf("receiver is %d in m but I'm %d", m.Receiver, s.C.MyID)
			}
		}
		t.Log("done with ", s.C.MyID)
		wg.Done()
	})
}
