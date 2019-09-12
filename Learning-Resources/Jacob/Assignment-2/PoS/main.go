package main

import (
	"fmt"
	"sync"
	// "bufio"
	// "crypto/sha256"
	// "encoding/hex"
	// "encoding/json"
	// "fmt"
	// "io"
	// "log"
	// "math/rand"
	// "net"
	// "os"
	// "strconv"
	// "sync"
	// "time"
	// "github.com/davecgh/go-spew/spew"
	// "github.com/joho/godotenv"
)

// Block represents each 'item' in the blockchain
type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
	Validator string
}

var Blockchain []Block                 // Confirmed/Validated blocks
var tempBlocks []Block                 //
var candidateBlocks = make(chan Block) // nodes that propose new blocks, send them to this channel

// announcements broadcasts winning validator to all nodes
var announcements = make(chan string)

var mutex = &sync.Mutex{} //prevents data races

// validators keeps track of open validators and balances
var validators = make(map[string]int)

func main() {
	fmt.Print("Proof of stake chain")
}
