package main

import (
	"fmt"
)

func main() {

	bc := NewBlockChain()
	fmt.Println("创建一个区块链")

	bc.AddBlock("新建一个区块，名字为first")

	for i, block := range bc.Blocks {

		fmt.Printf("--------------\n")
		fmt.Printf("区块高度: %d\n", i)
		fmt.Printf("PrevBlockHash : %x\n", block.PrevBlockHash)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("Data : %s\n", block.Data)
	}

}
