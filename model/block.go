package model

import (
	"bytes"
	"crypto/sha256"
)

type Block struct {
	Hash     []byte //Hash do bloco
	Data     []byte //O dado que deseja ser armazenado
	PrevHash []byte //Hash do ultimo bloco
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()

	return block
}
