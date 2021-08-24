package main

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
