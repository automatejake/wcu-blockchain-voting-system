//https://medium.com/@mycoralhealth/part-2-networking-code-your-own-blockchain-in-less-than-200-lines-of-go-17fe1dad46e1

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Block struct {
	Index     int
	Timestamp string
	Message   string
	PrevHash  string
	Nonce     string
	Hash      string
}

// Blockchain is a series of validated Blocks
var Blockchain []Block

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + block.Message + block.PrevHash + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block, message string) (Block, error) {
	var newBlock Block

	t := time.Now()
	newBlock.Timestamp = t.String()
	newBlock.Message = "this is a test" //insert function to prompt user for messsage
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Index = oldBlock.Index + 1
	//insert function to calculate nonce and check if nonce is valid hash
	return newBlock, nil
}

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

// bcServer handles incoming concurrent Blocks, this is used instead of mutex.Lock and mutex.Unlock
// what is the difference?
var bcServer chan []Block

func main() {
	fmt.Print("networking tutorial")
}
