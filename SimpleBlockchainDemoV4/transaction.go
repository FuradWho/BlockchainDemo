package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

/*
	定义交易结构
	定义input
	定义output
	设置交易ID
*/

type TXInput struct {
	TXID    []byte //交易id
	Index   int64  //output的索引
	Address string //解锁脚本，先使用地址来模拟
}

type TXOutput struct {
	Value   float64 //转账金额
	Address string  //锁定脚本
}

type Transaction struct {
	TXid      []byte     //交易id
	TXInputs  []TXInput  //所有的inputs
	TXOutputs []TXOutput //所有的outputs
}

func (tx *Transaction) SetTXID() {

	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(tx)

	if err != nil {
		log.Panic(err)
	}

	hash := sha256.Sum256(buffer.Bytes())
	tx.TXid = hash[:]
}

//把挖矿的人传递进来，因为有奖励
func NewCoinbaseTx(miner string, data string) *Transaction {

	//我们在后面的程序中，需要识别一个交易是否为coinbase，所以我们需要设置一些特殊的值，用于判断
	inputs := []TXInput{TXInput{nil, -1, data}}
	outputs := []TXOutput{TXOutput{12.5, miner}}

	tx := Transaction{nil, inputs, outputs}
	tx.SetTXID()

	return &tx
}

//内部逻辑：
//

func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {

	utxos := make(map[string][]int64) //标识能用的utxo
	var resValue float64              //这些utxo存储的金额

	//1. 遍历账本，找到属于付款人的合适的金额，把这个outputs找到
	utxos, resValue = bc.FindNeedUtxos(from, amount)

	//2. 如果找到钱不足以转账，创建交易失败。
	if resValue < amount {
		fmt.Printf("余额不足，交易失败!\n")
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput

	//3. 将outputs转成inputs
	for txid /*0x333*/, indexes := range utxos {
		for _, i /*0, 1*/ := range indexes {
			input := TXInput{[]byte(txid), i, from}
			inputs = append(inputs, input)
		}
	}

	//4. 创建输出，创建一个属于收款人的output
	output := TXOutput{amount, to}
	outputs = append(outputs, output)

	//5. 如果有找零，创建属于付款人output
	if resValue > amount {
		output1 := TXOutput{resValue - amount, from}
		outputs = append(outputs, output1)
	}

	//创建交易
	tx := Transaction{nil, inputs, outputs}

	//6. 设置交易id
	tx.SetTXID()

	//7. 返回交易结构
	return &tx
}
