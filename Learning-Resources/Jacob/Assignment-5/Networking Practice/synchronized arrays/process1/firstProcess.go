package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/davecgh/go-spew/spew"
)

type Block struct {
	Index   int
	Message string
}

var Clients []int
var SyncBlockchain []Block
var index int

/*************
*
*  Two primary threads with child processes:
*	1 - client, search for peers on the network and listen in on them
*		a - each peer found opens a new listening process
*		  - each listening process adds new things to the chain
*	2 - server, listen for incoming connections
*		a - each client that connects launches a listening process
*		  - whenever, there is new data in the array, each listening process sends data to the listener
*
*************/

func searchPeers() {
	for {

	}
}

func listenConnections() {
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

func handleConn(conn net.Conn) {

	defer fmt.Println("Client disconnected")
	defer conn.Close()

	fmt.Println("Client connected")
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Text()
		var newBlock Block
		newBlock.Index = index
		newBlock.Message = string(message)
		SyncBlockchain = append(SyncBlockchain, newBlock)
		spew.Println(SyncBlockchain)
		index++

		io.WriteString(conn, "\nEnter a message to write to the block:  ")
	}
}

func main() {
	index = 0

	go listenConnections()
	go searchPeers()

}
