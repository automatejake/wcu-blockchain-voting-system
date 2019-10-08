package main

import (
	"fmt"
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

type Peer struct {
	Port   int
	Socket net.Conn
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
*	2 - server, listen for incoming Sockets
*		a - each client that connects launches a listening process
*		  - whenever, there is new data in the array, each listening process sends data to the listener
*
*************/

/******* LISTENING PROCESS *******/
func listenSockets() {
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

		go peerProcess(conn)
	}
}

/***************************************************************/

/*** PEER PROCESS ***/
func peerProcess(conn net.Conn) {

	defer fmt.Println("TCP Connection Ended")
	defer conn.Close()

	fmt.Println("New TCP Connection")

	//listen for new connection types
	var buf [512]byte

	for {

		msgLength, err := conn.Read(buf[0:])

		if err != nil {
			return
		}

		if string(buf[0:7]) == "connect" {

			port, err := strconv.Atoi(string(buf[8:msgLength]))
			if err == nil {
				Nodes[port] = true
				peer := Peer{port, conn}
				Peers = append(Peers, peer)
			}

		} else if string(buf[0:9]) == "syncChain" {

			fmt.Println("syncing chain ", msgLength)

		} else if string(buf[0:7]) == "message" {
			// after connecting via "nc localhost [port]"
			// write to chain with "message [message inserted here]"
			tempBlock := Block{index, string(buf[8 : msgLength-1])}
			Blockchain = append(Blockchain, tempBlock)
			index++

			//broadcast message to all other connected nodes
			spew.Println(Blockchain)
			go broadcast(tempBlock)

		} else if string(buf[0:9]) == "broadcast" {

			fmt.Println(string(buf[10:msgLength]))

		} else {
			fmt.Println(msgLength, string(buf[0:10]), buf[0:10])
		}
	}

}

func broadcast(tempBlock Block) {
	for _, element := range Peers {
		element.Socket.Write([]byte("broadcast hi"))
	}
}

/***************************************************************/

/****** Discovery Process ******/
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	ignore, _ := strconv.Atoi(os.Getenv("PORT"))
	Nodes[ignore] = true

	index = 0

	go listenSockets()

	//Discovering Clients, there are 65,535 ports on a computer
	//I am using ports 7000-7020
	for {
		for port := 7000; port <= 7020; port++ {

			if !Nodes[port] {
				currentPort := strconv.Itoa(port)
				conn, _ := net.Dial("tcp", "127.0.0.1:"+currentPort)
				if conn != nil {
					peer := Peer{port, conn}
					Peers = append(Peers, peer)
					Nodes[port] = true
					conn.Write([]byte("connect " + strconv.Itoa(ignore)))
					go peerProcess(conn)
				}

			}

			time.Sleep(100 * time.Millisecond)
		}
	}

}

/******************************************************************************/
