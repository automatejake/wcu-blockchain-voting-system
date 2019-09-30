package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
)

type Block struct {
	Index     int
	RandomNum int
}

var SyncBlockchain []Block
var index int

func handleConn(conn net.Conn) {
	result, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)

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

	conn.Write([]byte(string(output) + "\n"))

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
