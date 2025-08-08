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
	"strings"
	"time"
	"sync"
)

// Length of the user connected gives them currency
// Register user onto the network using their public key

// TODO - Creating groups (Permanent or not?) (Protected?)
// TODO - Charge blocks for sending a message
// TODO - Increase the cost of blocks depending on total number of spam (Calculate the rate of messaging for a particular room)
// TODO - Prevent people from sending plain text
// TODO - Make sure you cannot send a message to yourself
// TODO - Rate limit connection
// TODO - Disable registration

// Username -> keys
// Store keys with server
//

type Server struct {
	listener net.Listener
	clients  map[string]*Client
    mu       sync.RWMutex
    closing  bool
    shutdownCh chan struct{}
    wg       sync.WaitGroup
    queue    *Queue
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
    // Initialize metrics and start metrics server (fixed port for now)
    initMetrics()
    startMetricsServer(":2112")
    s := &Server{
		listener: listener,
		clients:  make(map[string]*Client),
        shutdownCh: make(chan struct{}),
	}
    // Initialize message queue
    s.queue = newQueue(s)
    // Setup block distributor
    go blockDistribute(s)
	return s
}

func (s *Server) connectionHandler() {
    // Periodic log of connected users
    go func() {
        ticker := time.NewTicker(10 * time.Second)
        defer ticker.Stop()
        for {
            select {
            case <-ticker.C:
                s.mu.RLock()
                connected := len(s.clients)
                s.mu.RUnlock()
                log.Infof("connected users: %d", connected)
            case <-s.shutdownCh:
                return
            }
        }
    }()
    for {
		conn, err := s.listener.Accept()
		if err != nil {
            if s.isClosing() {
                break
            }
            log.Error(err)
            continue
		}
        metricConnectionsTotal.Inc()
        go s.newClient(conn)
	}
}

func (s *Server) newClient(conn net.Conn) {
    s.wg.Add(1)
    defer s.wg.Done()
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
		Conn: conn,
	}
	br := bufio.NewReader(client.Conn)
    // Set initial read deadline
    _ = client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	packet, err := plib.PacketRead(br)
	if err != nil {
        conn.Close()
        return
	}
	//// Handle initial request
	s.commandRouter(client, packet)
	// Handle subsequent requests
	for {
        // Refresh read deadline for each read
        _ = client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		packet, err = plib.PacketRead(br)
		if err != nil {
			break
		}
		s.commandRouter(client, packet)
	}
}

func (s *Server) commandRouter(c *Client, p []byte) {
	if len(p) <= 0 {
		log.Error("invalid packet -> ", string(p))
		c.Conn.Close()
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
        if !s.authCheck(c) { return }
		s.cmdUser(c, p[1:])
	case plib.CMD_MSGTO:
		log.Debug("message msg to command")
        if !s.authCheck(c) { return }
		s.cmdMsgTo(c, p[1:])
	default:
		log.Warn("received unknown command -> ", string(p))
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
    // Make sure our recipient user exists
    _, err := dbUserGet(msgObj.To)
	if err != nil {
		c.Send(plib.SVR_USER, utils.MarshalResponse(&models.UserResponseModel{
			Success: false,
			Message: err.Error(),
		}))
		return
	}
	// Debit our blocks
	totalBlocks, err := dbUserBlockDebit(c.Username, blockCalcCost())
	if err != nil {
		log.Debug("user has insufficient funds")
		// TODO: This should return a blocks error
		c.Send(plib.SVR_MSGTO, utils.MarshalResponse(&models.MsgToResponseModel{
			Success: false,
			MsgID:   msgObj.ID,
			To:      msgObj.To,
		}))
		return
	}
	// Set any data
	msgObj.From = c.Username
	msgObj.Date = time.Now()
    // Enqueue message for reliable delivery (online/offline)
    // The queue will notify sender upon successful delivery.
    s.queue.Add(&msgObj)
    // Optional: inform sender of new block balance and that message is queued
    c.Send(plib.SVR_MSGTO, utils.MarshalResponse(&models.MsgToResponseModel{
        Success: true,
        MsgID:   msgObj.ID,
        To:      msgObj.To,
        Blocks:  totalBlocks,
    }))
}

func (s *Server) cmdUser(c *Client, p []byte) {
	var userObj models.UserRequestModel
	if err := json.Unmarshal(p, &userObj); err != nil {
		log.Debug("unable to unmarshal packet")
		return
	}
	// Get our user from our users bucket
	user, err := dbUserGet(userObj.Username)
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
		User: models.User{
			Username: user.Username,
			PubKey:   user.PubKey,
		},
	}))
}

