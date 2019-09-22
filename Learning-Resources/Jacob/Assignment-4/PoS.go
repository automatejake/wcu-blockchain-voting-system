package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
	// "bufio"
	// "encoding/hex"
	// "encoding/json"
	// "io"
	// "log"
	// "net"
	// "os"
	// "strconv"
	// "time"
	// "github.com/davecgh/go-spew/spew"
	// "github.com/joho/godotenv"
)

// Block represents each 'item' in the blockchain
type Block struct {
	Index     int
	Timestamp string
	Message   string
	Hash      string
	PrevHash  string
}

// Blockchain is a series of validated Blocks
var Blockchain []Block

// SHA256 hashing
func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + block.Message + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// create a new block using previous block's hash
func generateBlock(oldBlock Block, message string) (Block, error) {
	t := time.Now()
	var newBlock Block
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Message = message
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)
	return newBlock, nil
}

// make sure block is valid by checking index, and comparing the hash of the previous block
func isBlockValid(newBlock, oldBlock Block) bool {
	if newBlock.Index != oldBlock.Index+1 {
		return false
	}
	if newBlock.PrevHash != oldBlock.Hash {
		return false
	}
	if newBlock.Hash != calculateHash(newBlock) {
		return false
	}
	return true
}

// make sure the chain we're checking is longer than the current blockchain
func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func main() {
	fmt.Println("Proof of stake chain")
}
