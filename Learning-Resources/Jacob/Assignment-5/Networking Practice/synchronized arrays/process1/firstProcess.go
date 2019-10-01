package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
)

type Block struct {
	Index     int
	RandomNum int
}

var Clients []int
var SyncBlockchain []Block
var index int

func handleConn(conn net.Conn) {

	defer fmt.Println("Client disconnected")
	defer conn.Close()

	fmt.Println("Client connected")

	newValue := rand.Intn(100)
	var tempBlock Block

	tempBlock.Index = index
	tempBlock.RandomNum = newValue

	fmt.Println("New value is ", newValue)
	SyncBlockchain = append(SyncBlockchain, tempBlock)
	fmt.Println(SyncBlockchain)
	index++

	output, err := json.Marshal(SyncBlockchain)
	if err != nil {
		log.Fatal(err)
	}

	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		fmt.Println(string(buf[0:]))
		_, err2 := conn.Write(buf[0:n])
		if err2 != nil {
			return
		}
	}

	conn.Write([]byte(string(output) + "\n"))

	// conn.Close()
}

func main() {
	index = 0
	port := ":1200"

	server, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server listening on port", port)

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConn(conn)
	}
}
