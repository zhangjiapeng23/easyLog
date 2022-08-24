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
		"Warn Filter":    filters.WarnFilter,
		"Info Filter":    filters.InfoFilter,
		"Debug Filter":   filters.DebugFilter,
		"Keyword Filter": filters.KeywordFilter,
		"All Filter":     filters.AllFilter,
	}, make([]string, 0)}
}

func (m *LogFilterMenu) ShowMenu() {
	m.ShowStatus()
	var option string
	for key := range m.filterRegister {
		m.filterRec = append(m.filterRec, key)
	}
	sort.Slice(m.filterRec, func(i, j int) bool {
		return m.filterRec[i] < m.filterRec[j]
	})

	for index, filter := range m.filterRec {
		fmt.Printf("[%d] %s\n", index+1, filter)
	}

	fmt.Println("[a] Select env")
	fmt.Println("[b] Select namespace")
	fmt.Println("[c] Select app")
	fmt.Println("[d] Select command")
	fmt.Println("[e] Exit")
	fmt.Print("Please select log filter (Ctrl+C quit): ")
	fmt.Scan(&option)
	if m.isDigit(option) {
		optionInt64, _ := strconv.ParseInt(option, 10, 32)
		optionInt := int(optionInt64) - 1
		if optionInt >= 0 && optionInt < len(m.filterRec) {
			m.SelectLogFilter(optionInt)
		} else {
			m.filterRec = make([]string, 0)
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
			m.filterRec = make([]string, 0)
			fmt.Println("Parameter parse error")
		}
	}
}

func (m *LogFilterMenu) SelectEnv(option int) {
	CurrentMenu = NewEnvMenu(m.Status)
}

func (m *LogFilterMenu) SelectNameSpace(option int) {
	CurrentMenu = NewNamespaceMenu(m.Status)
}

func (m *LogFilterMenu) SelectApp(option int) {
	CurrentMenu = NewAppMenu(m.Status)
}

func (m *LogFilterMenu) SelectCommand(command string) {
	CurrentMenu = NewCommandMenu(m.Status)
}

func (m *LogFilterMenu) SelectLogFilter(option int) {
	key := m.filterRec[option]
	m.Status.LogFilter = key
	m.Status.LogFilterObj = m.filterRegister[key]
	keyword := ""
	if strings.Contains(key, "Keyword") {
		fmt.Print("Please input keyword: ")
		fmt.Scan(&keyword)
	}
	m.PringLog(keyword)
	CurrentMenu = NewLogFilterMenu(m.Status)
}

func (m *LogFilterMenu) Close() {
	fmt.Println("Exiting...")
	os.Exit(0)
}

func (m *LogFilterMenu) PringLog(extra ...string) {
	podList := m.Status.Client.ListPodsForApp(m.Status.Namespace, m.Status.App)

	if m.Status.Command == "Print Log" {
		m.Status.Client.PrintLogForPods(m.Status.Namespace, podList, m.Status.LogFilterObj, extra...)
	} else if m.Status.Command == "Follow Log" {
		m.Status.Client.FollowLogForPods(m.Status.Namespace, podList, m.Status.LogFilterObj, extra...)
	} else {
		fmt.Println("Log model select error!")
	}

}
