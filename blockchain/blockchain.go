package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./temp/blocks"
)

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockchainInterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockchain() *Blockchain {
	var lastHash []byte
	opts := badger.DefaultOptions(dbPath)
	//Dir Chaves e valores
	opts.Dir = dbPath
	//ValueDir Valores
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)

	Handler(err)

	err = db.Update(func(txn *badger.Txn) error {
		//lh Last Hash
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Genesis proved")

			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handler(err)

			err = txn.Set([]byte("lh"), genesis.Hash)
			lastHash = genesis.Hash

			return err

		} else {
			item, err := txn.Get([]byte("lh"))
			Handler(err)

			errValue := item.Value(func(val []byte) error {
				lastHash = append([]byte{}, val...)
				return nil
			})
			Handler(errValue)
			return err
		}
	})
	Handler(err)

	blockchain := Blockchain{lastHash, db}
	return &blockchain
}

func (chain *Blockchain) AddBlock(data string) {
	var lastHash []byte
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handler(err)
		errValue := item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
		})
		Handler(errValue)
		return err
	})
	Handler(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {

		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handler(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		return err
	})

	Handler(err)
}

func (chain *Blockchain) Interator() *BlockchainInterator {
	interator := &BlockchainInterator{chain.LastHash, chain.Database}

	return interator
}

func (interator *BlockchainInterator) Next() *Block {
	var block *Block
	var blockData []byte
	err := interator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(interator.CurrentHash)
		Handler(err)
		errValue := item.Value(func(val []byte) error {
			blockData = append([]byte{}, val...)
			return nil
		})
		Handler(errValue)

		block = Deserialize(blockData)

		return err
	})
	Handler(err)

	interator.CurrentHash = block.PrevHash

	return block
}
