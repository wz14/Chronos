package consensus

import (
	"acc/config"
	"acc/crypto"
	"acc/pb"
	"bytes"
	"strconv"
	"testing"
)

func TestNewCHOS(t *testing.T) {
	nls := config.NewLocalStartWithReadLocalConfig(func(s config.Start) {
		conf := s.GetConfig()
		pid := s.Getpig().GetRootPID("CHOS1")
		c, err := NewCHOS(pid, s)
		if err != nil {
			t.Errorf("create CHOS fail: %s", err.Error())
		}
		ml, ts, err := c.Propose(&pb.Message{
			Id:       "CHOS1",
			Sender:   uint32(conf.MyID),
			Receiver: 0,
			Data:     []byte("mockdata" + strconv.Itoa(conf.MyID)),
		})
		if err != nil {
			t.Errorf("propose CHOS fail: %s", err.Error())
		}

		for i, m := range ml {
			if !bytes.Equal(m.Data, []byte("mockdata"+
				strconv.FormatUint(uint64(m.Sender), 10))) {
				t.Error("acs decided bad value")
			}
			if !bytes.Equal(ts[i].Hash, crypto.Hash(m.Data)) {
				t.Error("hash of ts is not that transaction")
			}
			if int(ts[i].Num) != len(m.Data)/250 {
				t.Error("number of timestamps is not same with length of tx")
			}
			t.Logf("%d values: %s", i, string(m.Data))
		}

	}, "./mock/config1.yaml")
	nls.Run()
}
