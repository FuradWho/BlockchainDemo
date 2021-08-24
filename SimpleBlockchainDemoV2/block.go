package main

import (
	"time"
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
	Version       uint64 //区块版本号
	PrevBlockHash []byte //前区块哈希
	MerKleRoot    []byte //先填写为空，后续使用
	TimeStamp     uint64 //从1970.1.1至今的秒数
	Difficulity   uint64 //挖矿的难度值, v2时使用
	Nonce         uint64 //随机数，挖矿找的就是它!
	Data          []byte //数据，目前使用字节流，v4开始使用交易代替
	Hash          []byte //当前区块哈希, 区块中本不存在的字段，为了方便我们添加进来
}

// 2. 创建区块
func NewBlock(data string, prevBlockHash []byte) *Block {

	block := Block{
		Version:       00,
		PrevBlockHash: prevBlockHash,
		MerKleRoot:    []byte{},
		TimeStamp:     uint64(time.Now().Unix()),
		Difficulity:   Bits, //随便写的，后续调整
		Data:          []byte(data),
		Hash:          []byte{}, //先填充为空，后续会填充数据
	}

	// block.SetHash() 生成Hash值

	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return &block
}

/*
// 我们实现一个简单的函数，去进行哈希值的计算，没有随机值，没有难度值
func (block *Block) SetHash() {
	var data []byte

	//uintToByte将数字转成[]byte{}, 在utils.go实现
	//data = append(data, uintToByte(block.Version)...)
	//data = append(data, block.PrevBlockHash...)
	//data = append(data, block.MerKleRoot...)
	//data = append(data, uintToByte(block.TimeStamp)...)
	//data = append(data, uintToByte(block.Difficulity)...)
	//data = append(data, block.Data...)
	//data = append(data, uintToByte(block.Nonce)...)

	tmp := [][]byte{
		uintToByte(block.Version),
		block.PrevBlockHash,
		block.MerKleRoot,
		uintToByte(block.TimeStamp),
		uintToByte(block.Difficulity),
		block.Data,
		uintToByte(block.Nonce),
	}

	data = bytes.Join(tmp, []byte{})

	hash := sha256.Sum256(data)
	block.Hash = hash[:]
}

*/
