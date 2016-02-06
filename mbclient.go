package main

import (
	"bufio"
	"fmt"
	"github.com/dustin/gomemcached"
	"net"
	"os"
	"sync"
)

var wg sync.WaitGroup

func sendSet(key, value string) {
	req := gomemcached.MCRequest{
		Opcode:  gomemcached.SET,
		Cas:     938424885,
		Opaque:  7242,
		VBucket: 824,
		Extras:  []byte{},
		Key:     []byte(key),
		Body:    []byte(value),
	}

	conn, _ := net.Dial("tcp", "localhost:9955")

	conn.Write(req.Bytes())

	res := gomemcached.MCResponse{}
	_, err := res.Receive(bufio.NewReader(conn), nil)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	fmt.Println(res.String())
	conn.Close()

}

func sendSet(key string, value int) {
	req := gomemcached.MCRequest{
		Opcode:  gomemcached.SET,
		Cas:     938424885,
		Opaque:  7242,
		VBucket: 824,
		Extras:  []byte{},
		Key:     []byte(key),
		Body:    []byte(value),
	}

	conn, _ := net.Dial("tcp", "localhost:9955")

	conn.Write(req.Bytes())

	res := gomemcached.MCResponse{}
	_, err := res.Receive(bufio.NewReader(conn), nil)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	fmt.Println(res.String())
	conn.Close()

}

func sendGet(key string) {
	req := gomemcached.MCRequest{
		Opcode:  gomemcached.GET,
		Cas:     938424885,
		Opaque:  7242,
		VBucket: 824,
		Extras:  []byte{},
		Key:     []byte(key),
		Body:    []byte{},
	}

	conn, _ := net.Dial("tcp", "localhost:9955")

	conn.Write(req.Bytes())

	res := gomemcached.MCResponse{}
	_, err := res.Receive(bufio.NewReader(conn), nil)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	fmt.Println(res.String())
	value := string(res.Body[:len(res.Body)])
	fmt.Println(value)
	conn.Close()

}

func sendUnknownCommand() {
	req := gomemcached.MCRequest{
		Opcode:  gomemcached.ADD,
		Cas:     938424885,
		Opaque:  7242,
		VBucket: 824,
		Extras:  []byte{},
		Key:     []byte("key"),
		Body:    []byte("somevalue"),
	}

	conn, _ := net.Dial("tcp", "localhost:9955")

	conn.Write(req.Bytes())

	res := gomemcached.MCResponse{}
	_, err := res.Receive(bufio.NewReader(conn), nil)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	fmt.Println(res.String())
	conn.Close()

}

func sendMalformedCommand() {
	req := "HERE ARE SOME BYTES TO SEE IF YOUR SERVER CAN HANDLE THEM!"

	conn, _ := net.Dial("tcp", "localhost:9955")

	conn.Write([]byte(req))

	res := gomemcached.MCResponse{}
	_, err := res.Receive(bufio.NewReader(conn), nil)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	fmt.Println(res.String())
	conn.Close()

}

func client() {
	sendUnknownCommand()
	sendMalformedCommand()
	sendSet("lynsey", 12345)
	sendGet("lynsey")
	sendSet("lynsey", "oh hello, oh world!")
	sendGet("lynsey")
	wg.Done()
}

func main() {

	//
	numClients := 1

	wg.Add(numClients)

	for i := 0; i < numClients; i++ {
		client()
	}

	fmt.Fprintf(os.Stderr, "Waiting for clients...\n")
	wg.Wait()
	fmt.Println("All clients finished!")

}
