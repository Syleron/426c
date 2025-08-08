package main

import (
    "errors"
    "github.com/labstack/gommon/log"
    "github.com/syleron/426c/common/packet"
    "net"
    "time"
)

type Client struct {
	Username string
    Conn net.Conn
    sendCh chan outboundMessage
    closed bool
}

type outboundMessage struct {
    cmdType int
    buf     []byte
}

func (c *Client) Send(cmdType int, buf []byte) (int, error) {
    log.Debug("Queueing message to " + c.Username)
    if c.sendCh == nil {
        // Fallback to direct write if no channel initialized
        _ = c.Conn.SetWriteDeadline(time.Now().Add(15 * time.Second))
        return c.Conn.Write(packet.PacketForm(byte(cmdType), buf))
    }
    select {
    case c.sendCh <- outboundMessage{cmdType: cmdType, buf: buf}:
        return 0, nil
    default:
        return 0, errors.New("send buffer full")
    }
}

func (c *Client) SendNotice(msg string) (int, error) {
    log.Debug("Queueing notice to " + c.Username)
    if c.sendCh == nil {
        _ = c.Conn.SetWriteDeadline(time.Now().Add(15 * time.Second))
        return c.Conn.Write(packet.PacketForm(byte(packet.SVR_NOTICE), []byte(msg)))
    }
    select {
    case c.sendCh <- outboundMessage{cmdType: int(packet.SVR_NOTICE), buf: []byte(msg)}:
        return 0, nil
    default:
        return 0, errors.New("send buffer full")
    }
}

func (c *Client) LoggedIn() bool {
	return c.Username != ""
}

func (c *Client) startWriter(done <-chan struct{}) {
    if c.sendCh == nil {
        c.sendCh = make(chan outboundMessage, 64)
    }
    go func() {
        for {
            select {
            case m, ok := <-c.sendCh:
                if !ok {
                    return
                }
                _ = c.Conn.SetWriteDeadline(time.Now().Add(15 * time.Second))
                _, _ = c.Conn.Write(packet.PacketForm(byte(m.cmdType), m.buf))
            case <-done:
                return
            }
        }
    }()
}

func (c *Client) stopWriter() {
    if c.sendCh != nil && !c.closed {
        close(c.sendCh)
        c.closed = true
    }
}
