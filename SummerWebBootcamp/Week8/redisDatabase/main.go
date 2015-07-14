package main

import (
	"bufio"
	"io"
	"net"
	"strings"
)

struct

func handleDatabase() {
	var db = map[string]string{}

}

func handleConn(conn net.Conn) {
	defer conn.Close()

	scn := bufio.NewScanner(conn)
	for scn.Scan() {
		line := scn.Text()
		lines := strings.Split(line, " ")
		switch lines[0] {
		case "GET":
			key := strings.Join(lines[1:], " ")
			data := db[key]
			io.WriteString(conn, data)
		case "SET":
		case "DEL":
		default:
			io.WriteString(conn, "Unknown command: "+line)
		}
	}
}

func main() {
	server, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			panic(err)
		}

		go handleConn(conn)
	}
}
