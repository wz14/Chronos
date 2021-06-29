package consensus

import (
	"acc/acs"
	"acc/config"
	"acc/core"
	"acc/crypto"
	"acc/idchannel"
	"acc/logger"
	"acc/pb"
	"github.com/golang/protobuf/proto"
	"strconv"
)

func NewHB(pid *idchannel.PrimitiveID, s config.Start) (Consen, error) {
	conf := s.GetConfig()
	return &HB{
		rootpid: pid,
		e:       s.GetEConfig(),
		pig:     s.Getpig(),
		c:       conf,
		s:       s,
		l:       logger.NewLoggerWithID("HB", conf.MyID),
	}, nil
}

type HB struct {
	rootpid *idchannel.PrimitiveID
	e       *crypto.TPKE
	pig     *idchannel.PIDGroup
	c       *config.Config
	s       config.Start
	l       *logger.Logger
}

func (h *HB) Propose(value *pb.Message) ([]*pb.Message, error) {
	// encrypt value.Data
	cpid := h.pig.GetChildPID("ACS", h.rootpid)

	// marshal value
	byt, err := proto.Marshal(value)
	if err != nil {
		h.l.Error("marshal fail")
	}

	ct := h.e.Enc(byt)

	msgs := h.acs(&pb.Message{
		Id:       cpid.Id,
		Sender:   uint32(h.c.MyID),
		Receiver: 0,
		Data:     ct,
	})

	return msgs, nil
}

func (h *HB) acs(value *pb.Message) []*pb.Message {
	values, err := acs.ACSDecided(value, h.s)
	if err != nil {
		h.l.Errorf("ACS run fail: %s", err.Error())
	}
	length := len(values)
	results := make(chan *pb.Message, length)
	// broadcast decrypt shares
	for i, v := range values {
		childp := h.pig.GetChildPID("Dec["+strconv.Itoa(i)+"]", h.rootpid)

		go func() {
			core.BroadCast(&pb.Message{
				Id:       childp.Id,
				Sender:   uint32(h.c.MyID),
				Receiver: 0,
				Data:     h.e.DecShare(v.Data),
			}, h.s)
		}()

		go func() {
			m := map[int][]byte{}
			for {
				decshare, err := core.Receive(childp.Id, h.s)
				if err != nil {
					h.l.Error("receive error")
				}
				m[int(decshare.Sender)] = decshare.Data
				if len(m) >= h.c.F+1 {
					byts := h.e.Dec(m, v.Data)
					// unmarshal realv
					realv := &pb.Message{}
					err := proto.Unmarshal(byts, realv)
					if err != nil {
						h.l.Error("unmarshal fail")
					}
					results <- realv
					return
				}
			}
		}()
	}
	msgs := []*pb.Message{}
	for i := 0; i < length; i++ {
		msg := <-results
		msgs = append(msgs, msg)
	}
	return msgs
}
