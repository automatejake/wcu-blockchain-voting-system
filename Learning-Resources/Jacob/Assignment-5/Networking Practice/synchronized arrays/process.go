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

/******* LISTENING PROCESS *******/
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
		fmt.Println("this far")

		go peerProcess(conn)
	}
}

/***************************************************************/

/*** PEER PROCESS ***/
func peerProcess(conn net.Conn) {

	defer fmt.Println("Client disconnected")
	defer conn.Close()

	//create peer object
	var peer Peer
	peer.Connection = conn
	peer.Port = 100

	//listen for new connection types
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}

		if string(buf[0:7]) == "connect" {
			fmt.Print("made it this far!")
			a := string(buf[8 : n-1])
			port, err := strconv.Atoi(a)
			if err == nil {
				Nodes[port] = true
				peer.Port = port
				Peers = append(Peers, peer)
			}
		} else if string(buf[0:9]) == "broadcast" {
			fmt.Println("recieved broadcast")
		}
	}

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
				currentPort := strconv.Itoa(port)
				conn, _ := net.Dial("tcp", "127.0.0.1:"+currentPort)
				if conn != nil {
					Nodes[port] = true
					conn.Write([]byte("connect " + currentPort))

					go peerProcess(conn)
				}

			}

			time.Sleep(100 * time.Millisecond)
		}
	}

}
