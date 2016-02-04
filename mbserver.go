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

// Set the requested key to the requested value and respond with success
func handleSet(req gomemcached.MCRequest, conn net.Conn) {
	key := string(req.Key[:len(req.Key)])
	val := string(req.Body[:len(req.Body)])

	rwmutex.Lock()
	kvmap[key] = val

	res := gomemcached.MCResponse{
		Opcode: gomemcached.SET,
		Status: gomemcached.SUCCESS,
		Opaque: 0,
		Cas:    0,
		Extras: []byte{},
		Key:    []byte{},
		Body:   []byte{},
	}
	conn.Write(res.Bytes())
	rwmutex.Unlock()

	conn.Close()
}

// --------------------------------------------------------------------------

func handleGet(req gomemcached.MCRequest, conn net.Conn) {
	key := string(req.Key[:len(req.Key)])

	rwmutex.RLock()
	val := kvmap[key]

	//flags := []byte{0xde, 0xad, 0xbe, 0xef}
	flags := []byte{0x00, 0x00, 0x00, 0x00}
	var res gomemcached.MCResponse
	if val != "" {
		res = gomemcached.MCResponse{
			Opcode: gomemcached.GET,
			Status: gomemcached.SUCCESS,
			Opaque: 0,
			Cas:    1,
			Extras: []byte(flags),
			Key:    []byte{},
			Body:   []byte(val),
		}
	} else {
		res = gomemcached.MCResponse{
			Opcode: gomemcached.GET,
			Status: gomemcached.KEY_ENOENT,
			Opaque: 0,
			Cas:    0,
			Extras: []byte{},
			Key:    []byte{},
			Body:   []byte("Not found"),
		}
	}
	fmt.Println(res.String())
	conn.Write(res.Bytes())
	rwmutex.RUnlock()

	conn.Close()
}

// --------------------------------------------------------------------------

// Handles incoming GET and SET requests.
func handleRequest(conn net.Conn) {

	req := gomemcached.MCRequest{}
	_, err := req.Receive(bufio.NewReader(conn), nil)
	if err != nil {
		fmt.Println("Error receiving message:", err.Error())
		conn.Close()
		return
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
