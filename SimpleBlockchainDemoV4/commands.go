package main

import (
	"bytes"
	"fmt"
	"time"
)

func (cli *CLI) AddBlock(txs []*Transaction) {
	cli.bc.AddBlock(txs)
	fmt.Printf("添加区块成功!\n")
}

func (cli *CLI) PrintChain() {

	it := cli.bc.NewIterator()

	for {
		block := it.Next()

		fmt.Printf("--------------\n")
		fmt.Printf("Version : %d\n", block.Version)
		fmt.Printf("PrevBlockHash : %x\n", block.PrevBlockHash)
		fmt.Printf("MerKleRoot : %x\n", block.MerKleRoot)

		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("TimeStamp : %s\n", timeFormat)
		fmt.Printf("Difficulity : %d\n", block.Difficulity)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Hash : %x\n", block.Hash)
		//fmt.Printf("Data : %s\n", block.Transactions)

		pow := NewProofOfWork(block)
		fmt.Printf("IsValid : %v\n", pow.IsValid())

		if bytes.Equal(block.PrevBlockHash, []byte{}) {
			fmt.Printf("区块链遍历结束!\n")
			break
		}
	}
}
