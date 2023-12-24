package blockchain

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
)

const (
	dbPath      = "./temp/blocks"
	dbFile      = "./temp/blocks/MANIFEST"
	genesisData = "First Transaction from Genesis"
)

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockchainInterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func DbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func ContinueBlockchain(address string) *Blockchain {
	if !DbExists() {
		fmt.Println("no existing Blockchain found, please create one.")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	//Dir Chaves e valores
	opts.Dir = dbPath
	//ValueDir Valores
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handler(err)

	err = db.Update(func(txn *badger.Txn) error {
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
	chain := Blockchain{lastHash, db}

	return &chain
}

func InitBlockchain(address string) *Blockchain {
	var lastHash []byte

	if DbExists() {
		fmt.Println("Blockchain already exist")
		runtime.Goexit()
	}
	opts := badger.DefaultOptions(dbPath)
	//Dir Chaves e valores
	opts.Dir = dbPath
	//ValueDir Valores
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)

	Handler(err)

	err = db.Update(func(txn *badger.Txn) error {
		//lh Last Hash
		cbtx := CoinBaseTx(address, genesisData)
		genesis := Genesis(cbtx)
		fmt.Println("Genesis has been created.")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handler(err)
		err = txn.Set([]byte("lh"), genesis.Hash)
		lastHash = genesis.Hash
		return err
	})
	Handler(err)

	blockchain := Blockchain{lastHash, db}
	return &blockchain
}

func (chain *Blockchain) AddBlock(transactions []*Transaction) {
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

	newBlock := CreateBlock(transactions, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handler(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash
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

func (chain *Blockchain) FindUnspentTransctions(address string) []Transaction {
	var unspentTxs []Transaction
	spentTxOs := make(map[string][]int)

	iter := chain.Interator()

	for {
		block := iter.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTxOs[txID] != nil {
					for _, spentOut := range spentTxOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.CanBrUnlocked(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}

			if !tx.IsCoinBase() {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTxOs[inTxID] = append(spentTxOs[inTxID], in.Out)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return unspentTxs
}

func (chain *Blockchain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := chain.FindUnspentTransctions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBrUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (chain *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransctions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.CanBrUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}
