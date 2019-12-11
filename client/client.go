package main

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	plib "github.com/syleron/426c/common/packet"
	"net"
)

type Client struct {
	Reader     *bufio.Reader
	Writer     *bufio.Writer
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
	return &Client{
		Writer: bufio.NewWriter(conn),
		Reader: bufio.NewReader(conn),
		Conn: conn,
	}
}

func (c *Client) Send(cmdType int, buf []byte) (int, error) {
	return c.Conn.Write(plib.PacketForm(byte(cmdType), buf))
}

func (c *Client) connectionHandler() {
	c.Send(plib.CMD_IDENT, []byte("testing"))
	c.Send(plib.CMD_MSGALL, []byte("more please"))
	c.Send(plib.CMD_WHO, []byte("more please"))
	for  {
		line, err := plib.PacketRead(c.Reader)
		if err != nil {
			break
		}
		fmt.Println(string(line[1:]))
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
