package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	fmt.Println("Hello World")
	ip := "0.0.0.0"
	port := "8080"
	address := fmt.Sprintf("%s:%s", ip, port)

	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("couldn't listen on address \"%s\", error: %s", address, err.Error())
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalf("couldn't accept conection: %s", err.Error())
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	request(conn)
	response(conn)
}

func request(conn net.Conn) {
	fmt.Println("request:")
	i := 0
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		if i == 0 {
			method := strings.Fields(line)[0]
			path := strings.Fields(line)[1]
			httpVersion := strings.Fields(line)[2]
			fmt.Printf("Method \"%s\", path \"%s\", http version \"%s\"\n", method, path, httpVersion)
		}
		fmt.Printf("> %s\n", line)
		if line == "" {
			break
		}
		i++
	}
}

func response(conn net.Conn) {
	// time.Sleep(5 * time.Second)

	body := `This Is Go Http
Server Using TCP
`

	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, body)
}
