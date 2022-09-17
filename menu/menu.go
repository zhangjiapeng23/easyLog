/*
 define Menu status struct and Menu inferface
*/

package menu

import (
	"easyLog/filters"
	"easyLog/k8s"
	"easyLog/util"
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

var (
	CurrentMenu Menu
	printGreen  = color.New(color.FgGreen)
	printCyan   = color.New(color.FgCyan)
	printRed    = color.New(color.FgHiRed)
	printBulue = color.New(color.FgBlue)
)

type MenuStatus struct {
	Env          string
	Namespace    string
	App          string
	Command      string
	LogFilter    string
	PodName      string
	Client       *k8s.Client
	NamespaceObj *v1.Namespace
	AppObj       *appsv1.Deployment
	LogFilterObj func(log chan []byte, filterLog chan *filters.Log, extra ...string)
}

type Menu interface {
	ShowMenu()
	SelectEnv(option int)
	SelectNameSpace(optioin int)
	SelectApp(option int)
	SelectCommand(command string)
	SelectLogFilter(option int)
	Close()
}

type MenuHelper struct {
	Status *MenuStatus
}

func (m *MenuHelper) ShowStatus() {
	flag := false
	if m.Status.Env != "" {
		printCyan.Printf("Env: %s  ", m.Status.Env)
		flag = true
	}
	if m.Status.Namespace != "" {
		printCyan.Printf("Namespace: %s  ", m.Status.Namespace)
		flag = true
	}
	if m.Status.App != "" {
		printCyan.Printf("App: %s  ", m.Status.App)
		flag = true
	}
	if m.Status.Command != "" {
		printCyan.Printf("Cmd: %s  ", m.Status.Command)
		flag = true
	}
	if m.Status.LogFilter != "" {
		printCyan.Printf("Filter: %s  ", m.Status.LogFilter)
		flag = true
	}
	if m.Status.PodName != "" {
		printCyan.Printf("Pod: %s", m.Status.PodName)
		flag = true
	}
	if flag {
		fmt.Println("")
		util.PrintSplitLine("-")
	}
}

func (m *MenuHelper) isDigit(s string) bool {
	_, err := strconv.ParseInt(s, 10, 32)
	return err == nil
}

func (m *MenuHelper) colorKeyWord(s string, keyword string) string {
	red := color.New(color.FgRed).SprintFunc()
	otherWords := strings.Split(s, keyword)
	return strings.Join(otherWords, fmt.Sprint(red(keyword)))
}

func init() {
	CurrentMenu = NewEnvMenu(
		&MenuStatus{
			Env:          "",
			Namespace:    "",
			App:          "",
			Command:      "",
			LogFilter:    "",
			PodName:      "",
			Client:       nil,
			NamespaceObj: nil,
			AppObj:       nil,
			LogFilterObj: nil,
		},
	)
}
