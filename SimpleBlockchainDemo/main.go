package main

import (
	"crypto/sha256"
	"fmt"
)

/*
大致流程：

	1. 定义结构（区块头的字段比正常的少）
		1. 前区块哈希
		2. 当前区块哈希
		3. 数据
	2. 创建区块
	3. 生成哈希
	4. 引入区块链
	5. 添加区块
	6. 重构代码

*/

const genesisInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

// 1.定义结构（区块头的字段比正常的少）
type Block struct {
	PrevBlockHash []byte //前区块哈希
	Hash          []byte //当前区块哈希
	Data          []byte //数据，目前使用字节流
}

// 2. 创建区块
func NewBlock(data string, prevBlockHash []byte) *Block {

	block := Block{
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{}, //先填充为空，后续会填充数据
		Data:          []byte(data),
	}

	block.SetHash() //生成Hash值

	return &block
}

// 我们实现一个简单的函数，去进行哈希值的计算，没有随机值，没有难度值
func (block *Block) SetHash() {
	var data []byte
	data = append(data, block.PrevBlockHash...) //使用前区块的hash值和该区块的数据
	data = append(data, block.Data...)

	hash := sha256.Sum256(data)

	block.Hash = hash[:]
}

// 4. 引入区块链，我们使用一个Block数组把它看作为区块链，每一次新增一个区块
type BlockChain struct {
	Blocks []*Block
}

// 实现创建区块链的方法
func NewBlockChain() *BlockChain {

	//在创建一个区块链的时候，添加一个区块，为 初始块

	genesisBlock := NewBlock(genesisInfo, []byte{0x0000000000000000})

	blockChain := BlockChain{
		Blocks: []*Block{genesisBlock},
	}

	return &blockChain

}

// 5. 添加区块
func (bc *BlockChain) AddBlock(data string) {

	//1.创建一个新的区块，并且去添加上一个区块的hash值，以及数据

	//bc.Blocks的最后一个区块的Hash值就是当前新区块的PrevBlockHash
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	prevHash := lastBlock.Hash

	block := NewBlock(data, prevHash)

	//2. 添加到bc.Blocks数组中
	bc.Blocks = append(bc.Blocks, block)

}

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
