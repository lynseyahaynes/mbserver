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
var kvmap map[string]string

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}
}

// --------------------------------------------------------------------------

// Send an error response to the connection and close it
func handleUnknown(conn net.Conn) {
	fmt.Println("received unknown message!")
	res := gomemcached.MCResponse{
		Opcode: gomemcached.NOOP,
		Status: gomemcached.UNKNOWN_COMMAND,
		Opaque: 0,
		Cas:    0,
		Extras: []byte(""),
		Key:    []byte(""),
	}

	conn.Write(res.Bytes())
	conn.Close()
}

// --------------------------------------------------------------------------

func handleSet(req gomemcached.MCRequest, conn net.Conn) {
	rwmutex.Lock()
	fmt.Println("inside the write lock!")
	key := string(req.Key[:len(req.Key)])
	val := string(req.Body[:len(req.Body)])
	kvmap[key] = val
	res := gomemcached.MCResponse{
		Opcode: gomemcached.SET,
		Status: gomemcached.SUCCESS,
		Opaque: 0,
		Cas:    0,
		Extras: []byte{0},
		Key:    req.Key,
		Body:   []byte{0},
	}
	conn.Write(res.Bytes())
	rwmutex.Unlock()

}

// --------------------------------------------------------------------------

func handleGet(req gomemcached.MCRequest, conn net.Conn) {
	rwmutex.RLock()
	key := string(req.Key[:len(req.Key)])
	val := kvmap[key]

	var res gomemcached.MCResponse
	if val != "" {
		res = gomemcached.MCResponse{
			Opcode: gomemcached.GET,
			Status: gomemcached.SUCCESS,
			Opaque: 0,
			Cas:    0,
			Extras: []byte{0},
			Key:    req.Key,
			Body:   []byte(val),
		}
	} else {
		res = gomemcached.MCResponse{
			Opcode: gomemcached.GET,
			Status: gomemcached.KEY_ENOENT,
			Opaque: 0,
			Cas:    0,
			Extras: []byte{0},
			Key:    req.Key,
			Body:   []byte("Not found"),
		}
	}
	conn.Write(res.Bytes())

	rwmutex.RUnlock()
}

// --------------------------------------------------------------------------

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

// --------------------------------------------------------------------------

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

// --------------------------------------------------------------------------

// Where the magic happens
func main() {

	kvmap = make(map[string]string)
	addr := os.Args[1]
	listenForConnections(addr)

}
