package main

import (
	"github.com/syleron/426c/common/packet"
	"net"
)

type Client struct {
	Username string
	Conn net.Conn
}

func (c *Client) Send(cmdType int, buf []byte) (int, error) {
	return c.Conn.Write(packet.PacketForm(byte(cmdType), buf))
}

func (c *Client) SendNotice(msg string) (int, error) {
	return c.Conn.Write(packet.PacketForm(byte(packet.SVR_NOTICE), []byte(msg)))
}
