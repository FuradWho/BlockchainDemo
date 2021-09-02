package main

import (
	"SimpleBlockchainDemoV5/base58"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

type TXInput struct {
	TXID      []byte //交易id
	Index     int64  //output的索引
	Signature []byte //交易签名
	PubKey    []byte //公钥本身，不是公钥哈希
}

type TXOutput struct {
	Value      float64 //转账金额
	PubKeyHash []byte  //是公钥的哈希，不是公钥本身
}

//给定转账地址，得到这个地址的公钥哈希，完成对output的锁定
func (output *TXOutput) Lock(address string) {

	//address -> public key hash
	//25字节
	decodeInfo := base58.Decode(address)

	pubKeyHash := decodeInfo[1 : len(decodeInfo)-4]

	output.PubKeyHash = pubKeyHash
}

func NewTXOutput(value float64, address string) TXOutput {
	output := TXOutput{Value: value}
	output.Lock(address)
	return output
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

//挖矿奖励
const reward = 12.5

//实现挖矿挖矿交易，
//特点：只有输出，没有有效的输入(不需要引用id，不需要索引，不需要签名)

//把挖矿的人传递进来，因为有奖励
func NewCoinbaseTx(miner string, data string) *Transaction {

	//我们在后面的程序中，需要识别一个交易是否为coinbase，所以我们需要设置一些特殊的值，用于判断
	inputs := []TXInput{TXInput{nil, -1, nil, []byte(data)}}
	//outputs := []TXOutput{TXOutput{12.5, miner}}

	output := NewTXOutput(reward, miner)
	outputs := []TXOutput{output}

	tx := Transaction{nil, inputs, outputs}
	tx.SetTXID()

	return &tx
}

func (tx *Transaction) IsCoinbase() bool {
	//特点：1. 只有一个input 2. 引用的id是nil 3. 引用的索引是-1
	inputs := tx.TXInputs
	if len(inputs) == 1 && inputs[0].TXID == nil && inputs[0].Index == -1 {
		return true
	}

	return false
}

func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	//1. 打开钱包
	ws := NewWallets()

	//获取秘钥对
	wallet := ws.WalletsMap[from]

	if wallet == nil {
		fmt.Printf("%s 的私钥不存在，交易创建失败!\n", from)
		return nil
	}

	//2. 获取公钥，私钥
	//privateKey := wallet.PrivateKey //目前使用不到，步骤三签名时使用
	publickKey := wallet.PublicKey

	pubKeyHash := HashPubKey(wallet.PublicKey)

	utxos := make(map[string][]int64) //标识能用的utxo
	var resValue float64              //这些utxo存储的金额
	//假如李四转赵六4，返回的信息为:
	//utxos[0x333] = int64{0, 1}
	//resValue : 5

	//1. 遍历账本，找到属于付款人的合适的金额，把这个outputs找到
	utxos, resValue = bc.FindNeedUtxos(pubKeyHash, amount)

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
			input := TXInput{[]byte(txid), i, nil, publickKey}
			inputs = append(inputs, input)
		}
	}

	//4. 创建输出，创建一个属于收款人的output
	//output := TXOutput{amount, to}
	output := NewTXOutput(amount, to)
	outputs = append(outputs, output)

	//5. 如果有找零，创建属于付款人output
	if resValue > amount {
		//output1 := TXOutput{resValue - amount, from}
		output1 := NewTXOutput(resValue-amount, from)
		outputs = append(outputs, output1)
	}

	//创建交易
	tx := Transaction{nil, inputs, outputs}

	//6. 设置交易id
	tx.SetTXID()

	//7. 返回交易结构
	return &tx
}
