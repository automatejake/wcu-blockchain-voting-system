package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"strings"

	//"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Peer struct {
	Port   int
	Socket net.Conn
}

type Block struct {
	Index     int
	Timestamp string
	Vote      int
	Hash      string
	PrevHash  string
}

var index int

var Blockchain []Block

var Peers []Peer
var Ports = make(map[int]bool)

func calculateHash(block Block) string {
	header := string(block.Index) + block.Timestamp + string(block.Vote) + block.PrevHash

	sha := sha256.New()
	sha.Write([]byte(header))
	hash := sha.Sum(nil)

	return hex.EncodeToString(hash)
}

func isBlockValid(oldBlock Block, newBlock Block) bool {
	if newBlock.PrevHash != calculateHash(oldBlock) {
		return false
	} else if newBlock.Index != (oldBlock.Index + 1) {
		return false
	} else if newBlock.Hash != calculateHash(newBlock) {
		return false
	}

	return true
}

func generateBlock(oldBlock Block, vote int) Block {
	var block Block

	block.Index = oldBlock.Index + 1
	block.Timestamp = time.Now().String()
	block.Vote = vote
	block.PrevHash = oldBlock.Hash
	block.Hash = calculateHash(block)

	return block
}

func listenConn() {
	port := ":" + os.Getenv("PORT")
	server, err := net.Listen("tcp", string(port))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listening on" + string(port))
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	var buf [1024]byte

	for {
		length, err := conn.Read(buf[0:])

		if err != nil {
			return
		}

		if string(buf[0:7]) == "connect" {
			port, err := strconv.Atoi(string(buf[8:length]))
			if err == nil {
				Ports[port] = true
				peer := Peer{port, conn}
				Peers = append(Peers, peer)
			}
		} else if string(buf[0:7]) == "message" {
			trimString := strings.TrimSuffix(string(buf[8:length]), "\n")
			tempVote, err := strconv.Atoi(trimString)
			if err == nil {
				chainLength := len(Blockchain)
				tempBlock := generateBlock(Blockchain[chainLength-1], tempVote)
				if isBlockValid(Blockchain[chainLength-1], tempBlock) {
					fmt.Println("Block Valid!")
					Blockchain = append(Blockchain, tempBlock)
					go broadcast(tempBlock)
				}
			}
		} else if string(buf[0:9]) == "broadcast" {
			tempBuff := bytes.NewBuffer(buf[10:length])
			tempStruct := new(Block)
			gobobj := gob.NewDecoder(tempBuff)
			err := gobobj.Decode(tempStruct)
			if err != nil {
				fmt.Println(err)
			}

			Blockchain = append(Blockchain, *tempStruct)

			fmt.Println(Blockchain)
		}
	}
}

func broadcast(tempBlock Block) {
	tempBuff := new(bytes.Buffer)
	tempStruct := tempBlock
	gobobj := gob.NewEncoder(tempBuff)
	err := gobobj.Encode(tempStruct)
	if err == nil {
		for _, sockets := range Peers {
			fmt.Println(tempStruct)
			sockets.Socket.Write(append([]byte("broadcast "), tempBuff.Bytes()...))
		}
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("File not loaded correctly.")
	}
	ignore, _ := strconv.Atoi(os.Getenv("PORT"))
	Ports[ignore] = true
	index = 0

	genesisBlock := Block{0, time.Now().String(), 0, "", ""}
	genesisBlock.Hash = calculateHash(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	go listenConn()

	for {
		for port := 9000; port <= 9001; port++ {
			if !Ports[port] {
				fmt.Println("Checking...")
				currentPort := strconv.Itoa(port)
				conn, _ := net.Dial("tcp", "127.0.0.1:"+currentPort)
				if conn != nil {
					fmt.Println("Found Peer!")
					peer := Peer{port, conn}
					Peers = append(Peers, peer)
					Ports[port] = true
					go handleConn(conn)
				}
			}

			time.Sleep(100 * time.Millisecond)
		}
	}
}
