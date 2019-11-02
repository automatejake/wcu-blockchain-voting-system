package main

import (
	// "bufio"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"io"
	"strings"
	"sync"
	"time"

	// "encoding/json"
	//"fmt"
	// "io"
	// "math/rand"
	"fmt"
	"net"
	"strconv"

	// "sync"
	// "time"
	// "github.com/davecgh/go-spew/spew"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Block struct {
	Index     int
	Timestamp string
	Vote      int
	Hash      string
	PrevHash  string
	Validator string
}

type Peer struct {
	Port   int
	Socket net.Conn
}

type Validator struct {
	Address string
	Tokens  int
}

var address string

var Blockchain []Block
var tempBlocks []Block
var Peers []Peer

var mutex = &sync.Mutex{}

var candidateBlocks = make(chan Block)
var announcements = make(chan string)

var dialedPorts = make(map[int]bool)
var validators = make(map[string]int)

func calculateHash(info string) string {
	sha := sha256.New()
	sha.Write([]byte(info))
	hash := sha.Sum(nil)

	return hex.EncodeToString(hash)
}

func calculateBlockHash(block Block) string {
	header := string(block.Index) + block.Timestamp + string(block.Vote) + block.PrevHash + block.Validator

	sha := sha256.New()
	sha.Write([]byte(header))
	hash := sha.Sum(nil)

	return hex.EncodeToString(hash)
}

func generateBlock(oldBlock Block, vote int, validator string) Block {
	var newBlock Block

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = time.Now().String()
	newBlock.Vote = vote
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Validator = validator
	newBlock.Hash = calculateBlockHash(newBlock)

	return newBlock
}

func isBlockValid(oldBlock Block, newBlock Block) bool {
	if newBlock.Index != oldBlock.Index+1 {
		return false
	} else if newBlock.PrevHash != oldBlock.Hash {
		return false
	} else if newBlock.Hash != calculateBlockHash(newBlock) {
		return false
	}

	return true
}

func listenConn() {
	portString := ":" + os.Getenv("PORT")
	listen, err := net.Listen("tcp", portString)
	if err != nil {
		log.Fatal(err)
	}

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	go func() {
		for {
			msg := <-announcements
			io.WriteString(conn, msg)
		}
	}()

	var buf [512]byte

	for {
		length, err := conn.Read(buf[0:])

		if err != nil {
			return
		}

		if string(buf[0:7]) == "connect" {
			port, err := strconv.Atoi(string(buf[8:length]))
			if err == nil {
				dialedPorts[port] = true
				tmpPeer := Peer{port, conn}
				Peers = append(Peers, tmpPeer)
			}
		} else if string(buf[0:7]) == "genesis" {
			genesisBlock := Block{0, time.Now().String(), 0, "", "", ""}
			genesisBlock.Hash = calculateBlockHash(genesisBlock)
			Blockchain = append(Blockchain, genesisBlock)
			fmt.Println(Blockchain)
			go broadcast(genesisBlock)
		} else if string(buf[0:8]) == "validate" {
			trimString := strings.TrimSuffix(string(buf[9:length]), "\n")
			tokens, err := strconv.Atoi(trimString)
			if err == nil {
				address = calculateHash(time.Now().String())
				validators[address] = tokens
				fmt.Println(validators)
				go updateValidators(Validator{address, tokens})
			}
		} else if string(buf[0:7]) == "propose" {
			trimString := strings.TrimSuffix(string(buf[8:length]), "\n")
			votes, err := strconv.Atoi(trimString)
			if err != nil {
				fmt.Println("Error")
				delete(validators, address)
				conn.Close()
			}

			prevBlock := Blockchain[len(Blockchain)-1]
			tempBlock := generateBlock(prevBlock, votes, address)

			if isBlockValid(prevBlock, tempBlock) {
				candidateBlocks <- tempBlock
				proposeCandidate(tempBlock)
				fmt.Println(tempBlocks)
			}
		} else if string(buf[0:15]) == "recieve propose" {
			buff := bytes.NewBuffer(buf[16:length])
			tmpStruct := new(Block)
			gobobj := gob.NewDecoder(buff)
			err := gobobj.Decode(tmpStruct)
			if err == nil {
				candidateBlocks <- *tmpStruct
			}
		} else if string(buf[0:16]) == "recieve validate" {
			buff := bytes.NewBuffer(buf[17:length])
			tmpStruct := new(Validator)
			gobobj := gob.NewDecoder(buff)
			err := gobobj.Decode(tmpStruct)
			if err == nil {
				validators[tmpStruct.Address] = tmpStruct.Tokens
				fmt.Println(validators)
			}
		} else if string(buf[0:9]) == "broadcast" {
			buff := bytes.NewBuffer(buf[10:length])
			tmpStruct := new(Block)
			gobobj := gob.NewDecoder(buff)
			err := gobobj.Decode(tmpStruct)
			if err == nil {
				Blockchain = append(Blockchain, *tmpStruct)
				fmt.Println(Blockchain)
			}
		}
	}
}

func broadcast(block Block) {
	buff := new(bytes.Buffer)
	tmpStruct := block
	gobobj := gob.NewEncoder(buff)
	err := gobobj.Encode(tmpStruct)
	if err == nil {
		for _, socket := range Peers {
			socket.Socket.Write(append([]byte("broadcast "), buff.Bytes()...))
		}
	}
}

func updateValidators(validator Validator) {
	buff := new(bytes.Buffer)
	tmpStruct := validator
	gobobj := gob.NewEncoder(buff)
	err := gobobj.Encode(tmpStruct)
	if err == nil {
		for _, socket := range Peers {
			fmt.Println("Sending Validators!")
			socket.Socket.Write(append([]byte("recieve validate "), buff.Bytes()...))
		}
	}
}

func proposeCandidate(block Block) {
	buff := new(bytes.Buffer)
	tmpStruct := block
	gobobj := gob.NewEncoder(buff)
	err := gobobj.Encode(tmpStruct)
	if err == nil {
		for _, socket := range Peers {
			fmt.Println("Sending Proposition!")
			socket.Socket.Write(append([]byte("recieve propose "), buff.Bytes()...))
		}
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file!")
	}

	ignore, _ := strconv.Atoi(os.Getenv("PORT"))
	dialedPorts[ignore] = true

	go listenConn()

	go func() {
		for candidate := range candidateBlocks {
			tempBlocks = append(tempBlocks, candidate)
			fmt.Println(tempBlocks)
		}
	}()

	for {
		for port := 9000; port <= 9001; port++ {
			if !dialedPorts[port] {
				fmt.Println("Checking...")
				sPort := strconv.Itoa(port)
				conn, _ := net.Dial("tcp", "127.0.0.1:"+sPort)
				if conn != nil {
					fmt.Println("Dial Successful!")
					tmpPeer := Peer{port, conn}
					Peers = append(Peers, tmpPeer)

					dialedPorts[port] = true

					conn.Write([]byte("connect " + strconv.Itoa(ignore)))
					go handleConn(conn)
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}
