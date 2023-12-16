package main

import (
	"fmt"
	"strconv"

	"github.com/Denis-Carlos-Farias/go-blockchain/blockchain"
)

func main() {
	chain := blockchain.InitBlockchain()

	chain.AddBlock("First block after Genesis's block")
	chain.AddBlock("Second block after Genesis's block")
	chain.AddBlock("Third block after Genesis's block")

	for _, block := range chain.Blocks {
		fmt.Printf("Previous Hash: %x\r\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\r\n", block.Data)
		fmt.Printf("Hash: %x\r\n", block.Hash)

		proofOfWork := blockchain.NewProof(block)
		fmt.Printf("Proof of work: %s\n", strconv.FormatBool(proofOfWork.Validate()))
		fmt.Println("---------------------------------------------------------")

	}
}
