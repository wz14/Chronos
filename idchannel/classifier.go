package idchannel

import (
	"acc/logger"
	"acc/pb"
	"context"
)

// create a classifier to distribute message to different primitiveID channel

func NewClassifier(C Config, id int, pidg *PIDGroup) *Classifier {
	c := Classifier{
		C:   C,
		l:   logger.NewLoggerWithID("idchannel", id),
		pig: pidg,
	}
	return &c
}

type Classifier struct {
	C   Config
	l   *logger.Logger
	pig *PIDGroup
}

func (c *Classifier) SendMessage(_ context.Context, m *pb.Message) (*pb.Zero, error) {
	c.l.Infof("receive %s message from %d", m.Id, m.Sender)
	c.l.Debugf("receive %s message from %d: %v:", m.Id, m.Sender, m.Data)
	pid := c.pig.GetRootPID(m.Id)
	pid.C <- m
	return &pb.Zero{}, nil
}
