package main

import (
	"easyLog/menu"
)


func main() {
	
	for {
		menu.CurrentMenu.ShowMenu()	
	}
	// k8s.Tmp()
	// fmt.Println(k8s.NewClient("pre").ListNamespaces())
	// k := k8s.NewClient("pre")
	// k.ListAppsForNamespace("office")
	// fmt.Print("end")
	// res := k.ListPodsForApp("office", "webull-auth-center")
	// for _, pod := range res.Items {
	// 	fmt.Println(pod.Name)
	// }

}

