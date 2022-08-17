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

func (Menu *EnvMenu) ShowMenu() {
	Menu.ShowStatus()
	var option string
	// fmt.Println("[1] pre")
	// fmt.Println("[2] uat")
	Menu.FilePathMap = searchK8SConifg()
	index := 1
	for key := range Menu.FilePathMap {
		fmt.Printf("[%d] %s\n", index, key)
		Menu.EnvRec = append(Menu.EnvRec, key)
		index++
	}
	fmt.Println("[a] Exit")
	fmt.Print("Please select env: ")
	fmt.Scan(&option)
	if Menu.isDigit(option) {
		optionInt, _ := strconv.ParseInt(option, 10, 8)
		Menu.SelectEnv(int(optionInt) - 1)
	} else if option == "a" {
		Menu.Close()
	} else {
		fmt.Println("Paramter parse error.")
	}
}

func (Menu *EnvMenu) SelectEnv(option int) {
	Menu.Status.Env = Menu.EnvRec[option]
	Menu.Status.Client = k8s.NewClient(Menu.FilePathMap[Menu.Status.Env])
	Menu.Status.Namespace = ""
	Menu.Status.App = ""
	CurrentMenu = NewNamespaceMenu(Menu.Status)
}

func (Menu *EnvMenu) SelectNameSpace(option int) {
	fmt.Fprintf(os.Stderr, "please select env first")
}

func (Menu *EnvMenu) SelectApp(option int) {
	fmt.Fprintf(os.Stderr, "please select env first")
}

func (Menu *EnvMenu) SelectLogModel(logModel string) {
	fmt.Fprintf(os.Stderr, "please select env first")
}

func (Menu *EnvMenu) SelectLogFilter(option int) {
	fmt.Fprintf(os.Stderr, "please select env first")
}

func (Menu *EnvMenu) Close() {
	fmt.Println("Exiting...")
	os.Exit(0)
}

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
			// get proj root dir
			// projDir, _ := os.Getwd()
			// kubeconfigDir = filepath.Join(projDir, "k8s", "conf")
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
				// projDir, _ := os.Getwd()
				file, _ := exec.LookPath(os.Args[0])
				path, _ := filepath.Abs(file)
				index := strings.LastIndex(path, string(os.PathSeparator))
				path = path[:index]
				kubeconfigDir = filepath.Join(path, "k8s", "conf")
			}
		}
	} else {
		// projDir, _ := os.Getwd()
		// kubeconfigDir = filepath.Join(projDir, "k8s", "conf")
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
