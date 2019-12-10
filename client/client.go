package main

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"os"
)

type Client struct {

}

func setupClient() {
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		panic(err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Rand = rand.Reader
	// connect to this socket
	conn, _ := tls.Dial("tcp", "127.0.0.1:8081", &config)
	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		// send to socket
		fmt.Fprintf(conn, text + "\n")
		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: "+message)
	}
}

func connectionHandler() {}

func msgReqShareKey() {}

func msgEncShareKey() {}

func msgSendShareKey() {}

func msgReqPubKey() {}

func msgSendPubKey() {}

func msgEncPubKey() {}

func ident() {}

func who() {}