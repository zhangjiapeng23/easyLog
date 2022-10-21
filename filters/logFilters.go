package filters

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var Done chan struct{}

func ErrorFilter(log chan []byte, filterLog chan *Log, extra ...string) {
	levelFilter("ERROR", log, filterLog, extra...)
}

func WarnFilter(log chan []byte, filterLog chan *Log, extra ...string) {
	levelFilter("WARN", log, filterLog, extra...)
}

func InfoFilter(log chan []byte, filterLog chan *Log, extra ...string) {
	levelFilter("INFO", log, filterLog, extra...)
}

func DebugFilter(log chan []byte, filterLog chan *Log, extra ...string) {
	levelFilter("DEBUG", log, filterLog, extra...)
}

func KeywordFilter(log chan []byte, filterLog chan *Log, extra ...string) {
	lineStartPattern, _ := regexp.Compile("[0-9]+-[0-9]+-[0-9]+.*")
	levelLogPattern, _ := regexp.Compile(".*(INFO|ERROR|DEBUG|WARN).*")
	dateTimePattern, _ := regexp.Compile("[0-9]+-[0-9]+-[0-9]+ [0-9]+:[0-9]+:[0-9]+")
	titlePattern, _ := regexp.Compile(`.*?\[.*?\](.*)`)
	keyword := ""
	if len(extra) > 0 {
		keyword = extra[0]
	}
	kwString := fmt.Sprintf(".*%s.*", keyword)
	keywordPattern, _ := regexp.Compile(kwString)
	var logEntity *Log
	for {
		select {
		case <-Done:
			return
		case logMsg := <-log:
			// 定位每条日志的开头
			if lineStartPattern.Match(logMsg) && levelLogPattern.Match(logMsg) {
				level := levelLogPattern.FindSubmatch(logMsg)[1]
				dateTimeStr := dateTimePattern.FindSubmatch(logMsg)[0]
				titleMatch := titlePattern.FindSubmatch(logMsg)
				title := ""
				if len(titleMatch) > 0 {
					title = strings.Trim(string(titleMatch[1]), " ")
				} else {
					title = strings.Trim(string(logMsg), " ")
				}
				// 当下一日志已经产生，就直接发送
				if logEntity != nil && logEntity.IsMatch && !logEntity.IsSend {
					logEntity.IsSend = true
					filterLog <- logEntity
				}
				logEntity = NewLog(string(dateTimeStr), string(level), title)
				if keywordPattern.Match(logMsg) {
					logEntity.IsMatch = true
					logEntity.Keyword = keyword
					go func(log *Log) {
						time.AfterFunc(2*time.Second, func() {
							if !log.IsSend {
								log.IsSend = true
								filterLog <- log
							}
						})
					}(logEntity)
				}
			} else if !lineStartPattern.Match(logMsg) && logEntity != nil {
				logEntity.AddLogDetail(logMsg)
				// if log detial match keyword and log is not ready to print
				if keywordPattern.Match(logMsg) && !logEntity.IsMatch {
					logEntity.IsMatch = true
					logEntity.Keyword = keyword
					go func(log *Log) {
						time.AfterFunc(2*time.Second, func() {
							if !log.IsSend {
								log.IsSend = true
								filterLog <- log
							}
						})
					}(logEntity)
				}
			}
		}
	}
}

func AllFilter(log chan []byte, filterLog chan *Log, extra ...string) {
	for {
		select {
		case <-Done:
			return
		case logMsg := <-log:
			filterLog <- NewOriginLog(logMsg)
		}
	}
}

func levelFilter(level string, log chan []byte, filterLog chan *Log, extra ...string) {
	levelPatterString := fmt.Sprintf(".*(%s).*", level)
	lineStartPattern, _ := regexp.Compile("[0-9]+-[0-9]+-[0-9]+.*")
	levelLogPattern, _ := regexp.Compile(levelPatterString)
	dateTimePattern, _ := regexp.Compile("[0-9]+-[0-9]+-[0-9]+ [0-9]+:[0-9]+:[0-9]+")
	titlePattern, _ := regexp.Compile(`.*?\[.*?\](.*)`)
	var logEntity *Log
	for {
		select {
		case <-Done:
			return
		case logMsg := <-log:
			if lineStartPattern.Match(logMsg) && levelLogPattern.Match(logMsg) {
				level := levelLogPattern.FindSubmatch(logMsg)[1]
				dateTimeStr := dateTimePattern.FindSubmatch(logMsg)[0]
				titleMatch := titlePattern.FindSubmatch(logMsg)
				title := ""
				if len(titleMatch) > 0 {
					title = strings.Trim(string(titleMatch[1]), " ")
				} else {
					title = strings.Trim(string(logMsg), " ")
				}
				// 当下一日志已经产生，就直接发送
				if logEntity != nil && !logEntity.IsSend {
					logEntity.IsSend = true
					filterLog <- logEntity
				}
				logEntity = NewLog(string(dateTimeStr), string(level), title)
				go func(log *Log) {
					time.AfterFunc(2*time.Second, func() {
						if !log.IsSend {
							log.IsSend = true
							filterLog <- log
						}
					})
				}(logEntity)
			} else if !lineStartPattern.Match(logMsg) && logEntity != nil {
				logEntity.AddLogDetail(logMsg)
			}
		}
	}
}
