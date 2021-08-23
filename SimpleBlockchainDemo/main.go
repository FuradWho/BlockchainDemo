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

func main() {

	block := NewBlock(genesisInfo, []byte{0x0000000000000000})

	fmt.Println("创建第一个区块")
	fmt.Printf("PrevBlockHash : %x\n", block.PrevBlockHash)
	fmt.Printf("Hash : %x\n", block.Hash)
	fmt.Printf("Data : %s\n", block.Data)
}
