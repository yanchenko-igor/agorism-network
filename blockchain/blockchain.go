package blockchain

import (
	"crypto/sha256"
	//"encoding/hex"
	"fmt"
	"strconv"
	//"strings"
	"bytes"
	// "database/sql"
	// "encoding/json"
	// _ "github.com/mattn/go-sqlite3"
	"math"
	"math/big"
	// "net/http"
	"time"
)

const targetBits = 24

// Block represents a block in the blockchain
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
	Difficulty    int
}

// Blockchain represents a chain of blocks
type Blockchain struct {
	Blocks []*Block
}

// NewBlock creates a new block and adds it to the chain
func (bc *Blockchain) FindBlock(data string) *Block {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	block := &Block{time.Now().Unix(), []byte(data), prevBlock.Hash, []byte{}, 0, 5}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// AddBlock adds a new block to the chain.
func (bc *Blockchain) AddBlock(block *Block) error {
	// Check if the block can be added to the current chain.
	if !bc.IsValidBlock(block) {
		return fmt.Errorf("cannot add block to chain")
	}
	bc.Blocks = append(bc.Blocks, block)
	return nil
}

func (bc *Blockchain) IsValidBlock(block *Block) bool {
	pow := NewProofOfWork(block)
	return pow.Validate()
}

// ProofOfWork represents a proof of work
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork creates a new proof of work
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-b.Difficulty))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(pow.block.Difficulty)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// Run performs a proof of work
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	//fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < math.MaxInt64 {
		//TODO remove next line, it's here to create a delay without loading the CPU
		time.Sleep(100 * time.Millisecond)
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		//fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	//fmt.Print("\n\n")

	return nonce, hash[:]
}

// Validate validates a proof of work
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}

// NewGenesisBlock creates the first block in the chain
func (bc *Blockchain) NewGenesisBlock() *Block {
	genesisBlock := &Block{time.Now().Unix(), []byte("Genesis Block"), []byte{}, []byte{}, 0, 10}
	genesisBlock.setHash()
	bc.Blocks = append(bc.Blocks, genesisBlock)
	return genesisBlock
}

func (b *Block) setHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

// targetTime is the target time in seconds for finding the previous n blocks
func adjustDifficulty(lastBlocks []*Block, targetTime int) int {
	newDifficulty := lastBlocks[len(lastBlocks)-1].Difficulty

	// Calculate the time it took to find the previous n blocks
	elapsedTime := lastBlocks[len(lastBlocks)-1].Timestamp - lastBlocks[0].Timestamp

	if elapsedTime < int64(targetTime)/2 {
		newDifficulty++
	} else if elapsedTime > int64(targetTime)*2 {
		newDifficulty--
	}

	return newDifficulty
}

func IntToHex(num int64) []byte {
	return []byte(strconv.FormatInt(num, 16))
}

func (b *Block) calcHash(nonce int) []byte {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	nonceStr := []byte(strconv.FormatInt(int64(nonce), 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp, nonceStr}, []byte{})
	hash := sha256.Sum256(headers)
	return hash[:]
}

//
// type CreateBlockRequest struct {
// 	Data string `json:"data"`
// }
//
// type CreateBlockResponse struct {
// 	Message string `json:"message"`
// }
//
// func wrapper(bc *Blockchain) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var request CreateBlockRequest
// 		err := json.NewDecoder(r.Body).Decode(&request)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 		fmt.Printf("%+v", bc)
// 		bc.NewBlock(request.Data)
// 		response := CreateBlockResponse{Message: "Block created successfully"}
// 		json.NewEncoder(w).Encode(response)
// 	}
// }
//
// func (bc *Blockchain) Load() error {
// 	db, err := sql.Open("sqlite3", "./blockchain.db")
// 	if err != nil {
// 		return err
// 	}
// 	defer db.Close()
//
// 	rows, err := db.Query("SELECT data, prev_block_hash, hash, nonce FROM blocks")
// 	if err != nil {
// 		return err
// 	}
// 	defer rows.Close()
//
// 	var blocks []*Block
// 	for rows.Next() {
// 		var data []byte
// 		var prevBlockHash []byte
// 		var hash []byte
// 		var nonce int
// 		if err := rows.Scan(&data, &prevBlockHash, &hash, &nonce); err != nil {
// 			return err
// 		}
//
// 		block := &Block{Data: data, PrevBlockHash: prevBlockHash, Hash: hash, Nonce: nonce}
// 		// if err := bc.ValidateBlock(block, prevBlockHash); err != nil {
// 		// 	return err
// 		// }
// 		blocks = append(blocks, block)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return err
// 	}
//
// 	bc.blocks = blocks
// 	return nil
// }

var bc *Blockchain

// func main() {
// 	bc := Blockchain{}
// 	bc.NewGenesisBlock()
// 	bc.NewBlock("Send 1 BTC to Alice")
// 	bc.NewBlock("Send 2 more BTC to Alice")
// 	bc.NewBlock("Send 5 more BTC to Alice")
// 	bc.NewBlock("Send 2 BTC to Bob")
// 	bc.NewBlock("Send 4 more BTC to Alice")
//
// 	for _, block := range bc.blocks {
// 		//fmt.Printf("%v\n", block)
// 		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
// 		fmt.Printf("Data: %s\n", block.Data)
// 		fmt.Printf("Hash: %x\n", block.Hash)
// 		fmt.Printf("Nonce: %x\n", block.Nonce)
// 		fmt.Println()
// 	}
// 	fmt.Printf("%+v", bc)
// 	http.HandleFunc("/createBlock", wrapper(&bc))
// 	http.ListenAndServe(":8080", nil)
// }
