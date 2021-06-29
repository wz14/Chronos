package aba

import (
	"acc/config"
	"acc/core"
	"acc/idchannel"
	"acc/logger"
	"acc/pb"
	"bytes"
	"sync"
)

/*

A MMR ABA is from Signature-free asynchronous Byzantine
consensus with t< n/3 and O (n2) messages paper 2014.

input: x=0/1

round increment value r:

	broadcast BVAL(x)

	wait f+1 BVAL(x)
		broadcast BVAL(x)
	wait 2f+1 BVAL(x)
		broadcast AUX(w)

	wait collect n-t AUX message
		generate values \in bin_values
		send part-sig

	wait collect t+1 part-sig
		generate random coin
		if decided, output
		open and enter next round. run forever
*/

func NewMMRABA(pid *idchannel.PrimitiveID, s config.Start) *MMRABA {
	c := s.GetConfig()
	m := &MMRABA{
		rootpid: pid,
		c:       c,
		nig:     s.Getnig(),
		pig:     s.Getpig(),
		l:       logger.NewLoggerWithID("MMRABA", c.MyID),
		s:       s,
		sid:     []byte(pid.Id),
	}
	return m
}

type MMRABA struct {
	rootpid *idchannel.PrimitiveID
	nig     *idchannel.NodeIDGroup
	pig     *idchannel.PIDGroup
	l       *logger.Logger
	c       *config.Config
	s       config.Start
	sid     []byte
}

// TODO: add context cancel mechanism
// Note: this function is blocking
func (m *MMRABA) Decided(message *pb.Message) (*pb.Message, error) {
	// go roundDecided (pid.0)
	// message <- pid.0
	// r = r + 1
	initRoundpid := m.pig.GetInitRoundPID(m.rootpid)
	rm := &roundmmr{
		pid:       initRoundpid,
		nig:       m.nig,
		pig:       m.pig,
		l:         logger.NewLoggerWithID(initRoundpid.Id, m.c.MyID),
		c:         m.c,
		s:         m.s,
		sid:       m.sid,
		values:    make(chan []byte, 1),
		est:       message.Data,
		isBVone:   false,
		isBVzero:  false,
		isAUX:     false,
		mu:        sync.Mutex{},
		binvalues: []byte{},
	}
	rm.roundDecided()
	// handle always decides mechaism
	go func() {
		isDecided := false
		for mess := range initRoundpid.C {
			if !isDecided {
				m.rootpid.C <- mess
				isDecided = true
			} else {
				m.l.Infof("mmr aba output mes: %s, but discard", string(mess.Data))
			}
		}
	}()
	mes := <-m.rootpid.C
	m.l.Infof("mmr aba output first mes: %s", string(mes.Data))
	return mes, nil
}

type roundmmr struct {
	pid *idchannel.PrimitiveID
	nig *idchannel.NodeIDGroup
	pig *idchannel.PIDGroup
	l   *logger.Logger
	c   *config.Config
	s   config.Start

	est     []byte
	nexest  []byte
	sid     []byte
	nextsid []byte
	values  chan []byte

	// lock area
	mu        sync.Mutex
	isBVone   bool
	isBVzero  bool
	isAUX     bool
	binvalues []byte
	// lock area

}

func (m *roundmmr) roundDecided() {
	// use pid in args and not use pid in struct
	// pid <- round
	// go bvBroadcast
	go m.bvBroadcast()
	go m.bvDeliver()
	go m.auxDeliver()
	go m.ccDeliver()
	// go bvDeliver && broadcast aux
	// go collect aux && run random coin
	// go wait random coin &&
	//(pid.0 <- message ||
	// new round
	//	go roundDecide(pid.1);message <- pid.1;pid.0 <- message)
	// mes := <-m.pid.C
	// m.pid.C <- mes
}

func (m *roundmmr) bvBroadcast() {
	childp := m.pig.GetChildPID("bvB", m.pid)
	if bytes.Equal(m.est, One) {
		m.mu.Lock()
		m.isBVone = true
		m.mu.Unlock()
	}
	if bytes.Equal(m.est, Zero) {
		m.mu.Lock()
		m.isBVzero = true
		m.mu.Unlock()
	}
	core.BroadCast(&pb.Message{
		Id:       childp.Id,
		Sender:   uint32(m.c.MyID),
		Receiver: 0,
		Data:     m.est,
	}, m.s)
}

