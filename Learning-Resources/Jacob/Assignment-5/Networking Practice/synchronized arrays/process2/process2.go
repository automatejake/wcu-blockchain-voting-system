package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

type Block struct {
	Index   int
	Message string
}

var Peers = make(map[int]bool)
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

/******* SERVER PORTION *******/
func listenConnections() {
	fmt.Println(os.Getenv("PORT"))
	port := ":" + os.Getenv("PORT")
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

	go broadcastChain(conn)

	for scanner.Scan() {
		io.WriteString(conn, "\nEnter a message to write to the block:  ")
		message := scanner.Text()
		var newBlock Block
		newBlock.Index = index
		newBlock.Message = string(message)
		SyncBlockchain = append(SyncBlockchain, newBlock)
		spew.Println(SyncBlockchain)
		index++
	}
}

func broadcastChain(conn net.Conn) {
	for {
		// conn.Write([]byte())
		time.Sleep(3 * time.Second)
	}
}

/***************************************************************/

/******* CLIENT PORTION *******/
func foundPeer(conn net.Conn, port int) {

	defer fmt.Println("Peer terminated process")

	Peers[port] = true
	fmt.Println("found peer!", conn)

	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Message from server: " + message)

	Peers[port] = false
}

/***************************************************************/

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	ignore, _ := strconv.Atoi(os.Getenv("PORT"))
	Peers[ignore] = true
	index = 0

	go listenConnections()

	//Discovering peers, there are 65,535 ports on a computer
	//I am using ports 7000-7020
	for {
		for port := 7000; port <= 7020; port++ {

			if !Peers[port] {

				conn, _ := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(port))
				if conn != nil {
					go foundPeer(conn, port)
				}

			}

			time.Sleep(100 * time.Millisecond)
		}
	}

}
