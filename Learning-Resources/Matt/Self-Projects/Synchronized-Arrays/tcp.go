package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var array []int

func handleConn(conn net.Conn) {
	defer conn.Close()

	io.WriteString(conn, "Enter a number: ")
	input := bufio.NewScanner(conn)

	for input.Scan() {
		inputNum, err := strconv.Atoi(input.Text())
		if err != nil {
			log.Printf("%v was not a number. %v", input.Text(), err)
		}

		array = append(array, inputNum)

		output, err := json.Marshal(array)
		if err != nil {
			log.Fatal(err)
		}

		io.WriteString(conn, string(output))
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	server, err := net.Listen("tcp", ":"+os.Getenv("ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}
