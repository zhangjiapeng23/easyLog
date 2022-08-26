/* Log model menu
Support to 1) select output log model, 2) back to env, namespace, app menu, 3) exit program
*/

package menu

import (
	"easyLog/util"
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
)

type CommandMenu struct {
	Status *MenuStatus
	*MenuHelper
}

func NewCommandMenu(menuStatus *MenuStatus) *CommandMenu {
	return &CommandMenu{menuStatus, &MenuHelper{menuStatus}}
}

func (m *CommandMenu) ShowMenu() {
	m.ShowStatus()
	var option string
	fmt.Println("[1] Print log")
	fmt.Println("[2] Follow log")
	fmt.Println("[3] Fetch pod info")
	fmt.Println("[4] Execute shell")
	fmt.Println("[a] Select env")
	fmt.Println("[b] Select namespace")
	fmt.Println("[c] Select app")
	fmt.Println("[d] Exit")
	fmt.Print("Please select command: ")
	fmt.Scan(&option)
	if m.isDigit(option) {
		switch option {
		case "1":
			m.SelectCommand("Print Log")
		case "2":
			m.SelectCommand("Follow Log")
		case "3":
			m.SelectCommand("Fetch Pod Info")
		case "4":
			m.SelectCommand("Execute Shell")
		default:
			fmt.Println("Paramter parse error")
		}
	} else {
		switch option {
		case "a":
			m.SelectEnv(-1)
		case "b":
			m.SelectNameSpace(-1)
		case "c":
			m.SelectApp(-1)
		case "d":
			m.Close()
		default:
			fmt.Println("Paramter parse error")
		}
	}
}

func (m *CommandMenu) SelectEnv(option int) {
	CurrentMenu = NewEnvMenu(m.Status)
}

func (m *CommandMenu) SelectNameSpace(option int) {
	CurrentMenu = NewNamespaceMenu(m.Status)
}

func (m *CommandMenu) SelectApp(option int) {
	CurrentMenu = NewAppMenu(m.Status)
}

func (m *CommandMenu) SelectCommand(command string) {
	m.Status.Command = command
	// clear log filter and pod name
	m.Status.LogFilter = ""
	m.Status.PodName = ""
	if command == "Fetch Pod Info" {
		m.FetchPodsInfo()
		CurrentMenu = NewCommandMenu(m.Status)
	} else if command == "Execute Shell" {
		CurrentMenu = NewPodExecMenu(m.Status)
	} else {
		CurrentMenu = NewLogFilterMenu(m.Status)
	}
}

func (m *CommandMenu) SelectLogFilter(option int) {
	fmt.Fprintf(os.Stderr, "please select log model first")
}

func (m *CommandMenu) Close() {
	fmt.Println("Exiting...")
	os.Exit(0)
}

func (m *CommandMenu) FetchPodsInfo() {
	podList := m.Status.Client.ListPodsForApp(m.Status.Namespace, m.Status.App)
	for _, pod := range podList.Items {
		util.PrintSplitLine("-")
		printPodInfo(pod)
	}
}

func printPodInfo(pod v1.Pod) {
	name := pod.ObjectMeta.Name
	ip := pod.Status.PodIP
	port := pod.Annotations["prometheus.io/port"]
	status := ""
	ready := false
	var restart int32 = 0
	if len(pod.Status.Conditions) > 0 {
		status = string(pod.Status.Conditions[0].Status)
	}
	if len(pod.Status.ContainerStatuses) > 0 {
		ready = pod.Status.ContainerStatuses[0].Ready
	}
	if len(pod.Status.ContainerStatuses) > 0 {
		restart = pod.Status.ContainerStatuses[0].RestartCount
	}
	fmt.Printf("Pod Name: %s\n", name)
	fmt.Printf("Pod IP: %s\n", ip)
	fmt.Printf("Pod Port: %s\n", port)
	fmt.Printf("Status: %v\n", status)
	fmt.Printf("Ready: %v\n", ready)
	fmt.Printf("Restart Count: %d\n", restart)
}
