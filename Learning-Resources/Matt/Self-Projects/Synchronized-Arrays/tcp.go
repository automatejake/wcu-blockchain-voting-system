package main

import (
	"fmt"
	"io"
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

var array []string

var Peers []Peer
var Ports = make(map[int]bool)

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

	var buf [512]byte

	for {
		//go readMessages(conn)

		length, err := conn.Read(buf[0:])

		if err != nil {
			return
		}

		if string(buf[0:7]) == "message" {
			tempString := string(buf[8 : length-1])
			array = append(array, tempString)

			go broadcast(tempString, conn)
		}
	}
}

func broadcast(tempString string, conn net.Conn) {
	bs := []byte(tempString)
	io.WriteString(conn, string(len(Peers)))
	for _, sockets := range Peers {
		sockets.Socket.Write(append([]byte("broadcast "), bs...))
		io.WriteString(conn, "Writing to "+string(sockets.Port)+" with bytes "+string(bs))
	}
}

func readMessages(conn net.Conn) {
	bs := make([]byte, 256)
	message, err := conn.Read(bs)
	if err == nil {
		array = append(array, string(message))
		io.WriteString(conn, string(len(array)))
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("File not loaded correctly.")
	}
	ignore, _ := strconv.Atoi(os.Getenv("PORT"))
	Ports[ignore] = true

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
					fmt.Println(string(len(Peers)))
					Ports[port] = true
					go handleConn(conn)
				}
			}

			time.Sleep(100 * time.Millisecond)
		}
	}
}
