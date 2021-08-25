package main

import (
	"fmt"
)

func main() {

	bc := NewBlockChain()
	fmt.Println("创建一个区块链")

	defer bc.db.Close()

	cli := CLI{bc}
	cli.Run()

}
