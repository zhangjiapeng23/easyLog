/* App menu
support to 1) select app, 2) back to env or namespace menu, 3) exit program
*/

package menu

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type AppMenu struct {
	Status *MenuStatus
	*MenuHelper
	appSearch string
}

func NewAppMenu(menuStatus *MenuStatus) *AppMenu {
	return &AppMenu{menuStatus, &MenuHelper{menuStatus}, ""}
}

func (m *AppMenu) ShowMenu() {
	m.ShowStatus()
	var option string
	appList := m.Status.Client.ListAppsForNamespace(m.Status.Namespace).Items
	searchMatch := false
	for index, app := range appList {
		// have search keyword, filter app name to show.
		if m.appSearch != "" {
			if strings.Contains(app.Name, m.appSearch) {
				searchMatch = true
				fmt.Printf("[%02d] %-40s\n", index+1, app.Name)
			}
		} else {
			searchMatch = true
			if (index+1)%2 == 0 {
				fmt.Printf("[%02d] %-40s\n", index+1, app.Name)
			} else {
				fmt.Printf("[%02d] %-40s\t", index+1, app.Name)
			}
			// break line app option and other option
			if index == len(appList)-1 && len(appList)%2 != 0 {
				fmt.Println("")
			}
		}
	}
	if !searchMatch {
		fmt.Println("No matching apps were found")
	}

	fmt.Println("[a] Search")
	fmt.Println("[b] Select env")
	fmt.Println("[c] Select namespace")
	fmt.Println("[d] Exit")
	fmt.Print("Please select app: ")
	fmt.Scan(&option)
	if m.isDigit(option) {
		optionInt64, _ := strconv.ParseInt(option, 10, 8)
		optionInt := int(optionInt64) - 1
		if optionInt >= 0 && optionInt < len(m.Status.Client.ListAppsForNamespace(m.Status.Namespace).Items) {
			m.SelectApp(optionInt)
		} else {
			fmt.Println("Paramter parse error.")
		}
	} else {
		switch option {
		case "a":
			fmt.Print("Please input app: ")
			fmt.Scan(&m.appSearch)
			m.ShowMenu()
		case "b":
			m.SelectEnv(-1)
		case "c":
			m.SelectNameSpace(-1)
		case "d":
			m.Close()
		default:
			fmt.Println("Paramter parse error")
		}
	}
}

func (m *AppMenu) SelectEnv(option int) {
	CurrentMenu = NewEnvMenu(m.Status)
}

func (m *AppMenu) SelectNameSpace(option int) {
	CurrentMenu = NewNamespaceMenu(m.Status)
}

func (m *AppMenu) SelectApp(option int) {
	appList := m.Status.Client.ListAppsForNamespace(m.Status.Namespace).Items
	m.Status.AppObj = &appList[option]
	m.Status.App = appList[option].Name
	// clear log model, log filter
	m.Status.Command = ""
	m.Status.LogFilter = ""
	CurrentMenu = NewLogModelMenu(m.Status)
}

func (m *AppMenu) SelectCommand(command string) {
	fmt.Fprintf(os.Stderr, "please select app first")
}

func (m *AppMenu) SelectLogFilter(option int) {
	fmt.Fprintf(os.Stderr, "please select app first")
}

func (m *AppMenu) Close() {
	fmt.Println("Exiting...")
	os.Exit(0)
}
