package filters

import (
	"easyLog/util"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

var (
	printBlue = color.New(color.FgBlue)
	printRed  = color.New(color.FgHiRed)
)

type Log struct {
	DatetimeStr string
	Level       string
	title       string
	LogDetail   [][]byte
	IsMatch     bool
	IsSend      bool
	Keyword     string
	originLog   []byte
}

func NewLog(dateTime string, level string, title string) (log *Log) {
	utcTime, _ := time.Parse("2006-01-02 15:04:05", dateTime)
	cstSh, _ := time.LoadLocation("Asia/Shanghai")
	dateTime = utcTime.In(cstSh).Format("2006/01/02 15:04:05")
	return &Log{dateTime, level, title, make([][]byte, 0), false, false, "", make([]byte, 0)}
}

func NewOriginLog(logLine []byte) (log *Log) {
	return &Log{"", "", "", make([][]byte, 0), false, false, "", logLine}
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
		if l.IsMatch && l.Keyword != "" {
			fmt.Printf("%s【%s】%s\n", colorKeyWord(l.DatetimeStr, l.Keyword),
				colorKeyWord(l.Level, l.Keyword), colorKeyWord(l.title, l.Keyword))
			if len(l.LogDetail) > 0 {
				printBlue.Println("日志详情:")
				for _, logLine := range l.LogDetail {
					fmt.Print(colorKeyWord(string(logLine), l.Keyword))
				}
			}
		} else {
			fmt.Printf("%s【%s】%s\n", l.DatetimeStr, l.Level, l.title)
			if len(l.LogDetail) > 0 {
				printBlue.Println("日志详情:")
				for index, logLine := range l.LogDetail {
					if index < 5 {
						printRed.Print(string(logLine))
					} else {
						fmt.Print(string(logLine))
					}

				}
			}
		}
		util.PrintSplitLine("-")
	}
}

func colorKeyWord(s string, keyword string) string {
	red := color.New(color.FgRed).SprintFunc()
	otherWords := strings.Split(s, keyword)
	return strings.Join(otherWords, fmt.Sprint(red(keyword)))
}
