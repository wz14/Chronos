package mock

import (
	"acc/config"
	"acc/idchannel"
	"acc/logger"
	"acc/pb"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

func NewLocalStart() *LocalStart {
	lo := logger.NewLogger("main")
	c, err := config.NewConfig("./mock/config1.yaml") // get a no-pointer config
	if err != nil {
		lo.Fatalf("read config fail: %s", err.Error())
	}
	return &LocalStart{
		c: c,
		l: logger.NewLogger("main"),
	}
}

type LocalStart struct {
	c   config.Config
	nig *idchannel.NodeIDGroup
	pig *idchannel.PIDGroup
	l   *logger.Logger
}

func (s *LocalStart) Getnig() *idchannel.NodeIDGroup {
	return s.nig
}

func (s *LocalStart) Getpig() *idchannel.PIDGroup {
	return s.pig
}

func (s *LocalStart) GetConfig() *config.Config {
	return &s.c
}

func (s *LocalStart) CopySelf(id int) LocalStart {
	newc := s.c
	newc.MyID = id
	return LocalStart{c: newc, l: logger.NewLoggerWithID("main", id)}
}

func (s *LocalStart) MockRun(f func(s *LocalStart)) {
	if s.c.Isremote {
		s.l.Fatal("no implement remote deployment setting")
	}
	wg.Add(s.c.N)
	for i := 0; i < s.c.N; i++ {
		news := s.CopySelf(i)
		go func(news *LocalStart) {
			localaddress := news.c.IpList[news.c.MyID] + ":" + strconv.Itoa(news.c.PortList[news.c.MyID])
			news.l.Infof("honest node id:%d in %s", news.c.MyID, localaddress)

			// init pid group
			var err error
			news.pig, err = idchannel.NewPIDGroup(&news.c)
			if err != nil {
				news.l.Fatal("primitive group create fail: %s", err.Error())
			}

			//create server && bind classifier
			lis, err := net.Listen("tcp", ":"+strconv.Itoa(news.c.PortList[news.c.MyID]))
			if err != nil {
				news.l.Fatalf("tcp port open fail: %s in %d server", err, news.c.MyID)
			}
			defer lis.Close()
			server := grpc.NewServer()
			pb.RegisterNodeConServer(server, idchannel.NewClassifier(&news.c, news.c.MyID, news.pig))
			go server.Serve(lis)

			// wait time may config in file
			time.Sleep(2 * time.Second)

			// create id group
			news.nig, err = idchannel.NewIDGroup(&news.c)
			if err != nil {
				news.l.Fatalf("group create fail: %s", err.Error())
			}

			f(news)
			wg.Done()

		}(&news)
	}
	wg.Wait()
}

func (s *LocalStart) Run() {

}

// TODO: wait others node start, e.g. 10s in config file
// build id system for self
// 1. start local grpc server for receive message
// 2. start classifier in id channel
// 3. wait some time for others' nodes start
// 4. create id group for all nodes

// start log
