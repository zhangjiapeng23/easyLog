/* Namespace menu
support to 1)select namespace, 2)Back to env menu, 3) Exit program
*/

package menu

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type NamespaceMenu struct {
	Status *MenuStatus
	*MenuHelper
	namespaceSearch string
}

func NewNamespaceMenu(menuStatus *MenuStatus) *NamespaceMenu {
	return &NamespaceMenu{menuStatus, &MenuHelper{menuStatus}, ""}
}

func (m *NamespaceMenu) ShowMenu() {
	m.ShowStatus()
	var option string
	searchMatch := false
	namespaceList := m.Status.Client.ListNamespaces().Items
	for index, namespace := range namespaceList {
		// have search keyword, will filter name to show
		if m.namespaceSearch != "" {
			if strings.Contains(namespace.Name, m.namespaceSearch) {
				searchMatch = true
				fmt.Printf("[%02d] %-40s\n", index+1, namespace.Name)
			}
		} else {
			searchMatch = true
			if (index+1)%2 == 0 {
				fmt.Printf("[%02d] %-40s\n", index+1, namespace.Name)
			} else {
				fmt.Printf("[%02d] %-40s\t", index+1, namespace.Name)
			}
			// break line namespace option and other option
			if index == len(namespaceList)-1 && len(namespaceList)%2 != 0 {
				fmt.Println("")
			}
		}
	}
	if !searchMatch {
		fmt.Println("No matching namespace were found")
	}

	fmt.Println("[a] Search")
	fmt.Println("[b] Select Env")
	fmt.Println("[c] Exit")
	fmt.Print("Please select namespace: ")
	fmt.Scan(&option)
	if m.isDigit(option) {
		optionInt64, _ := strconv.ParseInt(option, 10, 8)
		optionInt := int(optionInt64) - 1
		if optionInt >= 0 && optionInt < len(m.Status.Client.ListNamespaces().Items) {
			m.SelectNameSpace(optionInt)
		} else {
			fmt.Println("Paramter parse error.")
		}
	} else {
		switch option {
		case "a":
			fmt.Print("Please input namespace: ")
			fmt.Scan(&m.namespaceSearch)
			m.ShowMenu()
		case "b":
			m.SelectEnv(-1)
		case "c":
			m.Close()
		default:
			fmt.Println("Paramter parse error")
		}
	}
}

func (m *NamespaceMenu) SelectEnv(option int) {
	CurrentMenu = NewEnvMenu(m.Status)
}

func (m *NamespaceMenu) SelectNameSpace(option int) {
	namespace := m.Status.Client.ListNamespaces().Items[option]
	m.Status.Namespace = namespace.Name
	m.Status.NamespaceObj = &namespace
	// clear app, log model, log filter
	m.Status.App = ""
	m.Status.Command = ""
	m.Status.LogFilter = ""
	CurrentMenu = NewAppMenu(m.Status)
}

func (m *NamespaceMenu) SelectApp(option int) {
	fmt.Fprintf(os.Stderr, "please select namespace first")
}

func (m *NamespaceMenu) SelectCommand(command string) {
	fmt.Fprintf(os.Stderr, "please select namespace first")
}

func (m *NamespaceMenu) SelectLogFilter(option int) {
	fmt.Fprintf(os.Stderr, "please select namespace first")
}

func (m *NamespaceMenu) Close() {
	fmt.Println("Exiting...")
	os.Exit(0)
}
