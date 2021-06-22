package mock

import (
	"acc/idchannel"
	"acc/logger"
	"sync"
)

var log = logger.NewLogger("main")

func NewStart() *Start {
	c, err := NewConfig() // get a no-pointer config
	if err != nil {
		log.Fatalf("read config fail: %s", err.Error())
	}
	return &Start{C: c}
}

type Start struct {
	C   Config
	Nig *idchannel.NodeIDGroup
	Pig *idchannel.PIDGroup
}

func (s *Start) CopySelf(id int) Start {
	newc := s.C
	newc.MyID = id
	return Start{C: newc}
}

func (s *Start) Run(f func(s Start, wg *sync.WaitGroup)) {
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
