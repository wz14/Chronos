package idchannel

import (
	"acc/logger"
	"sync"
	"testing"
)

func GetMockPIG() *PIDGroup {
	pig := &PIDGroup{
		PIDmap: sync.Map{},
		l:      logger.NewLogger("test"),
	}
	return pig
}

func TestPIDGroup_GetRootPID(t *testing.T) {
	pig := GetMockPIG()
	p := pig.GetRootPID("feelgood")
	if p.Id != "feelgood" {
		t.Errorf("create pid fail with id: %s", p.Id)
	}
	if p.C == nil {
		t.Errorf("create pid fail and no init C with id: %s", p.Id)
	}
}

func TestPIDGroup_GetChildPID(t *testing.T) {
	pig := GetMockPIG()
	p := pig.GetRootPID("feelgood")
	if p.Id != "feelgood" {
		t.Errorf("create pid fail with id: %s", p.Id)
	}
	if p.C == nil {
		t.Errorf("create pid fail and no init C with id: %s", p.Id)
	}
	cpid := pig.GetChildPID("child", p)
	if cpid.Id != "feelgood"+sep+"child" {
		t.Errorf("get child pid fail")
	}
}

func TestPIDGroup_GetParentPID(t *testing.T) {
	pig := GetMockPIG()
	p := pig.GetRootPID("feelgood" + sep + "child")
	if p.Id != "feelgood"+sep+"child" {
		t.Errorf("create pid fail with id: %s", p.Id)
	}
	if p.C == nil {
		t.Errorf("create pid fail and no init C with id: %s", p.Id)
	}
	cpid := pig.GetParentPID(p)
	if cpid.Id != "feelgood" {
		t.Errorf("get child pid fail")
	}
}

func TestPIDGroup_GetInitRoundPID(t *testing.T) {
	pig := GetMockPIG()
	p := pig.GetRootPID("feelgood")
	if p.Id != "feelgood" {
		t.Errorf("create pid fail with id: %s", p.Id)
	}
	if p.C == nil {
		t.Errorf("create pid fail and no init C with id: %s", p.Id)
	}
	cpid := pig.GetInitRoundPID(p)
	if cpid.Id != "feelgood"+sep+"1" {
		t.Errorf("get child pid fail")
	}
}

func TestPIDGroup_GetNextRoundPID(t *testing.T) {
	pig := GetMockPIG()
	p := pig.GetRootPID("feelgood" + sep + "1")
	if p.Id != "feelgood"+sep+"1" {
		t.Errorf("create pid fail with id: %s", p.Id)
	}
	if p.C == nil {
		t.Errorf("create pid fail and no init C with id: %s", p.Id)
	}
	cpid := pig.GetNextRoundPID(p)
	if cpid.Id != "feelgood"+sep+"2" {
		t.Errorf("get child pid fail")
	}

}

func TestPIDGroup_GetNextRoundPID2(t *testing.T) {
	pig := GetMockPIG()
	p := pig.GetRootPID("feelgood" + sep + "999")
	if p.Id != "feelgood"+sep+"999" {
		t.Errorf("create pid fail with id: %s", p.Id)
	}
	if p.C == nil {
		t.Errorf("create pid fail and no init C with id: %s", p.Id)
	}
	cpid := pig.GetNextRoundPID(p)
	if cpid.Id != "feelgood"+sep+"1000" {
		t.Errorf("get child pid fail")
	}

}
