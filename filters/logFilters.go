package filters

import (
	"fmt"
	"regexp"
	"strings"
)

func ErrorFilter(log chan []byte, filterLog chan *Log, extra ...string) {
	lineStartPattern, _ := regexp.Compile("^[0-9]+-[0-9]+-[0-9]+.*")
	errorLogPattern, _ := regexp.Compile(".*(ERROR).*")
	dateTimePattern, _ := regexp.Compile("[0-9]+-[0-9]+-[0-9]+ [0-9]+:[0-9]+:[0-9]+,[0-9]+")
	titleTimePattern, _ := regexp.Compile(`.*?\[.*?\](.*)`)
	var logEntity *Log
	for {
		logMsg, ok := <-log
		if ok {
			if lineStartPattern.Match(logMsg) && errorLogPattern.Match(logMsg) {
				level := errorLogPattern.FindSubmatch(logMsg)[1]
				dateTimeStr := dateTimePattern.FindSubmatch(logMsg)[0]
				title := strings.Trim(string(titleTimePattern.FindSubmatch(logMsg)[1]), " ")
				logEntity = NewLog(string(dateTimeStr), string(level), title)
			} else if !lineStartPattern.Match(logMsg) && logEntity != nil {
				logEntity.AddErrorLog(logMsg)
			} else if logEntity != nil {
				filterLog <- logEntity
				logEntity = nil
			}
		} else {
			break
		}

	}
}

func InfoFilter(log chan []byte, filterLog chan *Log, extra ...string) {
	lineStartPattern, _ := regexp.Compile("^[0-9]+-[0-9]+-[0-9]+.*")
	infoLogPattern, _ := regexp.Compile(".*(INFO).*")
	dateTimePattern, _ := regexp.Compile("[0-9]+-[0-9]+-[0-9]+ [0-9]+:[0-9]+:[0-9]+,[0-9]+")
	titleTimePattern, _ := regexp.Compile(`.*?\[.*?\](.*)`)
	var logEntity *Log
	for {
		logMsg, ok := <-log
		if ok {
			if lineStartPattern.Match(logMsg) && infoLogPattern.Match(logMsg) {
				level := infoLogPattern.FindSubmatch(logMsg)[1]
				fmt.Println(string(logMsg))
				dateTimeStr := dateTimePattern.FindSubmatch(logMsg)[0]
				title := strings.Trim(string(titleTimePattern.FindSubmatch(logMsg)[1]), " ")
				logEntity = NewLog(string(dateTimeStr), string(level), title)
			} else if !lineStartPattern.Match(logMsg) && logEntity != nil {
				logEntity.AddErrorLog(logMsg)
			} else if logEntity != nil {
				filterLog <- logEntity
				logEntity = nil
			}
		} else {
			break
		}
	}
}

func KeywordFilter(log chan []byte, filterLog chan *Log, extra ...string) {
	lineStartPattern, _ := regexp.Compile("^[0-9]+-[0-9]+-[0-9]+.*")
	levelLogPattern, _ := regexp.Compile(".*(INFO|ERROR).*")
	dateTimePattern, _ := regexp.Compile("[0-9]+-[0-9]+-[0-9]+ [0-9]+:[0-9]+:[0-9]+,[0-9]+")
	titleTimePattern, _ := regexp.Compile(`.*?\[.*?\](.*)`)
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
				title := strings.Trim(string(titleTimePattern.FindSubmatch(logMsg)[1]), " ")
				if logEntity != nil && logEntity.IsMatch {
					filterLog <- logEntity
				}
				logEntity = NewLog(string(dateTimeStr), string(level), title)
				if keywordPattern.Match(logMsg) {
					logEntity.IsMatch = true
				}
			} else if !lineStartPattern.Match(logMsg) && logEntity != nil {
				logEntity.AddErrorLog(logMsg)
				if keywordPattern.Match(logMsg) {
					logEntity.IsMatch = true
				}
			} else if logEntity != nil {
				if logEntity.IsMatch {
					filterLog <- logEntity
				}
				logEntity = nil
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
