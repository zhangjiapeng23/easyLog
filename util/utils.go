package util

import "fmt"

func PrintSplitLine(symbol string) {
	for i := 0; i < 100; i++ {
		fmt.Print(symbol)
	}
	fmt.Println()
}