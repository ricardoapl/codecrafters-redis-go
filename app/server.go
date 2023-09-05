package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
)

var dataset = make(map[string][]byte)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		r, err := Deserialize(bufio.NewReader(conn))
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading from client: ", err.Error())
			os.Exit(1)
		}

		command := r.Elements[0].Value
		args := [][]byte{}
		for _, element := range r.Elements[1:] {
			args = append(args, element.Value)
		}

		switch string(command) {
		case "ping":
			conn.Write([]byte("+PONG\r\n"))
		case "echo":
			data := bytes.Join(args, []byte(" "))
			conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(data), data)))
		case "set":
			key, val := string(args[0]), args[1]
			dataset[key] = val
			conn.Write([]byte("+OK\r\n"))
		case "get":
			key := string(args[0])
			val, ok := dataset[key]
			if ok {
				conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)))
			}
			conn.Write([]byte("$-1\r\n"))
		default:
			conn.Write([]byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", command)))
		}
	}
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}