func (s *Server) cmdRegister(c *Client, p []byte) {
	var registerObj models.RegisterRequestModel
	if err := json.Unmarshal(p, &registerObj); err != nil {
		log.Debug("unable to unmarshal packet")
        c.Send(plib.SVR_REGISTER, utils.MarshalResponse(&models.RegisterResponseModel{
            Success: false,
            Message: "invalid register payload",
        }))
		return
	}
	// Some validation
	registerObj.Username = strings.ToLower(registerObj.Username)
	// Create our object to add to our DB
	user := &models.User{
		Username:       registerObj.Username,
		PassHash:       registerObj.PassHash,
		EncPrivKey:     registerObj.EncPrivKey,
		PubKey:         registerObj.PubKey,
		RegisteredDate: time.Now(),
		Blocks:         20,
		Access:         0,
	}
	// Register our user
    if err := dbUserAdd(user); err != nil {
        log.Debug(err)
        c.Send(plib.SVR_REGISTER, utils.MarshalResponse(&models.RegisterResponseModel{
            Success: false,
            Message: err.Error(),
        }))
        return
    }
    c.Send(plib.SVR_REGISTER, utils.MarshalResponse(&models.RegisterResponseModel{
        Success: true,
        Message: "registered",
    }))
}

func (s *Server) cmdLogin(c *Client, p []byte) {
	var loginObj models.LoginRequestModel
	if err := json.Unmarshal(p, &loginObj); err != nil {
		log.Debug("unable to unmarshal packet")
		return
	}
	// Some validation
	loginObj.Username = strings.ToLower(loginObj.Username)
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
	user, err := dbUserGet(loginObj.Username)
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
    if user, ok := s.getClient(c.Username); ok {
		// TODO: Should be a dedicated packet logging out the user.
		user.Conn.Close()
	}
	// Add client to our online list
	if err := s.clientAdd(c.Username, c); err != nil {
		log.Error(err)
	}
	c.Send(plib.SVR_LOGIN, utils.MarshalResponse(&models.LoginResponseModel{
		Success:    true,
		Message:    "success",
		Blocks:     user.Blocks,
		MsgCost:    blockCalcCost(),
		EncPrivKey: user.EncPrivKey,
	}))
}

func (s *Server) broadcast(cmdType int, buf []byte) {
    // Take a snapshot to avoid holding locks during network IO
    s.mu.RLock()
    clients := make([]*Client, 0, len(s.clients))
    for _, c := range s.clients { clients = append(clients, c) }
    s.mu.RUnlock()
    for _, c := range clients {
        go c.Send(cmdType, buf)
    }
}

func (s *Server) clientAdd(username string, c *Client) error {
	log.Debug("Adding client " + username)
    s.mu.Lock()
    defer s.mu.Unlock()
    _, exists := s.clients[username]
	if exists {
		return errors.New("user already exists")
	}
	s.clients[username] = c
	return nil
}

func (s *Server) clientRemoveByUsername(username string) {
	log.Debug("Removing client " + username)
    s.mu.Lock()
    defer s.mu.Unlock()
    if _, exists := s.clients[username]; exists {
        delete(s.clients, username)
    }
}

func (s *Server) clientRemoveByConnection(conn net.Conn) {
    s.mu.Lock()
    defer s.mu.Unlock()
    for _, c := range s.clients {
        if c.Conn == conn {
            log.Info("Removing client " + c.Username)
            delete(s.clients, c.Username)
        }
    }
}

func (s *Server) shutdown() {
	log.Debug("Server shutdown")
    s.mu.Lock()
    if s.closing {
        s.mu.Unlock()
        return
    }
    s.closing = true
    close(s.shutdownCh)
    s.mu.Unlock()
	if err := s.listener.Close(); err != nil {
		panic(err)
	}
    // Snapshot clients and close their connections
    s.mu.RLock()
    clients := make([]*Client, 0, len(s.clients))
    for _, c := range s.clients { clients = append(clients, c) }
    s.mu.RUnlock()
    for _, c := range clients {
        _ = c.Conn.Close()
    }
    // Wait for client goroutines to finish
    s.wg.Wait()
}

func (s *Server) authCheck(c *Client) bool {
    // Make sure we have a session set otherwise we kill their connection
    if c.Username == "" {
        log.Warn("Unauthorised. Connection closed. " + c.Conn.RemoteAddr().String())
        c.Conn.Close()
        return false
    }
    return true
}

// Helper methods for safe client access
func (s *Server) getClient(username string) (*Client, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    c, ok := s.clients[username]
    return c, ok
}

func (s *Server) listClients() []*Client {
    s.mu.RLock()
    defer s.mu.RUnlock()
    clients := make([]*Client, 0, len(s.clients))
    for _, c := range s.clients {
        clients = append(clients, c)
    }
    return clients
}

func (s *Server) isClosing() bool {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.closing
}
