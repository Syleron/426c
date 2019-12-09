package main

import "net"

type Client struct {
	Username string
	Conn net.Conn
}

func (c *Client) Send(cmdType int, buf []byte) (int, error) {
	return c.Conn.Write(packetForm(byte(cmdType), buf))
}

func (c *Client) SendNotice(msg string) (int, error) {
	return c.Conn.Write(packetForm(byte(SVR_NOTICE), []byte(msg)))
}
