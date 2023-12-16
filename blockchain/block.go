package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash     []byte //Hash do bloco
	Data     []byte //O dado que deseja ser armazenado
	PrevHash []byte //Hash do ultimo bloco
	nonce    int
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
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	proofOfWork := NewProof(block)

	nonce, hash := proofOfWork.Run()

	block.nonce = nonce
	block.Hash = hash[:]

	return block
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	Handler(err)

	return result.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)
	Handler(err)

	return &block
}

func Handler(err error) {
	if err != nil {
		log.Panic(err)
	}
}
