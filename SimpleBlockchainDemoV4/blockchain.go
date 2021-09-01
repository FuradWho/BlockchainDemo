package main

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

// 使用bolt改写
type BlockChain struct {
	db   *bolt.DB //数据库的一个句柄
	tail []byte   //最后一个区块的Hash值

}

const blockChainName = "blockChain.db"
const blockBucketName = "blockBucket"
const lastHashKey = "lastHashKey"

//我们想把FindMyUtoxs和FindNeedUTXO进行整合
//我们可以定义一个结构，同时包含output已经定位信息
type UTXOInfo struct {
	TXID   []byte   //交易id
	Index  int64    //output的索引值
	Output TXOutput //output本身
}

func CreateBlockChain(miner string) *BlockChain {

	if IsFileExist(blockChainName) {
		fmt.Printf("区块链已经存在，不需要重复创建!\n")
		return nil
	}

	//功能分析：
	//1. 获得数据库的句柄，打开数据库，读写数据

	db, err := bolt.Open(blockChainName, 0600, nil)
	//向数据库中写入数据
	//从数据库中读取数据

	if err != nil {
		log.Panic(err)
	}

	//defer db.Close()

	var tail []byte

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blockBucketName))

		if err != nil {
			log.Panic(err)
		}

		//抽屉准备完毕，开始添加创世块
		//创世块中只有一个挖矿交易，只有Coinbase
		coinbase := NewCoinbaseTx(miner, genesisInfo)
		genesisBlock := NewBlock([]*Transaction{coinbase}, []byte{})

		b.Put(genesisBlock.Hash, genesisBlock.Serialize() /*将区块序列化，转成字节流*/)
		b.Put([]byte(lastHashKey), genesisBlock.Hash)

		//为了测试，我们把写入的数据读取出来，如果没问题，注释掉这段代码
		//blockInfo := b.Get(genesisBlock.Hash)
		//block := Deserialize(blockInfo)
		//fmt.Printf("解码后的block数据:%s\n", block)

		tail = genesisBlock.Hash

		return nil
	})

	return &BlockChain{db, tail}
}

//返回区块链实例
func NewBlockChain() *BlockChain {

	if !IsFileExist(blockChainName) {
		fmt.Printf("区块链不存在，请先创建!\n")
		return nil
	}

	//功能分析：
	//1. 获得数据库的句柄，打开数据库，读写数据

	db, err := bolt.Open(blockChainName, 0600, nil)

	if err != nil {
		log.Panic(err)
	}

	//defer db.Close()

	var tail []byte

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucketName))

		if b == nil {
			fmt.Printf("区块链bucket为空，请检查!\n")
			os.Exit(1)
		}

		tail = b.Get([]byte(lastHashKey))

		return nil
	})

	return &BlockChain{db, tail}
}

// 5. 添加区块
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	//1. 创建一个区块
	bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blockBucket"))

		if b == nil {
			fmt.Printf("bucket不存在，请检查!\n")
			os.Exit(1)
		}

		block := NewBlock(txs, bc.tail)
		b.Put(block.Hash, block.Serialize() /*将区块序列化，转成字节流*/)
		b.Put([]byte("lastHashKey"), block.Hash)

		bc.tail = block.Hash

		return nil
	})
}

//实现思路：
func (bc *BlockChain) FindMyUtoxs(pubKeyHash []byte) []UTXOInfo {
	fmt.Printf("FindMyUtoxs\n")
	//var UTXOs []TXOutput //返回的结构
	var UTXOInfos []UTXOInfo //新的返回结构

	it := bc.NewIterator()

	//这是标识已经消耗过的utxo的结构，key是交易id，value是这个id里面的output索引的数组
	spentUTXOs := make(map[string][]int64)

	//1. 遍历账本
	for {

		block := it.Next()

		//2. 遍历交易
		for _, tx := range block.Transactions {
			//遍历交易输入:inputs

			if tx.IsCoinbase() == false {
				//如果不是coinbase，说明是普通交易，才有必要进行遍历
				for _, input := range tx.TXInputs {

					//判断当前被使用input是否为目标地址所有
					if bytes.Equal(HashPubKey(input.PubKey), pubKeyHash) {

						fmt.Printf("找到了消耗过的output! index : %d\n", input.Index)
						key := string(input.TXID)
						spentUTXOs[key] = append(spentUTXOs[key], input.Index)
					}
				}
			}

			key := string(tx.TXid)
			indexes := spentUTXOs[key]

		OUTPUT:
			//3. 遍历output
			for i, output := range tx.TXOutputs {

				if len(indexes) != 0 {
					fmt.Printf("当前这笔交易中有被消耗过的output!\n")
					for _, j := range indexes {
						if int64(i) == j {
							fmt.Printf("i == j, 当前的output已经被消耗过了，跳过不统计!\n")
							continue OUTPUT
						}
					}
				}

				//4. 找到属于我的所有output
				if bytes.Equal(pubKeyHash, output.PubKeyHash) {
					//fmt.Printf("找到了属于 %s 的output, i : %d\n", address, i)
					//UTXOs = append(UTXOs, output)
					utxoinfo := UTXOInfo{tx.TXid, int64(i), output}
					UTXOInfos = append(UTXOInfos, utxoinfo)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			fmt.Printf("遍历区块链结束!\n")
			break
		}
	}

	return UTXOInfos
}

func (bc *BlockChain) GetBalance(address string) {
	utxos := bc.FindMyUtoxs(address)

	var total = 0.0

	for _, utxo := range utxos {
		total += utxo.Value
	}

	fmt.Printf("%s 的余额为: %f\n", address, total)
}

//1. 遍历账本，找到属于付款人的合适的金额，把这个outputs找到
//utxos, resValue = bc.FindNeedUtxos(from, amount)
func (bc *BlockChain) FindNeedUtxos(pubKeyHash []byte, amount float64) (map[string][]int64, float64) {

	needUtxos := make(map[string][]int64) //标识能用的utxo, //返回的结构
	var resValue float64                  //统计的金额

	//复用FindMyUtxo函数，这个函数已经包含了所有信息
	utxoinfos := bc.FindMyUtoxs(pubKeyHash)

	for _, utxoinfo := range utxoinfos {
		key := string(utxoinfo.TXID)

		needUtxos[key] = append(needUtxos[key], int64(utxoinfo.Index))
		resValue += utxoinfo.Output.Value

		//2. 判断一下金额是否足够
		if resValue >= amount {
			//a. 足够， 直接返回
			break
		}
	}
	return needUtxos, resValue
}
