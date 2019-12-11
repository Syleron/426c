package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	plib "github.com/syleron/426c/common/packet"
	"net"
	"os"
)

// Length of the user connected gives them currency
// Register user onto the network using their public key

// TODO - Registering user onto the network
// TODO - Creating groups (Permanent or not?) (Protected?)
// TODO - Distribution of "blocks"
// TODO - Charge blocks for sending a message
// TODO - Increase the cost of blocks depending on total number of spam (Calculate the rate of messaging for a particular room)

type Server struct {
	listener net.Listener
	clients map[string]*Client
}

func setupServer(laddr string) *Server {
	// Generate new tls keys
	// Read our certs
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		panic(err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Rand = rand.Reader
	// Setup our server listener
	listener, err := tls.Listen("tcp", laddr, &config)
	if err != nil {
		panic(err)
	}
	log.Printf("Listening on port %v\n", port)
	return &Server{
		listener: listener,
		clients: make(map[string]*Client),
	}
}

func (s *Server) connectionHandler() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Error(conn)
		}
		go s.newClient(conn)
	}
}

func (s *Server) newClient(conn net.Conn) {
	log.Print("New client connection")
	defer func() {
		// Remove our client
		s.clientRemoveByConnection(conn)
		// Close our connection
		if err := conn.Close(); err != nil {
			log.Error(err)
		}
	}()
	client := &Client{
		Conn:     conn,
	}
	br := bufio.NewReader(client.Conn)
	packet, err := plib.PacketRead(br)
	if (err != nil) || (packet[0] != plib.CMD_IDENT) {
		log.Error(err.Error())
		return
	}
	if ok := s.cmdIdent(client, packet[1:]); ok {
		log.Printf("new user - %s", client.Username)
	}
	for {
		packet, err = plib.PacketRead(br)
		if err != nil {
			log.Error(err)
			break
		}
		s.commandRouter(client, packet)
	}
}

func (s *Server) commandRouter(client *Client, packet []byte) {
	cmd := packet[0]
	switch {
	case cmd == plib.CMD_MSGALL:
		log.Print("Message all command")
		s.cmdMsgAll(client, packet[1:])
	case cmd == plib.CMD_MSGTO:
		log.Print("Message to command")
		s.cmdMsgTo(client, packet[1:])
	case cmd == plib.CMD_WHO:
		log.Print("Message who command")
		s.cmdWho(client)
	default:
		client.SendNotice("Unknown command")
	}
}

func (s *Server) broadcast(cmdType int, buf []byte) {
	for _, c := range s.clients {
		go c.Send(cmdType, buf)
	}
}

func (s *Server) cmdMsgAll(client *Client, packet []byte) (int, error) {
	// Make sure we have a valid username set
	if client.Username == "" {
		_, err := client.SendNotice("please register yourself with the server")
		if err != nil {
			log.Error(err)
		}
		return -1, nil
	}
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, uint16(len(client.Username)))
	buf.Write([]byte(client.Username))
	buf.Write(packet)
	// Broadcast the message to everyone
	s.broadcast(plib.SVR_MSG, buf.Bytes())
	// Return success
	return 0, nil
}

func (s *Server) cmdMsgTo(client *Client, packet []byte) (int, error) {
	// Make sure we have a valid username set
	if client.Username == "" {
		_, err := client.SendNotice("please register yourself with the server")
		if err != nil {
			log.Error(err)
		}
		return -1, nil
	}
	targetlen := int(binary.BigEndian.Uint16(packet[0:2]))
	target := string(packet[2:2+targetlen])
	data := packet[2+targetlen:]
	targetClient, exists := s.clients[target]
	if !exists {
		c, err := client.SendNotice(fmt.Sprintf("unknown recipient %s", target))
		if err != nil {
			return c, err
		}
		return -1, nil
	}
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, uint16(len(client.Username)))
	buf.Write([]byte(client.Username))
	buf.Write(data)
	return targetClient.Send(plib.SVR_MSG, buf.Bytes())
}

func (s *Server) cmdIdent(client *Client, packet []byte) bool {
	username := string(packet)
	// Make sure our we have a valid username length
	if len(username) > plib.MAX_NAME_LENGTH {
		_, err := client.SendNotice("username is too long")
		if err != nil {
			log.Error(err)
			return false
		}
	}
	// Add our new client to our map
	if err := s.clientAdd(username, client); err != nil {
		log.Error(err)
		return false
	}
	// Set our username for our client connection
	client.Username = username
	// Let everyone know that we have connected
	s.broadcast(plib.SVR_NOTICE, []byte(username + " connected"))
	return true
}

func (s *Server) cmdWho(client *Client) {
	log.Debug("WHO command requested")
	if client.Username == "" {
		_, err := client.SendNotice("please register yourself with the server")
		if err != nil {
			log.Error(err)
		}
		return
	}
	msg := fmt.Sprintf("Who (%v users):\n", len(s.clients))
	for username, _ := range s.clients {
		msg += fmt.Sprintf(" %v\n", username)
	}
	_, err := client.SendNotice(msg)
	if err != nil {
		log.Error(err)
	}
}

func (s *Server) clientAdd(username string, c *Client) error {
	log.Info("Adding client " + username)
	_, exists := s.clients[username]
	if exists {
		return errors.New("user already exists")
	}
	s.clients[username] = c
	return nil
}

func (s *Server) clientRemoveByUsername(username string) {
	log.Info("Removing client " + username)
	_, exists := s.clients[username]
	if exists {
		delete(s.clients, username)
	}
}

func (s *Server) clientRemoveByConnection(conn net.Conn) {
	for _, c := range s.clients {
		if c.Conn == conn {
			log.Info("Removing client " + c.Username)
			delete(s.clients, c.Username)
		}
	}
}

func (s *Server) shutdown() {
	log.Info("Server shutdown")
	if err := s.listener.Close(); err != nil {
		panic(err)
	}
	os.Exit(0)
}
