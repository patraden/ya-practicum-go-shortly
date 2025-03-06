package main

import (
	"fmt"
	"os"
)

func callExit() {
	fmt.Println("Exiting using func call")
	os.Exit(4) // indirect call from main function (should be ignored)
}

func main() {
	os.Exit(1) // want "direct call to os.Exit found in ..."

	defer func() {
		fmt.Println("Exiting using defer")
		os.Exit(2) // want "direct call to os.Exit found in ..."
	}()

	if true {
		os.Exit(3) // want "direct call to os.Exit found in ..."
	}

	callExit()
}
