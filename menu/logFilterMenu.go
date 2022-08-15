/* Log filter menu
Support to 1) select print log filter, 2) back to env, namespace, app, log model menu, 3) exit program
*/

package menu

import (
	"easyLog/filters"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type LogFilterMenu struct {
	Status *MenuStatus
	*MenuHelper
	filterRegister map[string]func(log chan []byte, filterLog chan *filters.Log, extra ...string)
	filterRec      []string
}

func NewLogFilterMenu(menuStatus *MenuStatus) *LogFilterMenu {
	return &LogFilterMenu{menuStatus, &MenuHelper{menuStatus}, map[string]func(log chan []byte,
		filterLog chan *filters.Log, extra ...string){
		"Error Filter":   filters.ErrorFilter,
		"Info Filter":    filters.InfoFilter,
		"Keyword Filter": filters.KeywordFilter,
		"All Filter":     filters.AllFilter,
	}, make([]string, 0)}
}

func (Menu *LogFilterMenu) ShowMenu() {
	Menu.ShowStatus()
	var option string
	for key := range Menu.filterRegister {
		Menu.filterRec = append(Menu.filterRec, key)
	}
	sort.Slice(Menu.filterRec, func(i, j int) bool {
		return Menu.filterRec[i] < Menu.filterRec[j]
	})

	for index, filter := range Menu.filterRec {
		fmt.Printf("[%d] %s\n", index+1, filter)
	}

	fmt.Println("[a] Select env")
	fmt.Println("[b] Select namespace")
	fmt.Println("[c] Select app")
	fmt.Println("[d] Select log model")
	fmt.Println("[e] Exit")
	fmt.Print("Please select log filter: ")
	fmt.Scan(&option)
	if Menu.isDigit(option) {
		optionInt, _ := strconv.ParseInt(option, 10, 8)
		Menu.SelectLogFilter(int(optionInt) - 1)
	} else if option == "a" {
		Menu.SelectEnv(-1)
	} else if option == "b" {
		Menu.SelectNameSpace(-1)
	} else if option == "c" {
		Menu.SelectApp(-1)
	} else if option == "d" {
		Menu.SelectLogModel("")
	} else if option == "e" {
		Menu.Close()
	} else {
		fmt.Println("Parameter parse error")
	}

}

func (Menu *LogFilterMenu) SelectEnv(option int) {
	CurrentMenu = NewEnvMenu(Menu.Status)
}

func (Menu *LogFilterMenu) SelectNameSpace(option int) {
	CurrentMenu = NewNamespaceMenu(Menu.Status)
}

func (Menu *LogFilterMenu) SelectApp(option int) {
	CurrentMenu = NewAppMenu(Menu.Status)
}

func (Menu *LogFilterMenu) SelectLogModel(logModel string) {
	CurrentMenu = NewLogModelMenu(Menu.Status)
}

func (Menu *LogFilterMenu) SelectLogFilter(option int) {
	// Menu.Status.LogFilter = logFilter
	key := Menu.filterRec[option]
	Menu.Status.LogFilter = key
	Menu.Status.LogFilterObj = Menu.filterRegister[key]
	keyword := ""

	if strings.Contains(key, "Keyword") {
		fmt.Print("Please input keyword: ")
		fmt.Scan(&keyword)
	}

	Menu.PringLog(keyword)
	CurrentMenu = NewLogFilterMenu(Menu.Status)
}

func (Menu *LogFilterMenu) Close() {
	fmt.Println("Exiting...")
	os.Exit(0)
}

func (Menu *LogFilterMenu) PringLog(extra ...string) {
	podList := Menu.Status.Client.ListPodsForApp(Menu.Status.Namespace, Menu.Status.App)

	if Menu.Status.LogModel == "Print Log" {
		Menu.Status.Client.PrintLogForPods(Menu.Status.Namespace, podList, Menu.Status.LogFilterObj, extra...)
	} else if Menu.Status.LogModel == "Follow Log" {
		Menu.Status.Client.FollowLogForPods(Menu.Status.Namespace, podList, Menu.Status.LogFilterObj, extra...)
	} else {
		fmt.Println("Log model select error!")
	}

}
