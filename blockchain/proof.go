package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

//Pegar o Data do bloco
//Criar um contador (nonce) na qual iniará em 0
//Criar um Hash do Data + o contador
//Verificar o Hash para ver se ele atende o conjunto de requisitos
//Requisitos:
//Os primeiros bytes devem conter zeros

const Difficulty = 5

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	proofOfWork := &ProofOfWork{b, target}

	return proofOfWork
}

func (proofOfWork *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			proofOfWork.Block.PrevHash,
			proofOfWork.Block.HashTransactions(),
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)

	return data
}

func (proofOfWork *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := proofOfWork.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(proofOfWork.Target) == -1 {
			break
		}
		nonce++
	}
	fmt.Println()
	return nonce, hash[:]
}

func (proofOfWork *ProofOfWork) Validate() bool {
	var intHash big.Int
	data := proofOfWork.InitData(proofOfWork.Block.nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(proofOfWork.Target) == -1
}

func ToHex(number int64) []byte {
	buff := new(bytes.Buffer)

	err := binary.Write(buff, binary.BigEndian, number)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
