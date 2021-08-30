package main

import (
	"fmt"
	"os"
)

func mn() {
	cmds := os.Args

	for i, cmd := range cmds {
		fmt.Printf("cmd[%d] : %s\n", i, cmd)
	}
}
