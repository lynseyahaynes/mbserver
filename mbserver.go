/* A basic binary memcache server

Run as: go run mbserver.go <ip>:<port>

Written by : Lynsey Haynes */

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/dustin/gomemcached"
)

var rwmutex sync.RWMutex

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}
}

func handleUnknown(conn net.Conn) {
	fmt.Println("received unknown message!")
}

func handleSet(req gomemcached.MCRequest, conn net.Conn) {
	rwmutex.Lock()
	fmt.Println("inside the write lock!")
	rwmutex.Unlock()

}

func handleGet(req gomemcached.MCRequest, conn net.Conn) {
	rwmutex.RLock()
	fmt.Println("inside the read lock!")
	rwmutex.RUnlock()
}

// Handles incoming GET and SET requests.
func handleRequest(conn net.Conn) {

	req := gomemcached.MCRequest{}
	_, err := req.Receive(bufio.NewReader(conn), nil)
	if err != nil {
		fmt.Println("Error receiving message:", err.Error())
	}

	if req.Opcode == gomemcached.GET {
		handleGet(req, conn)
	} else if req.Opcode == gomemcached.SET {
		handleSet(req, conn)
	} else {
		handleUnknown(conn)
	}

	fmt.Println(req)
}

//Listens for connections from possibly concurrent clients
func listenForConnections(addr string) {
	server, err := net.Listen("tcp", addr)
	checkError(err)
	defer server.Close()

	for {
		conn, err := server.Accept()
		checkError(err)
		go handleRequest(conn)
	}
}

// Where the magic happens
func main() {

	addr := os.Args[1]
	listenForConnections(addr)

}
