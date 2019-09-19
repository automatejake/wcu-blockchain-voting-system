//https://medium.com/@mycoralhealth/part-2-networking-code-your-own-blockchain-in-less-than-200-lines-of-go-17fe1dad46e1

package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

type Block struct {
	Index      int
	Timestamp  string
	Message    string
	PrevHash   string
	Nonce      string
	Difficulty int
	Hash       string
}

// Blockchain is a series of validated Blocks
var Blockchain []Block

const difficulty = 4

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
	newBlock.Message = message //insert function to prompt user for messsage
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Index = oldBlock.Index + 1
	newBlock.Difficulty = difficulty

	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		if !isHashValid(calculateHash(newBlock), newBlock.Difficulty) {
			fmt.Println(calculateHash(newBlock), " do more work!")
			continue
		} else {
			fmt.Println(calculateHash(newBlock), " work done!")
			newBlock.Hash = calculateHash(newBlock)
			break
		}

	}
	return newBlock, nil
}

func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
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

func handleConn(conn net.Conn) {
	defer conn.Close()
	io.WriteString(conn, "Enter a message to the block:")

	scanner := bufio.NewScanner(conn)

	// take in BPM from stdin and add it to blockchain after conducting necessary validation
	go func() {
		for scanner.Scan() {
			message := scanner.Text()
			newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], message)
			if err != nil {
				log.Println(err)
				continue
			}
			if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
				newBlockchain := append(Blockchain, newBlock)
				replaceChain(newBlockchain)
			}

			bcServer <- Blockchain
			io.WriteString(conn, "\nEnter a message to write to the block:  ")
		}
	}()
	// simulate receiving broadcast
	go func() {
		for {
			time.Sleep(5 * time.Second)
			output, err := json.Marshal(Blockchain)
			if err != nil {
				log.Fatal(err)
			}
			io.WriteString(conn, string(output))
		}
	}()

	for _ = range bcServer {
		spew.Dump(Blockchain)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	bcServer = make(chan []Block)

	// Index     int
	// Timestamp string
	// Message   string
	// PrevHash  string
	// Nonce     string
	// Hash      string

	t := time.Now()
	genesisBlock := Block{0, t.String(), "Genesis Block", "", "", 0, ""}
	spew.Dump(genesisBlock) //pretty prints block before appending it
	Blockchain = append(Blockchain, genesisBlock)

	// start TCP and serve TCP server
	server, err := net.Listen("tcp", ":"+os.Getenv("ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}

}
