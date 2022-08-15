/* Namespace menu
support to 1)select namespace, 2)Back to env menu, 3) Exit program
*/

package menu

import (
	"fmt"
	"os"
	"strconv"
)

type NamespaceMenu struct {
	Status *MenuStatus
	*MenuHelper
}

func NewNamespaceMenu(menuStatus *MenuStatus) *NamespaceMenu{
	return &NamespaceMenu{menuStatus, &MenuHelper{menuStatus}}
}

func (Menu *NamespaceMenu) ShowMenu() {
	Menu.ShowStatus()
	var option string
	namespaceList := Menu.Status.Client.ListNamespaces().Items
	for index, namespace := range namespaceList {
		if (index+1) % 2 == 0 {
			fmt.Printf("[%02d] %-40s\n", index+1, namespace.Name)
		} else {
			fmt.Printf("[%02d] %-40s\t", index+1, namespace.Name)
		}
	}

	if len(namespaceList) % 2 != 0 {
		fmt.Println("")
	}

	fmt.Println("[a] Select Env")
	fmt.Println("[b] Exit")
	fmt.Print("Please select namespace: ")
	fmt.Scan(&option)
	if Menu.isDigit(option) {
		optionInt, _ := strconv.ParseInt(option, 10, 8)
		Menu.SelectNameSpace(int(optionInt)-1)
	
	// if option == "1" || option == "2" {
	// 	if option == "1" {
	// 		Menu.SelectNameSpace("office")
	// 	} else {
	// 		Menu.SelectNameSpace("web")
	// 	}
	} else if option == "a" {
		Menu.SelectEnv(-1)
	} else if option == "b" {
		Menu.Close()
	} else {
		fmt.Println("Paramter parse error")
	}
}

func (Menu *NamespaceMenu) SelectEnv(option int) {
	CurrentMenu = NewEnvMenu(Menu.Status)
}

func (Menu *NamespaceMenu) SelectNameSpace(option int) {
	namespace := Menu.Status.Client.ListNamespaces().Items[option]
	Menu.Status.Namespace = namespace.Name
	Menu.Status.NamespaceObj = &namespace
	CurrentMenu = NewAppMenu(Menu.Status)
}

func (Menu *NamespaceMenu) SelectApp(option int) {
	fmt.Fprintf(os.Stderr, "please select namespace first")
}

func (Menu *NamespaceMenu) SelectLogModel(logModel string) {
	fmt.Fprintf(os.Stderr, "please select namespace first")
}

func (Menu *NamespaceMenu) SelectLogFilter(option int) {
	fmt.Fprintf(os.Stderr, "please select namespace first")
}

func (Menu *NamespaceMenu) Close() {
	fmt.Println("Exiting...")
	os.Exit(0)
}