package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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

func main() {
	fmt.Println("Proof of stake chain")
}
