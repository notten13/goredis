package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	listener, error := net.Listen("tcp", ":9876")

	if error != nil {
		log.Fatal(error)
	}

	fmt.Println("Server ready on port 9876")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the incoming data
	fmt.Printf("Received: %s", buf)
}
