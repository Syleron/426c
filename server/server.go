package main

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/labstack/gommon/log"
	"github.com/syleron/426c/common/models"
	plib "github.com/syleron/426c/common/packet"
	"github.com/syleron/426c/common/utils"
	"net"
	"os"
	"time"
)

// Length of the user connected gives them currency
// Register user onto the network using their public key

// TODO - Creating groups (Permanent or not?) (Protected?)
// TODO - Distribution of "blocks"
// TODO - Charge blocks for sending a message
// TODO - Increase the cost of blocks depending on total number of spam (Calculate the rate of messaging for a particular room)
// TODO - Client/Server Version Validation
// TODO - Pending messages for offline people
// TODO - Prevent people from sending plain text
// TODO - Session timeout
// TODO - Make sure you cannot send a message to yourself
// TODO - Rate limit connection

// Username -> keys
// Store keys with server
//


// Total Chat's per second (TCS) / total number of users
// For example 1 / 10 = 0.10 *

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
	log.Infof("listening on port %v\n", port)
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
	log.Debug("New client connection")
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
	if err != nil {
		log.Error(err)
	}
	//// Handle initial request
	s.commandRouter(client, packet)
	// Handle subsequent requests
	for {
		packet, err = plib.PacketRead(br)
		if err != nil {
			log.Error(err)
			break
		}
		s.commandRouter(client, packet)
	}
}

func (s *Server) commandRouter(c *Client, p []byte) {
	if len(p) <= 0 {
		log.Error("invalid packet", p)
		return
	}
	switch p[0] {
	case plib.CMD_LOGIN:
		log.Debug("message login command")
		s.cmdLogin(c, p[1:])
	case plib.CMD_REGISTER:
		log.Debug("message register command")
		s.cmdRegister(c, p[1:])
	case plib.CMD_USER:
		log.Debug("message user command")
		s.authCheck(c)
		s.cmdUser(c, p[1:])
	case plib.CMD_MSGTO:
		log.Debug("message msg to command")
		s.authCheck(c)
		s.cmdMsgTo(c, p[1:])
	default:
		// TODO: This shouldn't ever happen. Perhaps handle block or handle this.
		log.Debug("received unknown command")
		c.Conn.Close()
	}
}

func (s *Server) cmdMsgTo(c *Client, p []byte) {
	var msgObj models.MsgToRequestModel
	if err := json.Unmarshal(p, &msgObj); err != nil {
		log.Debug("unable to unmarshal packet")
		return
	}
	// Make sure we have a username to send to
	if msgObj.To == "" {
		log.Debug("unable to send message. user not specified")
		return
	}
	if msgObj.ID == 0 {
		log.Debug("msg id is zero.. msgto failed")
		return
	}
	// Make sure our user is online
	if _, ok := s.clients[msgObj.To]; !ok {
		log.Debug("unable to send message as user is offline")
		return
	}
	// Set any data
	msgObj.From = c.Username
	msgObj.Date = time.Now()
	// Send our message to our recipient
	_, err := s.clients[msgObj.To].Send(plib.SVR_MSG, utils.MarshalResponse(&models.MsgResponseModel{
		Message: msgObj.Message,
	}))
	if err != nil {
		c.Send(plib.SVR_MSGTO, utils.MarshalResponse(&models.MsgToResponseModel{
			Success: false,
		}))
		return
	}
	// reply to our sender to say it was successful
	c.Send(plib.SVR_MSGTO, utils.MarshalResponse(&models.MsgToResponseModel{
		Success: true,
		MsgID: msgObj.ID,
		To: msgObj.To,
	}))
}

func (s *Server) cmdUser(c *Client, p []byte) {
	var userObj models.UserRequestModel
	if err := json.Unmarshal(p, &userObj); err != nil {
		log.Debug("unable to unmarshal packet")
		return
	}
	// Get our user from our users bucket
	user, err := userGet(userObj.Username)
	if err != nil {
		log.Error(err)
		c.Send(plib.SVR_USER, utils.MarshalResponse(&models.UserResponseModel{
			Success: false,
			Message: err.Error(),
		}))
		return
	}
	// Success, send details back
	c.Send(plib.SVR_USER, utils.MarshalResponse(&models.UserResponseModel{
		Success: true,
		Message: "user found",
		User:    models.User{
			Username:       user.Username,
			PubKey:         user.PubKey,
		},
	}))
}

func (s *Server) cmdRegister(c *Client, p []byte) {
	var registerObj models.RegisterRequestModel
	if err := json.Unmarshal(p, &registerObj); err != nil {
		log.Debug("unable to unmarshal packet")
		return
	}
	user := &models.User{
		Username:       registerObj.Username,
		PassHash:       registerObj.PassHash,
		EncPrivKey:     registerObj.EncPrivKey,
		PubKey:         registerObj.PubKey,
		RegisteredDate: time.Now(),
		Access:         0,
	}
	// Register our user
	if err := userAdd(user); err != nil {
		log.Debug(err)
	}
}

func (s *Server) cmdLogin(c *Client, p []byte) {
	var loginObj models.LoginRequestModel
	if err := json.Unmarshal(p, &loginObj); err != nil {
		log.Debug("unable to unmarshal packet")
		return
	}
	// Check version
	if loginObj.Version != VERSION {
		log.Debug("client version mismatch")
		c.Send(plib.SVR_LOGIN, utils.MarshalResponse(&models.LoginResponseModel{
			Success: false,
			Message: "Version mismatch. Please make sure you are running the latest version.",
		}))
		return
	}
	// Get our user account
	user, err := userGet(loginObj.Username)
	if err != nil {
		log.Debug("unable to find user account")
		c.Send(plib.SVR_LOGIN, utils.MarshalResponse(&models.LoginResponseModel{
			Success: false,
			Message: "Unable to find user account",
		}))
		return
	}
	// Compare credentials
	if user.PassHash != loginObj.Password {
		log.Debug("invalid login password")
		c.Send(plib.SVR_LOGIN, utils.MarshalResponse(&models.LoginResponseModel{
			Success: false,
			Message: "Invalid user account",
		}))
		return
	}
	// Set our connection details
	c.Username = loginObj.Username
	// If our user already is connected, disconnect them.
	if user, ok := s.clients[c.Username]; ok {
		// TODO: Should be a dedicated packet logging out the user.
		user.Conn.Close()
	}
	// Add client to our online list
	if err := s.clientAdd(c.Username, c); err != nil {
		log.Error(err)
	}
	c.Send(plib.SVR_LOGIN, utils.MarshalResponse(&models.LoginResponseModel{
		Username: loginObj.Username,
		Success: true,
		Message: "success",
	}))
}

func (s *Server) broadcast(cmdType int, buf []byte) {
	for _, c := range s.clients {
		go c.Send(cmdType, buf)
	}
}

func (s *Server) clientAdd(username string, c *Client) error {
	log.Debug("Adding client " + username)
	_, exists := s.clients[username]
	if exists {
		return errors.New("user already exists")
	}
	s.clients[username] = c
	return nil
}

func (s *Server) clientRemoveByUsername(username string) {
	log.Debug("Removing client " + username)
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
	log.Debug("Server shutdown")
	if err := s.listener.Close(); err != nil {
		panic(err)
	}
	os.Exit(0)
}

func (s *Server) authCheck(c *Client) {
	// Make sure we have a session set otherwise we kill their connection
	if c.Username == "" {
		log.Warn("Unauthorised. Connection closed. " + c.Conn.RemoteAddr().String())
		c.Conn.Close()
	}
}
