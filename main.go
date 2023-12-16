package main

import (
	"fmt"

	"github.com/Denis-Carlos-Farias/go-blockchain/model"
)

func main() {
	chain := model.InitBlockchain()

	chain.AddBlock("First block after Genesis's block")
	chain.AddBlock("Second block after Genesis's block")
	chain.AddBlock("Third block after Genesis's block")

	for _, block := range chain.Blocks {
		fmt.Printf("Previous Hash: %x\r\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\r\n", block.Data)
		fmt.Printf("Hash: %x\r\n", block.Hash)
		fmt.Println("---------------------------------------------------------")
	}
}
