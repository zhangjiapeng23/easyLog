/* App menu
support to 1) select app, 2) back to env or namespace menu, 3) exit program
*/

package menu

import (
	"fmt"
	"os"
	"strconv"
)

type AppMenu struct {
	Status *MenuStatus
	*MenuHelper
}

func NewAppMenu(menuStatus *MenuStatus) *AppMenu {
	return &AppMenu{menuStatus, &MenuHelper{menuStatus}}
}

func (Menu *AppMenu) ShowMenu() {
	Menu.ShowStatus()
	var option string
	appList := Menu.Status.Client.ListAppsForNamespace(Menu.Status.Namespace).Items
	for index, app := range appList {
		if (index+1) % 2 == 0 {
			fmt.Printf("[%02d] %-40s\n", index+1, app.Name)
		} else {
			fmt.Printf("[%02d] %-40s\t", index+1, app.Name)
		}
	}

	if len(appList) % 2 != 0 {
		fmt.Println("")
	}

	fmt.Println("[a] Select env")
	fmt.Println("[b] Select namespace")
	fmt.Println("[c] Exit")
	fmt.Print("Please select app: ")
	fmt.Scan(&option)
	if Menu.isDigit(option) {
		optionInt, _ := strconv.ParseInt(option, 10, 8)
		Menu.SelectApp(int(optionInt)-1)
	} else if option == "a" {
		Menu.SelectEnv(-1)
	} else if option == "b" {
		Menu.SelectNameSpace(-1)
	} else if option == "c" {
		Menu.Close()
	} else {
		fmt.Println("Paramter parse error")
	}

}

func (Menu *AppMenu) SelectEnv(option int) {
	CurrentMenu = NewEnvMenu(Menu.Status)
}

func (Menu *AppMenu) SelectNameSpace(option int) {
	CurrentMenu = NewNamespaceMenu(Menu.Status)
}

func (Menu *AppMenu) SelectApp(option int) {
	appList := Menu.Status.Client.ListAppsForNamespace(Menu.Status.Namespace).Items
	Menu.Status.AppObj = &appList[option]
	Menu.Status.App = appList[option].Name
	CurrentMenu = NewLogModelMenu(Menu.Status)
}

func (Menu *AppMenu) SelectLogModel(logModel string) {
	fmt.Fprintf(os.Stderr, "please select app first")
}

func (Menu *AppMenu) SelectLogFilter(option int) {
	fmt.Fprintf(os.Stderr, "please select app first")
}

func (Menu *AppMenu) Close() {
	fmt.Println("Exiting...")
	os.Exit(0)
}