package Benchmark

import (
	"acc/logger"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

// benchmark in go test seems can't handle multi-routine program

var l *logger.Logger = logger.NewLogger("Benchmark")
var BenchmarkTable sync.Map = sync.Map{}
var fileLocation, _ = os.LookupEnv("STATISTIC")

const (
	BEGIN = 0
	END   = 1
)

type internal struct {
	who  int
	T    int
	time int64
}

func Begin(id string, who int) {
	i := &internal{
		time: time.Now().UnixNano(),
		T:    BEGIN,
		who:  who,
	}
	BenchmarkTable.Store(id+" B "+strconv.Itoa(who), i)
}

func End(id string, who int) {
	i := &internal{
		time: time.Now().UnixNano(),
		T:    END,
		who:  who,
	}
	BenchmarkTable.Store(id+" E "+strconv.Itoa(who), i)
}

func BenchmarkOuput() {
	type internalWithId struct {
		id   string
		who  int
		time int64
		T    int
	}
	list := []*internalWithId{}
	BenchmarkTable.Range(
		func(key interface{}, value interface{}) bool {
			k := key.(string)
			v := value.(*internal)
			list = append(list, &internalWithId{
				id:   k,
				who:  v.who,
				time: v.time,
				T:    v.T,
			})
			return true
		})
	sort.Slice(list, func(i int, j int) bool {
		return len(list[i].id) < len(list[j].id)
	})

	f, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		l.Errorf("wirte log fail: %s", err.Error())
		return
	}
	defer f.Close()
	for _, s := range list {
		_, err := f.Write([]byte(s.id + " " + strconv.Itoa(s.who) +
			" " + strconv.FormatInt(s.time, 10) + " " + strconv.Itoa(s.T) + "\n"))
		if err != nil {
			l.Errorf("wirte log fail: %s", err.Error())
			return
		}
	}
	return
}
