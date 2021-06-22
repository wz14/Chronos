package mock

import (
	"acc/config"
	"acc/idchannel"
	"acc/logger"
	"sync"
)

var log = logger.NewLogger("main")

func NewStart() *Start {
	c, err := config.NewConfig("./mock/config1.yaml") // get a no-pointer config
	if err != nil {
		log.Fatalf("read config fail: %s", err.Error())
	}
	return &Start{C: c}
}

type Start struct {
	C   config.Config
	Nig *idchannel.NodeIDGroup
	Pig *idchannel.PIDGroup
}

func (s *Start) Getnig() *idchannel.NodeIDGroup {
	return s.Nig
}

func (s *Start) Getpig() *idchannel.PIDGroup {
	return s.Pig
}

func (s *Start) GetConfig() *config.Config {
	return &s.C
}

func (s *Start) CopySelf(id int) Start {
	newc := s.C
	newc.MyID = id
	return Start{C: newc}
}

func (s *Start) Run() {

}

func (s *Start) MockRun(f func(s Start, wg *sync.WaitGroup)) {
	wg := sync.WaitGroup{}
	if s.C.Isremote {
		log.Fatal("no implement remote deployment setting")
	}
	for i := 0; i < s.C.N; i++ {
		news := s.CopySelf(i)
		wg.Add(1)
		go f(news, &wg)
	}
	wg.Wait()
}