func (m *roundmmr) bvDeliver() {
	childp := m.pig.GetChildPID("bvB", m.pid)
	onecount := 0
	zerocount := 0
	for {
		mes, err := core.Receive(childp.Id, m.s)
		if err != nil {
			m.l.Error("receive bv value fail")
		}

		if bytes.Equal(mes.Data, One) {
			onecount += 1
		} else if bytes.Equal(mes.Data, Zero) {
			zerocount += 1
		} else {
			m.l.Error("receive undefined value (not 0/1)")
		}

		m.l.Infof("receive %s from %d: %d one, %d zero",
			string(mes.Data), mes.Sender, onecount, zerocount)

		// reuse code maybe ??
		if onecount >= m.c.F+1 {
			m.mu.Lock()
			var tmp bool
			if !m.isBVone {
				m.isBVone = true
				tmp = false
			} else {
				tmp = true
			}
			m.mu.Unlock()
			if !tmp {
				core.BroadCast(&pb.Message{
					Id:       childp.Id,
					Sender:   uint32(m.c.MyID),
					Receiver: 0,
					Data:     One,
				}, m.s)

			}
		}

		if zerocount >= m.c.F+1 {
			m.mu.Lock()
			var tmp bool
			if !m.isBVzero {
				m.isBVzero = true
				tmp = false
			} else {
				tmp = true
			}
			m.mu.Unlock()
			if !tmp {
				core.BroadCast(&pb.Message{
					Id:       childp.Id,
					Sender:   uint32(m.c.MyID),
					Receiver: 0,
					Data:     Zero,
				}, m.s)

			}
		}

		// 2f + 1 conditions, reuse code
		if zerocount >= m.c.F*2+1 {
			m.l.Infof("0 b_val enough and add to binvalues")
			m.mu.Lock()
			if !bytes.Contains(m.binvalues, Zero) {
				m.binvalues = append(m.binvalues, Zero...)
			}
			m.mu.Unlock()
		}

		if onecount >= m.c.F*2+1 {
			m.l.Infof("1 b_val enough and add to binvalues")
			m.mu.Lock()
			if !bytes.Contains(m.binvalues, One) {
				m.binvalues = append(m.binvalues, One...)
			}
			m.mu.Unlock()
		}

		m.mu.Lock()
		if m.isAUX || bytes.Equal(m.binvalues, []byte("")) {
			m.mu.Unlock()
			continue
		} else {
			m.isAUX = true
			m.mu.Unlock()
		}

		m.l.Infof("broadcat aux message")
		childp2 := m.pig.GetChildPID("aux", m.pid)
		if bytes.Contains(m.binvalues, One) {
			core.BroadCast(&pb.Message{
				Id:       childp2.Id,
				Sender:   uint32(m.c.MyID),
				Receiver: 0,
				Data:     One,
			}, m.s)
		} else if bytes.Contains(m.binvalues, Zero) {
			core.BroadCast(&pb.Message{
				Id:       childp2.Id,
				Sender:   uint32(m.c.MyID),
				Receiver: 0,
				Data:     Zero,
			}, m.s)
		} else {
			m.l.Error("code error")
		}
	}
}

func (m *roundmmr) auxDeliver() {
	childp := m.pig.GetChildPID("aux", m.pid)
	auxcount := 0
	values := []byte{}
	isRandomCoin := false
	for {
		mes, err := core.Receive(childp.Id, m.s)
		if err != nil {
			m.l.Error("recieve aux mes error")
		}
		m.l.Infof("receive aux message from %d", mes.Sender)
		auxcount += 1
		if !bytes.Contains(values, mes.Data) {
			// may check mes.Data is one or zero
			values = append(values, mes.Data...)
		}
		//TODO: neet do check all n-t messages is in bin_values

		if !isRandomCoin && auxcount >= m.c.N-m.c.F {
			m.l.Infof("collect enough aux message, values: %s", string(values))
			m.values <- values
			isRandomCoin = true
			id := m.pig.GetChildPID("ccoin", m.pid).Id
			partSignature, err := m.s.GetCConfig().Sign(m.sid)
			if err != nil {
				m.l.Errorf("sign fail to in %s: %s", id, err.Error())
			}
			m.l.Infof("broadcast part signature")
			core.BroadCast(&pb.Message{
				Id:       id,
				Sender:   uint32(m.c.MyID),
				Receiver: 0,
				Data:     partSignature, // add cConfig to config and mmr structure
				// add random sid to mmr structure.
			}, m.s)
			return
		}
	}
}

func (m *roundmmr) ccDeliver() {
	sigs := [][]byte{}
	for {
		mes, err := core.Receive(m.pig.GetChildPID("ccoin", m.pid).Id, m.s)
		if err != nil {
			m.l.Errorf("receive fail in ccoin: %s", err.Error())
			continue
		}
		m.l.Infof("receive ccoin from %d", mes.Sender)
		//combine sigs to get common coin
		sigs = append(sigs, mes.Data)
		if len(sigs) == m.s.GetCConfig().T {
			sig, err := m.s.GetCConfig().Combine(sigs, m.sid)
			if err != nil {
				m.l.Errorf("combine part signature fail: %s", err.Error())
			}
			m.nextsid = sig
			coin := []byte{sig[len(sig)-1]%2 + 48} // ascii("0/1") -> 48/49
			m.l.Infof("get coin is %s", string(coin))
			// TODO: judge decide is ok or not, try to decide
			vs := <-m.values
			m.l.Infof("vs:%s, coin:%s", string(vs), string(coin))
			if len(vs) == 1 {
				if bytes.Equal(vs, coin) {
					m.l.Infof("send %s to self parent channel", string(vs))
					core.Send(&pb.Message{
						Id:       m.pid.Id,
						Sender:   uint32(m.c.MyID),
						Receiver: uint32(m.c.MyID),
						Data:     vs,
					}, m.s)
				}
				m.nexest = vs
			} else {
				m.nexest = coin
			}
			// continue next round
			npid := m.pig.GetNextRoundPID(m.pid)
			rm := &roundmmr{
				pid:       npid,
				nig:       m.nig,
				pig:       m.pig,
				l:         logger.NewLoggerWithID(npid.Id, m.c.MyID),
				c:         m.c,
				s:         m.s,
				sid:       m.nextsid,
				values:    make(chan []byte, 1),
				est:       m.nexest,
				isBVone:   false,
				isBVzero:  false,
				isAUX:     false,
				mu:        sync.Mutex{},
				binvalues: []byte{},
			}

			go rm.roundDecided()
			for v := range npid.C {
				m.pid.C <- v
			}
			return
		}
	}
}
