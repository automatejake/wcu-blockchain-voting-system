package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
)

type Block struct {
	index     int
	randomNum int
}

var SyncBlockchain []Block
var index int

func handleConn(conn net.Conn) {
	newValue := rand.Intn(100)
	var tempBlock Block

	tempBlock.index = index
	tempBlock.randomNum = newValue

	fmt.Println("New value is ", newValue)
	SyncBlockchain = append(SyncBlockchain, tempBlock)
	fmt.Println(SyncBlockchain)
	index++

	conn.Write([]byte(string("\nhelodocmderefrf\n")))
	conn.Close()
}

func main() {
	index = 0
	port := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", port)
	checkError(err)

	server, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConn(conn)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
