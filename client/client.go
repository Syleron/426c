package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"net"
)

type Client struct {
	Conn net.Conn
}

func setupClient() *Client {

	// Setup our listener
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		panic(err)
	}
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	config.Rand = rand.Reader
	// connect to this socket
	// TODO This should be a client command rather done automagically.
	conn, err := tls.Dial("tcp", "127.0.0.1:9000", &config)
	if err != nil {
		panic(err)
	}
	//ui.addInfoMessage("connected")
	return &Client{
		Conn: conn,
	}
}

func (c *Client) connectionHandler() {
	for  {
		fmt.Fprintf(c.Conn, "testing like a boss" + "\n")
	}
}

func (c *Client) msgReqShareKey() {}

func (c *Client) msgEncShareKey() {}

func (c *Client) msgSendShareKey() {}

func (c *Client) msgReqPubKey() {}

func (c *Client) msgSendPubKey() {}

func (c *Client) msgEncPubKey() {}

func (c *Client) ident() {}

func (c *Client) who() {}

func (c *Client) Close() {}
