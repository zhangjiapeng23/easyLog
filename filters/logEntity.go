package filters

import (
	"easyLog/util"
	"fmt"
	"sync"
	"time"
)

type Log struct {
	DatetimeStr string
	Level       string
	title       string
	LogDetail   [][]byte
	IsMatch     bool
	originLog   []byte
}

func NewLog(dateTime string, level string, title string) (log *Log) {
	utcTime, _ := time.Parse("2006-01-02 15:04:05", dateTime)
	cstSh, _ := time.LoadLocation("Asia/Shanghai")
	dateTime = utcTime.In(cstSh).Format("2006/01/02 15:04:05")
	return &Log{dateTime, level, title, make([][]byte, 0), false, make([]byte, 0)}
}

func NewOriginLog(logLine []byte) (log *Log) {
	return &Log{"", "", "", make([][]byte, 0), false, logLine}
}

func (l *Log) AddLogDetail(error []byte) {
	l.LogDetail = append(l.LogDetail, error)
}

func (l *Log) String() {
	var printLock sync.Mutex
	printLock.Lock()
	defer printLock.Unlock()
	if len(l.originLog) > 0 {
		fmt.Print(string(l.originLog))
	} else {
		fmt.Printf("%s【%s】%s\n", l.DatetimeStr, l.Level, l.title)
		if len(l.LogDetail) > 0 {
			fmt.Println("日志详情:")
			for _, logLine := range l.LogDetail {
				fmt.Print(string(logLine))
			}
		}
		util.PrintSplitLine("-")
	}
}
