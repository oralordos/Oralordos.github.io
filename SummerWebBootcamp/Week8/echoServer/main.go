package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Connection accepted from:", conn.RemoteAddr())
	for {
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err == io.EOF {
			if n == 0 {
				break
			}
		} else if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Printf("Got data from \"%s\" with a payload of: %s", conn.RemoteAddr(), data[:n])
		conn.Write(data[:n])
	}
	fmt.Println("Connection closed at:", conn.RemoteAddr())
}

func main() {
	server, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatalln(err.Error())
		}
		go handleConn(conn)
	}
}
