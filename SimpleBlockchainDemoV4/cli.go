package main

import (
	"fmt"
	"os"
	"strconv"
)

const Usage = `
	./blockchain addBlock "xxxx"   添加数据到区块链
	./blockchain printChain          打印区块链
	./blockchain getBalance          获取地址的余额
	./blockchain send from to amount miner		转账命令
	
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

	case "printChain":
		fmt.Printf("打印区块链命令调用\n") //打印区块链时候：遍历区块链，不需要外部输入数据
		cli.PrintChain()

	case "getBalance":
		fmt.Printf("获取地址的余额\n") //打印区块链时候：遍历区块链，不需要外部输入数据
		cli.bc.GetBalance(cmds[2])

	case "send":
		fmt.Printf("转账命令\n")
		//./blockchain send FROM TO AMOUNT MINER DATA "转账命令"
		if len(cmds) != 7 {
			fmt.Printf("send命令发现无效参数，请检查!\n")
			fmt.Printf(Usage)
			os.Exit(1)
		}

		from := cmds[2]
		to := cmds[3]
		amount, _ := strconv.ParseFloat(cmds[4], 64)
		miner := cmds[5]
		data := cmds[6]
		cli.Send(from, to, amount, miner, data)

	default:
		fmt.Printf("无效的命令，请检查\n")
		fmt.Printf(Usage)
	}

}