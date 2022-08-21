package filters

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func ErrorFilter(log chan []byte, filterLog chan *Log, extra ...string) {
	levelFilter("ERROR", log, filterLog, extra...)
}

func WarnFilter(log chan []byte, filterLog chan *Log, extra ...string) {
	levelFilter("WARN", log, filterLog, extra...)
}

func InfoFilter(log chan []byte, filterLog chan *Log, extra ...string) {
	levelFilter("INFO", log, filterLog, extra...)
}

func KeywordFilter(log chan []byte, filterLog chan *Log, extra ...string) {
	lineStartPattern, _ := regexp.Compile("[0-9]+-[0-9]+-[0-9]+.*")
	levelLogPattern, _ := regexp.Compile(".*(INFO|ERROR).*")
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
		logMsg, ok := <-log
		if ok {
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
				logEntity = NewLog(string(dateTimeStr), string(level), title)
				if keywordPattern.Match(logMsg) {
					logEntity.IsMatch = true
					go func(log *Log) {
						time.AfterFunc(1*time.Second, func() {
							filterLog <- log
						})
					}(logEntity)
				}
			} else if !lineStartPattern.Match(logMsg) && logEntity != nil {
				logEntity.AddLogDetail(logMsg)
				// if log detial match keyword and log is not ready to print
				if keywordPattern.Match(logMsg) && !logEntity.IsMatch {
					logEntity.IsMatch = true
					go func(log *Log) {
						time.AfterFunc(1*time.Second, func() {
							filterLog <- log
						})
					}(logEntity)
				}
			}
		} else {
			break
		}
	}
}

func AllFilter(log chan []byte, filterLog chan *Log, extra ...string) {
	for {
		logMsg, ok := <-log
		if ok {
			filterLog <- NewOriginLog(logMsg)
		} else {
			break
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
		logMsg, ok := <-log
		if ok {
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
				logEntity = NewLog(string(dateTimeStr), string(level), title)
				go func(log *Log) {
					time.AfterFunc(1*time.Second, func() {
						filterLog <- log
					})
				}(logEntity)
			} else if !lineStartPattern.Match(logMsg) && logEntity != nil {
				logEntity.AddLogDetail(logMsg)
			}
		} else {
			break
		}

	}
}
