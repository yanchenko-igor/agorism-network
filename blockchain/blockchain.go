// blockchain/blockchain.go

package blockchain

import (
	"bytes"
	"crypto/sha256"
	"math/big"
	"math/rand"
	"time"
	//"strconv"
	//"log"
	"fmt"
)

const (
	initialDifficulty = 3
)

// ProofOfWork represents a proof-of-work
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// Block represents a block in the blockchain
type Block struct {
	Timestamp     int64  // Timestamp of when the block was created
	Data          []byte // Data stored in the block
	PrevBlockHash []byte // Hash of the previous block
	Hash          []byte // Hash of the block
	Nonce         int    // Nonce used to create the proof of work
	Difficulty    int    // Difficulty of the proof of work
}

type Blockchain struct {
	Blocks        []*Block
	updateChannel chan struct{}
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}, make(chan struct{})}
}

func (bc *Blockchain) Difficulty() int {
	latestBlock := bc.GetLatestBlock()
	if latestBlock == nil {
		return initialDifficulty
	}
	return latestBlock.Difficulty
}

func NewBlock(prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte{}, prevBlockHash, []byte{}, 0, 0}
	block.BuildBlockData()
	block.Hash = block.CalculateHash()
	block.Difficulty = CalculateDifficulty()
	return block
}

func (b *Block) BuildBlockData() {
	// TODO: Implement function to build block data
	b.Data = []byte("Empty block")
}

func CalculateDifficulty() int {
	// TODO: Implement function to calculate difficulty
	return initialDifficulty
}

// NewBlock creates a new block and calculates its hash
func NewBlockOld(data string, prevBlockHash []byte, difficulty int) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0, difficulty}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewProofOfWork creates a new proof-of-work
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-b.Difficulty))

	pow := &ProofOfWork{b, target}

	return pow
}

// Run performs a proof-of-work
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	maxNonce := int(^uint32(0))
	rand.Seed(time.Now().UnixNano())

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	return nonce, hash[:]
}

// Validate validates a proof-of-work
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	return bytes.Join([][]byte{
		pow.block.PrevBlockHash,
		pow.block.Data,
		intToHex(pow.block.Timestamp),
		intToHex(int64(nonce)),
		intToHex(int64(pow.block.Difficulty)),
	}, []byte{})
}

func (bc *Blockchain) GetLatestBlock() *Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

func MineBlock(blockchain *Blockchain, newBlockChannel chan *Block) {
	for {
		latestBlock := blockchain.GetLatestBlock()
		newBlock := NewBlock(latestBlock.Hash)
		newBlock.BuildBlockData()
		newBlock.Difficulty = CalculateDifficulty()
		select {
		case <-blockchain.updateChannel:
			continue
		default:
		}
		// Keep mining until the proof of work is valid
		for !newBlock.HasValidProofOfWork() {
			newBlock.Nonce++
			newBlock.Hash = newBlock.CalculateHash()
		}

		// Send the new block to the channel
		newBlockChannel <- newBlock
	}
}

func isBlockValid(newBlock, latestBlock *Block) bool {
	// Check if the new block's previous hash value is the same as the latest block's hash value
	if bytes.Compare(newBlock.PrevBlockHash, latestBlock.Hash) != 0 {
		return false
	}

	// Check if the new block's hash value is valid
	if !isHashValid(newBlock.Hash, newBlock.Difficulty) {
		return false
	}

	// If both checks pass, the block is valid
	return true
}

func NewGenesisBlock() *Block {
	return NewBlockOld("Genesis block", []byte{}, initialDifficulty)
}

func (bc *Blockchain) AddBlock(block *Block) error {
	latestBlock := bc.GetLatestBlock()
	if !bytes.Equal(block.PrevBlockHash, latestBlock.Hash) {
		return fmt.Errorf("invalid block")
	}
	if !block.HasValidProofOfWork() {
		return fmt.Errorf("invalid proof of work")
	}
	bc.Blocks = append(bc.Blocks, block)
	// send update to the update channel
	bc.updateChannel <- struct{}{}
	return nil
}

func (b *Block) HasValidProofOfWork() bool {
	return isHashValid(b.Hash, b.Difficulty)
}

func (b *Block) CalculateHash() []byte {
	data := []byte(fmt.Sprintf("%d%x%x%d%d", b.Timestamp, b.PrevBlockHash, b.Data, b.Nonce, b.Difficulty))
	hash := sha256.Sum256(data)
	return hash[:]
}

func HandleNewBlocks(bc *Blockchain, newBlockChannel chan *Block) {
	for {
		select {
		case newBlock := <-newBlockChannel:
			//fmt.Printf("%x", newBlock)
			if err := bc.AddBlock(newBlock); err != nil {
				fmt.Printf("Error: %+v\n", err)
			} else {
				fmt.Printf("Prev. hash: %x\n", newBlock.PrevBlockHash)
				fmt.Printf("Data: %s\n", newBlock.Data)
				fmt.Printf("Hash: %x\n", newBlock.Hash)
				fmt.Printf("Nonce: %x\n", newBlock.Nonce)
			}
			//bc.AddBlock(newBlock)

			// TODO: Validate the new block
			// TODO: Check if the new block has a valid proof of work
			// TODO: Check if the new block's previous hash matches the latest block in the blockchain
			// TODO: If all checks pass, add the new block to the blockchain
			// TODO: Send the new block to all other peers in the network
		default:
		}
	}
}
