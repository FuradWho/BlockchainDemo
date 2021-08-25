package main

import (
	"fmt"
	"os"
)

const Usage = `
	./blockchain addBlock "xxxx"   添加数据到区块链
	./blockchain printChain          打印区块链
`

type CLI struct {
	bc *BlockChain
}

//命令解析，方法调用
func (cli *CLI) Run() {

	cmds := os.Args //获取命令数组

	if len(cmds) < 2 {
		fmt.Printf(Usage)
		os.Exit(1)
	}

	switch cmds[1] {
	case "addBlock":
		data := cmds[2]
		fmt.Printf("添加区块命令调用, 数据：%s\n", data) //添加区块的时候： bc.addBlock(data), data 通过os.Args拿回来
		cli.AddBlock(data)

	case "printChain":
		fmt.Printf("打印区块链命令调用\n") //打印区块链时候：遍历区块链，不需要外部输入数据
		cli.PrintChain()

	default:
		fmt.Printf("无效的命令，请检查\n")
		fmt.Printf(Usage)
	}

}
