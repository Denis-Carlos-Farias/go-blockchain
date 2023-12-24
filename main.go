package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/Denis-Carlos-Farias/go-blockchain/blockchain"
)

type CommandLine struct{}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usege: ")
	fmt.Println("getbalance -address ADDRESS - get balance for address")
	fmt.Println("createblockchain -address ADDRESS - creats a blockchain")
	fmt.Println("printchain - Prints the blocks in the chain")
	fmt.Println("send -from FROM -to TO -amount AMOUNT - send amount from to")
}

func (cli *CommandLine) ValidateArgs() {
	if len(os.Args) < 2 {
		{
			cli.printUsage()
			runtime.Goexit()
		}
	}
}

func (cli *CommandLine) PrintChain() {
	chain := blockchain.ContinueBlockchain("")
	defer chain.Database.Close()

	interator := chain.Interator()

	for {
		block := interator.Next()

		fmt.Printf("Previous Hash: %x\r\n", block.PrevHash)
		fmt.Printf("Hash: %x\r\n", block.Hash)

		proofOfWork := blockchain.NewProof(block)
		fmt.Printf("Proof of work: %s\n", strconv.FormatBool(proofOfWork.Validate()))
		fmt.Println("---------------------------------------------------------")

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) createBlockchain(address string) {
	chain := blockchain.InitBlockchain(address)
	chain.Database.Close()
	fmt.Println("Finished")
}

func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlockchain(address)
	defer chain.Database.Close()

	balance := 0

	UTXOs := chain.FindUTXO(address)
	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockchain(from)
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)

	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Sucess!")
}

func (cli *CommandLine) Run() {
	cli.ValidateArgs()

	createBlockchaincmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalance := flag.NewFlagSet("getbalance", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	getBalanceAddress := getBalance.String("address", "", "address of wallet")
	createBlockchainAdress := createBlockchaincmd.String("address", "", "address of miner")
	sendFrom := sendCmd.String("from", "", "adress of wallet")
	sendTo := sendCmd.String("to", "", "address of wallet")
	sendAmount := sendCmd.Int("amount", 0, "amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalance.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchaincmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalance.Parsed() {
		if *getBalanceAddress == "" {
			getBalance.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)

	}
	if createBlockchaincmd.Parsed() {
		if *createBlockchainAdress == "" {
			createBlockchaincmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockchain(*createBlockchainAdress)

	}
	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)

	}
	if printChainCmd.Parsed() {
		cli.PrintChain()
	}
}

func main() {
	defer os.Exit(0)

	cli := CommandLine{}
	cli.Run()
}
