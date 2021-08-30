package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

/*

定义一个工作量证明的结构ProofOfWork
	a. block
	b. 目标值
提供创建POW的函数
	NewProofOfWork(参数)
提供计算不断计算hash的哈数
	Run()
提供一个校验函数
	IsValid()
*/

type ProofOfWork struct {
	block *Block

	//来存储哈希值，它内置一些方法Cmp:比较方法
	// SetBytes : 把bytes转成big.int类型 []byte("0x00000919011eeb8fbdf0c476d8510b8e1e632eba7b584ac04c11ad20cbbdd394")
	// SetString : 把string转成big.int类型 "0x00000919011eeb8fbdf0c476d8510b8e1e632eba7b584ac04c11ad20cbbdd394"
	target *big.Int //系统提供的，是固定的
}

const Bits = 20

func NewProofOfWork(block *Block) *ProofOfWork {

	pow := ProofOfWork{
		block: block,
	}

	//先去初始化  0000000000000000000000000000000000000000000000000000000000000001
	bigIntTmp := big.NewInt(1)
	//左移256 - Bits 动态调整1的位置
	bigIntTmp.Lsh(bigIntTmp, 256-Bits)
	pow.target = bigIntTmp

	return &pow

}

//这是pow的运算函数，为了获取挖矿的随机数，同时返回区块的哈希值
func (pow *ProofOfWork) Run() ([]byte, uint64) {

	//1. 获取block数据
	//2. 拼接nonce
	//3. sha256
	//4. 与难度值比较
	//		a. 哈希值大于难度值，nonce++
	//		b. 哈希值小于难度值，挖矿成功,退出

	//定义一个nonce
	var nonce uint64
	//定义返回值的Hash值
	var hash [32]byte

	for {

		//data := block + nonce
		hash = sha256.Sum256(pow.prepareData(nonce))

		//将hash（数组类型）转成big.int, 然后与pow.target进行比较, 需要引入局部变量
		var bigIntTmp big.Int
		bigIntTmp.SetBytes(hash[:])

		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		if bigIntTmp.Cmp(pow.target) == -1 {
			// x < y ， 挖矿成功！
			fmt.Printf("挖矿成功！nonce: %d, 哈希值为: %x\n", nonce, hash)
			break
		} else {
			nonce++
		}
	}
	return hash[:], nonce

}

//定义一个prepareData函数，用来进行 区块数据 和 nonce 的拼接
func (pow *ProofOfWork) prepareData(nonce uint64) []byte {
	block := pow.block

	tmp := [][]byte{
		uintToByte(block.Version),
		block.PrevBlockHash,
		block.MerKleRoot,
		uintToByte(block.TimeStamp),
		uintToByte(block.Difficulity),
		//block.Data,
		uintToByte(nonce),
	}

	data := bytes.Join(tmp, []byte{})
	return data
}

//校验 ，block的数据是完整的，我们要做的是校验一下，Hash，block数据，和Nonce是否满足难度值要求
func (pow *ProofOfWork) IsValid() bool {
	//获取block数据
	//拼接nonce
	//做sha256
	//比较

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)

	var tmp big.Int
	tmp.SetBytes(hash[:])

	return tmp.Cmp(pow.target) == -1
}
