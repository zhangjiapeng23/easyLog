package menu

import (
	"fmt"
	"os"
	"strconv"
)

type PodExecMenu struct {
	Status *MenuStatus
	*MenuHelper
}

func NewPodExecMenu(menuStatus *MenuStatus) *PodExecMenu {
	return &PodExecMenu{menuStatus, &MenuHelper{menuStatus}}
}

func (m *PodExecMenu) ShowMenu() {
	m.ShowStatus()
	podName := ""
	podList := m.Status.Client.ListPodsForApp(m.Status.Namespace, m.Status.App)
	if len(podList.Items) <= 0 {
		fmt.Println("The number of pods is 0")
		CurrentMenu = NewCommandMenu(m.Status)
	} else {
		option := ""
		for index, pod := range podList.Items {

			fmt.Printf("[%d] %s\n", index+1, pod.Name)
		}
		fmt.Println("[a] Select env")
		fmt.Println("[b] Select namespace")
		fmt.Println("[c] Select app")
		fmt.Println("[d] Select command")
		fmt.Println("[e] Exit")
		fmt.Print("Please select pod (input 'exit' quit container): ")
		fmt.Scan(&option)
		if m.isDigit(option) {
			optionInt64, _ := strconv.ParseInt(option, 10, 32)
			optionInt := int(optionInt64) - 1
			if optionInt >= 0 && optionInt < len(podList.Items) {
				podName = podList.Items[optionInt].Name
				m.SelectPod(podName)
			} else {
				fmt.Println("Paramter parse error.")
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
				m.SelectCommand("")
			case "e":
				m.Close()
			default:
				fmt.Println("Paramter parse error")
			}
		}
	}
}

func (m *PodExecMenu) SelectEnv(option int) {
	CurrentMenu = NewEnvMenu(m.Status)
}

func (m *PodExecMenu) SelectNameSpace(option int) {
	CurrentMenu = NewNamespaceMenu(m.Status)
}

func (m *PodExecMenu) SelectApp(option int) {
	CurrentMenu = NewAppMenu(m.Status)
}

func (m *PodExecMenu) SelectCommand(command string) {
	CurrentMenu = NewCommandMenu(m.Status)
}

func (m *PodExecMenu) SelectLogFilter(option int) {
	fmt.Fprintf(os.Stderr, "please select log model first")
}

func (m *PodExecMenu) Close() {
	fmt.Println("Exiting...")
	os.Exit(0)
}

func (m *PodExecMenu) SelectPod(podName string) {
	m.Status.PodName = podName
	m.Status.Client.ExecPod(m.Status.Namespace, podName)
}
