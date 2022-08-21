/* Env menu
Support to 1) select env, 2) exeit program
*/

package menu

import (
	"easyLog/k8s"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"k8s.io/client-go/util/homedir"
)

type EnvMenu struct {
	Status *MenuStatus
	*MenuHelper
	FilePathMap map[string]string
	EnvRec      []string
}

func NewEnvMenu(menuStatus *MenuStatus) *EnvMenu {
	return &EnvMenu{menuStatus, &MenuHelper{menuStatus}, make(map[string]string), make([]string, 0)}
}

func (m *EnvMenu) ShowMenu() {
	m.ShowStatus()
	var option string
	m.FilePathMap = searchK8SConifg()

	for key := range m.FilePathMap {
		m.EnvRec = append(m.EnvRec, key)
	}
	sort.Slice(m.EnvRec, func(i, j int) bool {
		if m.EnvRec[i] < m.EnvRec[j] {
			return true
		} else {
			return false
		}
	})
	index := 1
	for _, env := range m.EnvRec {
		fmt.Printf("[%d] %s\n", index, env)
		m.EnvRec = append(m.EnvRec, env)
		index++
	}
	fmt.Println("[a] Exit")
	fmt.Print("Please select env: ")
	fmt.Scan(&option)
	if m.isDigit(option) {
		optionInt64, _ := strconv.ParseInt(option, 10, 8)
		optionInt := int(optionInt64) - 1
		if optionInt >= 0 && optionInt < len(m.EnvRec) {
			m.SelectEnv(optionInt)
		} else {
			m.EnvRec = make([]string, 0)
			fmt.Println("Paramter parse error.")
		}
	} else {
		switch option {
		case "a":
			m.Close()
		default:
			m.EnvRec = make([]string, 0)
			fmt.Println("Paramter parse error.")
		}
	}
}

func (m *EnvMenu) SelectEnv(option int) {
	m.Status.Env = m.EnvRec[option]
	m.Status.Client = k8s.NewClient(m.FilePathMap[m.Status.Env])
	// clear namespace, app, log model, log filter
	m.Status.Namespace = ""
	m.Status.App = ""
	m.Status.Command = ""
	m.Status.LogFilter = ""
	CurrentMenu = NewNamespaceMenu(m.Status)
}

func (m *EnvMenu) SelectNameSpace(option int) {
	fmt.Fprintf(os.Stderr, "please select env first")
}

func (m *EnvMenu) SelectApp(option int) {
	fmt.Fprintf(os.Stderr, "please select env first")
}

func (m *EnvMenu) SelectCommand(command string) {
	fmt.Fprintf(os.Stderr, "please select env first")
}

func (m *EnvMenu) SelectLogFilter(option int) {
	fmt.Fprintf(os.Stderr, "please select env first")
}

func (m *EnvMenu) Close() {
	fmt.Println("Exiting...")
	os.Exit(0)
}

// This function to find k8s config, strategy is first to find local kube homedir is whether exits
// config. if find, will use this to show env menu. Otherwise, will use dir /k8s/conf default config.
func searchK8SConifg() (filePath map[string]string) {
	var kubeconfigDir string
	filePath = make(map[string]string)
	if home := homedir.HomeDir(); home != "" {
		// first search .kube dir is whether exits
		kubeconfigDir = filepath.Join(home, ".kube")
		exits, err := PathExits(kubeconfigDir)
		if err != nil {
			panic(fmt.Sprintf("load k8s config error: %v", err.Error()))
		}
		if !exits {
			file, _ := exec.LookPath(os.Args[0])
			path, _ := filepath.Abs(file)
			index := strings.LastIndex(path, string(os.PathSeparator))
			path = path[:index]
			kubeconfigDir = filepath.Join(path, "k8s", "conf")
		} else {
			// check is whether empty dir
			rd, err := ioutil.ReadDir(kubeconfigDir)
			if err != nil {
				panic(err.Error())
			}
			if len(filterFile(rd)) == 0 {
				file, _ := exec.LookPath(os.Args[0])
				path, _ := filepath.Abs(file)
				index := strings.LastIndex(path, string(os.PathSeparator))
				path = path[:index]
				kubeconfigDir = filepath.Join(path, "k8s", "conf")
			}
		}
	} else {
		file, _ := exec.LookPath(os.Args[0])
		path, _ := filepath.Abs(file)
		index := strings.LastIndex(path, string(os.PathSeparator))
		path = path[:index]
		kubeconfigDir = filepath.Join(path, "k8s", "conf")
	}
	rd, err := ioutil.ReadDir(kubeconfigDir)
	if err != nil {
		panic(err.Error())
	}
	for _, file := range filterFile(rd) {
		filePath[file.Name()] = filepath.Join(kubeconfigDir, file.Name())
	}
	return
}

func PathExits(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func filterFile(rd []fs.FileInfo) (fileList []fs.FileInfo) {
	for _, fi := range rd {
		if !fi.IsDir() {
			fileList = append(fileList, fi)
		}
	}
	return
}
