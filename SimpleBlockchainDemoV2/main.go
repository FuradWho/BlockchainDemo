package main

import (
	"fmt"
	"time"
)

func main() {

	bc := NewBlockChain()
	fmt.Println("创建一个区块链")

	bc.AddBlock("新建一个区块，名字为No.1")
	bc.AddBlock("新建一个区块，名字为No.2")

	for i, block := range bc.Blocks {

		fmt.Printf("--------------\n")
		fmt.Printf("区块高度: %d\n", i)
		fmt.Printf("Version : %d\n", block.Version)
		fmt.Printf("PrevBlockHash : %x\n", block.PrevBlockHash)
		fmt.Printf("MerKleRoot : %x\n", block.MerKleRoot)

		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("TimeStamp : %s\n", timeFormat)

		fmt.Printf("Difficulity : %d\n", block.Difficulity)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("Data : %s\n", block.Data)

		pow := NewProofOfWork(block)
		fmt.Printf("IsValid : %v\n", pow.IsValid())

	}

}
