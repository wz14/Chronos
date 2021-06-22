package idchannel

import (
	"acc/logger"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"strconv"
)

type Config interface {
	GetN() (int, error)
	GetIPList() ([]string, error)
	GetPortList() ([]int, error)
	GetMyID() (int, error)
}

type NodeIDGroup struct {
	NodeIDmap map[int]*NodeID
	log       *logger.Logger
}

type NodeID struct {
	ID      int
	Address string
	Connect *grpc.ClientConn
}

func NewIDGroup(c Config) (*NodeIDGroup, error) {
	myid, _ := c.GetMyID()
	l := logger.NewLoggerWithID("idchannel", myid)

	ipl, err := c.GetIPList()
	if err != nil {
		return nil, err
	}
	pol, err := c.GetPortList()
	if err != nil {
		return nil, err
	}
	n, err := c.GetN()
	if err != nil {
		return nil, err
	}

	l.Infof("create new group in 0 to %d", n-1)

	nig := NodeIDGroup{
		NodeIDmap: map[int]*NodeID{},
	}
	for i := 0; i < n; i++ {
		ID, err := NewID(i, ipl[i], pol[i])
		if err != nil {
			return nil, errors.Wrap(err, "create group fail")
		}
		nig.NodeIDmap[i] = ID
	}
	return &nig, nil
}

func NewID(id int, ip string, port int) (*NodeID, error) {
	i := NodeID{
		ID:      id,
		Address: ip + ":" + strconv.Itoa(port),
		Connect: nil,
	}
	c, err := grpc.Dial(i.Address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	i.Connect = c
	return &i, nil
}

func (n *NodeIDGroup) GetID(id int) (*NodeID, error) {
	i, ok := n.NodeIDmap[id]
	if !ok {
		return nil, errors.New("no such ID in ID system")
	}
	return i, nil
}
