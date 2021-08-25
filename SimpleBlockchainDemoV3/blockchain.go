package main

import (
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

// 实现创建区块链的方法
func NewBlockChain() *BlockChain {

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
		b := tx.Bucket([]byte(blockBucketName))

		if b == nil {
			//如果b1为空，说明名字为"buckeName1"这个桶不存在，我们需要创建之
			fmt.Printf("bucket不存在，准备创建!\n")
			b, err = tx.CreateBucket([]byte(blockBucketName))

			if err != nil {
				log.Panic(err)
			}

			//抽屉准备完毕，开始添加创世块
			genesisBlock := NewBlock(genesisInfo, []byte{})
			b.Put(genesisBlock.Hash, genesisBlock.Serialize() /*将区块序列化，转成字节流*/)
			b.Put([]byte(lastHashKey), genesisBlock.Hash)

			//为了测试，我们把写入的数据读取出来，如果没问题，注释掉这段代码
			//blockInfo := b.Get(genesisBlock.Hash)
			//block := Deserialize(blockInfo)
			//fmt.Printf("解码后的block数据:%s\n", block)

			tail = genesisBlock.Hash
		} else {
			tail = b.Get([]byte(lastHashKey))
		}

		return nil
	})

	return &BlockChain{db, tail}

}

// 5. 添加区块
func (bc *BlockChain) AddBlock(data string) {
	//1. 创建一个区块
	bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blockBucket"))

		if b == nil {
			fmt.Printf("bucket不存在，请检查!\n")
			os.Exit(1)
		}

		block := NewBlock(data, bc.tail)
		b.Put(block.Hash, block.Serialize() /*将区块序列化，转成字节流*/)
		b.Put([]byte("lastHashKey"), block.Hash)

		bc.tail = block.Hash

		return nil
	})
}
