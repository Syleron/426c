package main

import (
	"github.com/labstack/gommon/log"
	"github.com/syleron/426c/common/packet"
	"net"
    "time"
)

type Client struct {
	Username string
	Conn net.Conn
}

func (c *Client) Send(cmdType int, buf []byte) (int, error) {
	log.Debug("Sending message to " + c.Username)
    // Set a write deadline to avoid hanging writes
    _ = c.Conn.SetWriteDeadline(time.Now().Add(15 * time.Second))
    return c.Conn.Write(packet.PacketForm(byte(cmdType), buf))
}

func (c *Client) SendNotice(msg string) (int, error) {
	log.Debug("Sending notice to " + c.Username)
    _ = c.Conn.SetWriteDeadline(time.Now().Add(15 * time.Second))
    return c.Conn.Write(packet.PacketForm(byte(packet.SVR_NOTICE), []byte(msg)))
}

func (c *Client) LoggedIn() bool {
	return c.Username != ""
}
