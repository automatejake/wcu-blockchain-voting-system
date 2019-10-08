package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Block struct {
	Index   int
	Message string
}

type Peer struct {
	Port       int
	Connection net.Conn
}

var Peers []Peer
var Nodes = make(map[int]bool)
var Blockchain []Block
var index int

/*************
*
*  Two primary threads with child processes:
*	1 - client, search for Clients on the network and listen in on them
*		a - each peer found opens a new listening process
*		  - each listening process adds new things to the chain
*	2 - server, listen for incoming connections
*		a - each client that connects launches a listening process
*		  - whenever, there is new data in the array, each listening process sends data to the listener
*
*************/

/******* SERVER PORTION *******/
func listenConnections() {
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

	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		fmt.Println(string(buf[0:]))
		_, err2 := conn.Write(buf[0:n]) //ONLY WRITING TO THE NODE THAT SENT THE DATA, NEED WAY TO BROADCAST TO ALL LISTENING NODES
		if err2 != nil {
			return
		}
	}
	// scanner := bufio.NewScanner(conn)

	// for scanner.Scan() {
	// 	io.WriteString(conn, "\nEnter a message to write to the block:  ")
	// 	message := scanner.Text()
	// 	var newBlock Block
	// 	newBlock.Index = index
	// 	newBlock.Message = string(message)
	// 	Blockchain = append(Blockchain, newBlock)
	// 	spew.Println(Blockchain)
	// 	index++
	// }

}

/***************************************************************/

/******* CLIENT PORTION *******/
func foundPeer(conn net.Conn, port int) {

	defer fmt.Println("Peer terminated process")
	defer closeConnection(port)

	Nodes[port] = true
	fmt.Println("found peer!", conn)

	//read incoming messages from Clients (new blocks)
	for {

	}

	// message, _ := bufio.NewReader(conn).ReadString('\n')
	// fmt.Println("Message from server: " + message)

}

func closeConnection(port int) {
	Nodes[port] = false
}

/***************************************************************/

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	ignore, _ := strconv.Atoi(os.Getenv("PORT"))
	Nodes[ignore] = true
	index = 0

	go listenConnections()

	//Discovering Clients, there are 65,535 ports on a computer
	//I am using ports 7000-7020
	for {
		for port := 7000; port <= 7020; port++ {

			if !Nodes[port] {

				conn, _ := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(port))
				if conn != nil {
					go foundPeer(conn, port)
				}

			}

			time.Sleep(100 * time.Millisecond)
		}
	}

}
