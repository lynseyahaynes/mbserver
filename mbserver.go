/* A basic binary memcache server

Run as: go run mbserver.go <ip>:<port>

Written by : Lynsey Haynes */

package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"sync"
)

const (
	HEADER_SIZE = 24
	REQUEST     = 0x80
	RESPONSE    = 0x81
	GET         = 0x00
	SET         = 0x01
	NOOP        = 0x0A
	UNKNOWN     = 0x81
	SUCCESS     = 0x0000
	NOTFOUND    = 0x0001
	MAX_INT     = 2147483647
)

type mbpacket struct {
}

type data struct {
	val  []byte
	kind []byte
}

var rwmutex sync.RWMutex
var kvmap map[string]data

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}
}

// --------------------------------------------------------------------------

func constructPacket(opcode uint8, extras []byte, status uint8, body []byte) []byte {

	var magic []byte = []byte{0x81}
	var opaque []byte = []byte{0x00, 0x00, 0x00, 0x00}
	var CAS []byte = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	var keylen []byte = []byte{0x00, 0x00}
	var bodylen = make([]byte, 4)
	var extralen = uint8(len(extras))
	binary.BigEndian.PutUint32(bodylen, uint32(len(body)+len(extras)))

	msg := append(magic, []byte{opcode}...)
	msg = append(msg, keylen...)
	msg = append(msg, []byte{extralen}...)
	msg = append(msg, []byte{0x00}...)
	msg = append(msg, []byte{0x00, status}...)
	msg = append(msg, bodylen...)
	msg = append(msg, opaque...)
	msg = append(msg, CAS...)
	msg = append(msg, extras...)
	msg = append(msg, body...)

	return msg
}

// --------------------------------------------------------------------------

// Send an error response to the connection and close it
func handleUnknown(conn *net.TCPConn) {
	var empty []byte
	msg := constructPacket(NOOP, empty, UNKNOWN, empty)

	fmt.Printf("%X", msg)
	fmt.Println(" ", len(msg))
	conn.Write(msg)

}

// --------------------------------------------------------------------------

// Set the requested key to the requested value and respond with success
func handleSet(header []byte, conn *net.TCPConn) {
	var keylen uint16
	var extralen uint8
	var bodylen uint32
	keylen = binary.BigEndian.Uint16(header[2:4])
	extralen = header[4]
	bodylen = binary.BigEndian.Uint32(header[8:12])

	if bodylen < uint32(keylen)+uint32(extralen) {
		fmt.Printf("Error! Body length not long enough. Header: %X\n", header)
		fmt.Println("BodyLength = ", bodylen)
		return
	}

	body := make([]byte, bodylen)
	if bodylen <= MAX_INT {
		conn.SetReadBuffer(int(bodylen))
		conn.Read(body)
	} else { //size of body is larger than max ReadBuffer size
		restlen := int(bodylen - MAX_INT)
		rest := make([]byte, restlen)
		conn.SetReadBuffer(MAX_INT)
		conn.Read(body)
		conn.SetReadBuffer(restlen)
		conn.Read(rest)
		body = append(body, rest...)
	}

	extras := body[:extralen-4]
	key := string(body[extralen:(keylen + uint16(extralen))])
	val := body[(uint16(extralen) + keylen):bodylen]
	rwmutex.Lock()
	kvmap[key] = data{val, extras}

	var empty []byte
	msg := constructPacket(SET, empty, SUCCESS, empty)

	conn.Write(msg)
	rwmutex.Unlock()

}

// --------------------------------------------------------------------------

// handle GET request. if value is there, return value, otherwise 'not found'
func handleGet(header []byte, conn *net.TCPConn) {
	var keylen uint16
	keylen = binary.BigEndian.Uint16(header[2:4])
	keybuf := make([]byte, int(keylen))
	conn.SetReadBuffer(int(keylen))
	conn.Read(keybuf)
	key := string(keybuf)

	// ===============
	rwmutex.RLock()
	val := kvmap[key].val
	kind := kvmap[key].kind

	var msg []byte
	if val != nil {
		msg = constructPacket(GET, kind, SUCCESS, val)
	} else {
		var empty []byte
		msg = constructPacket(GET, empty, NOTFOUND, []byte("Not found"))
	}
	conn.Write(msg)
	rwmutex.RUnlock()

}

// --------------------------------------------------------------------------

// Handles incoming GET, SET, and UNKNOWN requests based on the packet header.
func handleRequest(conn *net.TCPConn) {

	header := make([]byte, 24)
	conn.SetReadBuffer(24)
	n, err := conn.Read(header)
	if err != nil || n != 24 {
		return
	}

	magic := header[0]
	if magic != REQUEST || n < HEADER_SIZE {
		handleUnknown(conn)
		return
	}

	opcode := header[1]
	if opcode == GET {
		handleGet(header, conn)
	} else if opcode == SET {
		handleSet(header, conn)
	} else {
		handleUnknown(conn)
	}
}

// --------------------------------------------------------------------------

func handleRequests(conn *net.TCPConn) {
	defer conn.Close()
	for {
		handleRequest(conn)
	}
}

// --------------------------------------------------------------------------

//Listens for connections from possibly concurrent clients
func listenForConnections(addr string) {
	tcpaddr, _ := net.ResolveTCPAddr("tcp", addr)
	server, err := net.ListenTCP("tcp", tcpaddr)
	checkError(err)

	for {
		conn, err := server.AcceptTCP()
		checkError(err)
		go handleRequests(conn)
	}
}

// --------------------------------------------------------------------------

// Where the magic happens
func main() {

	kvmap = make(map[string]data)
	addr := os.Args[1]
	listenForConnections(addr)

}
