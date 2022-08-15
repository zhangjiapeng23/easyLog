/* Log model menu
 Support to 1) select output log model, 2) back to env, namespace, app menu, 3) exit program
 */

package menu

import (
	"fmt"
	"os"
)

type LogModelMenu struct {
	Status *MenuStatus
	*MenuHelper
}

func NewLogModelMenu(menuStatus *MenuStatus) *LogModelMenu{
	return &LogModelMenu{menuStatus, &MenuHelper{menuStatus}}
}

func (Menu *LogModelMenu) ShowMenu() {
	Menu.ShowStatus()
	var option string
	fmt.Println("[1] Print log")
	fmt.Println("[2] Follow log")
	fmt.Println("[a] Select env")
	fmt.Println("[b] Select namespace")
	fmt.Println("[c] Select app")
	fmt.Println("[d] Exit")
	fmt.Print("Please select log model: ")
	fmt.Scan(&option)
	if Menu.isDigit(option) {
		if option == "1" {
			Menu.SelectLogModel("Print Log")
		} else if option == "2" {
			Menu.SelectLogModel("Follow Log")
		} else {
			fmt.Println("Paramter parse error")
		}
	} else if option == "a" {
		Menu.SelectEnv(-1)
	} else if option == "b" {
		Menu.SelectNameSpace(-1)
	} else if option == "c" {
		Menu.SelectApp(-1)
	} else if option == "d" {
		Menu.Close()
	} else {
		fmt.Println("Paramter parse error")
	}
}

func (Menu *LogModelMenu) SelectEnv(option int) {
	CurrentMenu = NewEnvMenu(Menu.Status)
}

func (Menu *LogModelMenu) SelectNameSpace(option int) {
	CurrentMenu = NewNamespaceMenu(Menu.Status)
}

func (Menu *LogModelMenu) SelectApp(option int) {
	CurrentMenu = NewAppMenu(Menu.Status)
}

func (Menu *LogModelMenu) SelectLogModel(logModel string) {
	Menu.Status.LogModel = logModel
	CurrentMenu = NewLogFilterMenu(Menu.Status)
}

func (Menu *LogModelMenu) SelectLogFilter(option int) {
	fmt.Fprintf(os.Stderr, "please select log model first")
}

func (Menu *LogModelMenu) Close() {
	fmt.Println("Exiting...")
	os.Exit(0)
}
