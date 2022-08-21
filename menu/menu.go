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

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

var (
	CurrentMenu Menu
)

type MenuStatus struct {
	Env string
	Namespace string
	App string
	Command string
	LogFilter string
	Client *k8s.Client
	NamespaceObj *v1.Namespace
	AppObj *appsv1.Deployment
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
		fmt.Printf("Env: %s\t", m.Status.Env)
		flag = true
	}
	if m.Status.Namespace != "" {
		fmt.Printf("Namespace: %s\t", m.Status.Namespace)
		flag = true
	}
	if m.Status.App != "" {
		fmt.Printf("App: %s\t", m.Status.App)
		flag = true
	}
	if m.Status.Command != "" {
		fmt.Printf("Command: %s\t", m.Status.Command)
		flag = true
	}
	if m.Status.LogFilter != "" {
		fmt.Printf("Filter: %s\t", m.Status.LogFilter)
		flag = true
	}
	if flag {
		fmt.Println("")
		util.PrintSplitLine("-")
	}
}

func (m *MenuHelper) isDigit(s string) bool {
	_, err := strconv.ParseInt(s, 10, 8)
	return err == nil
}

func init() {
	CurrentMenu = NewEnvMenu(
		&MenuStatus{
			Env:       "",
			Namespace: "",
			App:       "",
			Command:  "",
			LogFilter: "",
			Client: nil,
			NamespaceObj: nil,
			AppObj: nil,
			LogFilterObj: nil,
		},
	)
}