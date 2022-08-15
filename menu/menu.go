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
	LogModel string
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
	SelectLogModel(logModel string)
	SelectLogFilter(option int)
	Close()
}

type MenuHelper struct {
	Status *MenuStatus
}

func (menu *MenuHelper) ShowStatus() {
	flag := false
	if menu.Status.Env != "" {
		fmt.Printf("Env: %s\t", menu.Status.Env)
		flag = true
	}
	if menu.Status.Namespace != "" {
		fmt.Printf("Namespace: %s\t", menu.Status.Namespace)
		flag = true
	}
	if menu.Status.App != "" {
		fmt.Printf("App: %s\t", menu.Status.App)
		flag = true
	}
	if menu.Status.LogModel != "" {
		fmt.Printf("Log Model: %s\t", menu.Status.LogModel)
		flag = true
	}
	if menu.Status.LogFilter != "" {
		fmt.Printf("Filter: %s\t", menu.Status.LogFilter)
		flag = true
	}
	if flag {
		fmt.Println("")
		util.PrintSplitLine("-")
	}
}

func (menu *MenuHelper) isDigit(s string) bool {
	_, err := strconv.ParseInt(s, 10, 8)
	return err == nil
}

func init() {
	CurrentMenu = NewEnvMenu(
		&MenuStatus{
			Env:       "",
			Namespace: "",
			App:       "",
			LogModel:  "",
			LogFilter: "",
			Client: nil,
			NamespaceObj: nil,
			AppObj: nil,
			LogFilterObj: nil,
		},
	)
}