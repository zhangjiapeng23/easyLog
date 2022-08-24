package util

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func PrintSplitLine(symbol string) {
	// get terminal width
	width, _, _ := term.GetSize(int(os.Stdin.Fd()))
	// fix windows does not work
	if width == 0 {
		width = 100
	}
	for i := 0; i < width; i++ {
		fmt.Print(symbol)
	}
	fmt.Println()
}
