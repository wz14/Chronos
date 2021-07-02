package consensus

import (
	"acc/config"
	"acc/pb"
	"bytes"
	"strconv"
	"testing"
)

func TestNewHB(t *testing.T) {
	nls := config.NewLocalStart(func(s config.Start) {
		conf := s.GetConfig()
		ml, _, err := Consensus(HBCon, &pb.Message{
			Id:       "Consen",
			Sender:   uint32(conf.MyID),
			Receiver: 0,
			Data:     []byte("mockdata" + strconv.Itoa(conf.MyID)),
		}, s)
		if err != nil {
			t.Error("consensus fail")
		}
		for i, m := range ml {
			if !bytes.Equal(m.Data, []byte("mockdata"+
				strconv.FormatUint(uint64(m.Sender), 10))) {
				t.Error("hb consensus bad value")
			}
			t.Logf("%d values: %s", i, string(m.Data))
		}

		if len(ml) < 2*conf.F+1 {
			t.Error("too few values in common set")
		}
	}, "./mock/config1.yaml")
	nls.Run()
}
