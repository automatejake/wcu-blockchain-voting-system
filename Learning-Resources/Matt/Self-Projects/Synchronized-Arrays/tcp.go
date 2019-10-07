package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var array []int

func listenConn(env string) {
	port := ":" + os.Getenv(env)
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

	fmt.Println("Client Connected!")

	io.WriteString(conn, "Enter a number: ")
	input := bufio.NewScanner(conn)

	for input.Scan() {
		num, err := strconv.Atoi(input.Text())
		if err != nil {
			log.Println("The input was not a number.")
		}
		array = append(array, num)

		output, err := json.Marshal(array)
		if err != nil {
			log.Fatal(err)
		}

		io.WriteString(conn, "Successfully input to array.\n"+string(output))
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("File not loaded correctly.")
	}

	go listenConn("PORT_1")
	listenConn("PORT_2")
}
