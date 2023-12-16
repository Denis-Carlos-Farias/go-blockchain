package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/Denis-Carlos-Farias/go-blockchain/blockchain"
)

type CommandLine struct {
	blockchain *blockchain.Blockchain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usege: ")
	fmt.Println("add -block BLOCK_DATA - Adiciona um bloco na blockchain")
	fmt.Println("print - Imprime os blocos da blockchain")
}

func (cli *CommandLine) ValidateArgs() {
	if len(os.Args) < 2 {
		{
			cli.printUsage()
			runtime.Goexit()
		}
	}
}

func (cli *CommandLine) AddBlockchain(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Bloco adicionado com sucesso.")
}

func (cli *CommandLine) PrintChain() {
	interator := cli.blockchain.Interator()

	for {

		block := interator.Next()

		fmt.Printf("Previous Hash: %x\r\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\r\n", block.Data)
		fmt.Printf("Hash: %x\r\n", block.Hash)

		proofOfWork := blockchain.NewProof(block)
		fmt.Printf("Proof of work: %s\n", strconv.FormatBool(proofOfWork.Validate()))
		fmt.Println("---------------------------------------------------------")

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) Run() {
	cli.ValidateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)

	addBlockData := addBlockCmd.String("block", "", "Block Data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handler(err)
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handler(err)
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.AddBlockchain(*addBlockData)

	}
	if printChainCmd.Parsed() {
		cli.PrintChain()
	}
}

func main() {
	defer os.Exit(0)
	chain := blockchain.InitBlockchain()

	defer chain.Database.Close()

	cli := CommandLine{chain}
	cli.Run()
}
