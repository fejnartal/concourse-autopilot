package cmd

import (
	"fmt"
	"os"
)

func Init() {
	fmt.Println("Init world")
	if len(os.Args) != 3 {
		fmt.Println("Wrong parameters")
	}
}
