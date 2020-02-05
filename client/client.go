package main

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"errors"
	gopenpgp "github.com/ProtonMail/gopenpgp/crypto"
	"github.com/labstack/gommon/log"
	"github.com/syleron/426c/common/models"
	plib "github.com/syleron/426c/common/packet"
	"github.com/syleron/426c/common/utils"
	"net"
	"strings"
)

type Client struct {
	Reader *bufio.Reader
	Writer *bufio.Writer
	Conn   net.Conn
	MQ     *MessageQueue
}

func setupClient() (*Client, error) {
	// Setup our listener
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		app.Stop()
	}
	config := tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	config.Rand = rand.Reader
	// connect to this socket
	// TODO This should be a client command rather done automagically.
	conn, err := tls.Dial("tcp", "127.0.0.1:9000", &config)
	if err != nil {
		return &Client{}, errors.New("unable to connect to host")
	}
	// Put our handlers into a go routine
	c := &Client{
		Writer: bufio.NewWriter(conn),
		Reader: bufio.NewReader(conn),
		Conn:   conn,
		MQ:     NewMessageQueue(),
	}
	// Put our handlers into a go routine
	go c.connectionHandler()
	return c, nil
}

func (c *Client) Send(cmdType int, buf []byte) (int, error) {
	return c.Conn.Write(plib.PacketForm(byte(cmdType), buf))
}

func (c *Client) connectionHandler() {
	for {
		p, err := plib.PacketRead(c.Reader)
		if err != nil {
			app.Stop()
		}
		c.commandRouter(p)
	}
}

func (c *Client) commandRouter(p []byte) {
	switch p[0] {
	case plib.SVR_LOGIN:
		c.svrLogin(p[1:])
	case plib.SVR_USER:
		c.svrUser(p[1:])
	case plib.SVR_MSGTO:
		c.svrMsgTo(p[1:])
	case plib.SVR_MSG:
		c.svrMsg(p[1:])
	default:
	}
}

// ||
// Client Requests
// ||

func (c *Client) cmdRegister(username string, password string) {
	var pgp = gopenpgp.GetGopenPGP()
	// Generate password hash
	hashString := hashPassword(password)
	// Calculate hash key
	hashKey := hashString[:32]
	// Calculate hash remainder
	hashRemainder := hashString[32:48]
	// Generate RSA key
	rsaKey, err := pgp.GenerateKey(
		username,
		"secure.426c.net",
		hashString,
		"rsa",
		4096,
	)
	// save our key
	if err := utils.WriteFile(rsaKey, username); err != nil{
		app.Stop()
	}
	if err != nil {
		app.Stop()
	}
	keyRing, err := gopenpgp.ReadArmoredKeyRing(strings.NewReader(rsaKey))
	if err != nil {
		app.Stop()
	}
	publicKey, err := keyRing.GetArmoredPublicKey()
	if err != nil {
		app.Stop()
	}
	// Encrypt our private RSA key
	encryptedKey, err := encryptRSA([]byte(rsaKey), []byte(hashRemainder), []byte(hashKey))
	if err != nil {
		app.Stop()
	}
	// Create our object to send
	registerObject := &models.RegisterRequestModel{
		Username:   username,
		PassHash:   hashRemainder,
		EncPrivKey: encryptedKey,
		PubKey:     publicKey,
	}
	// Send our username, hash remainder, encrypted private key, and readable public key.
	_, err = c.Send(
		plib.CMD_REGISTER,
		utils.MarshalResponse(registerObject),
	)
	if err != nil {
		app.Stop()
	}
}

func (c *Client) cmdLogin(username string, password string) {
	// Generate password hash
	hashString := hashPassword(password)
	// Calculate hash remainder
	hashRemainder := hashString[32:48]
	// Create our object to send
	registerObject := &models.LoginRequestModel{
		Username: username,
		Password: hashRemainder,
		Version:  VERSION,
	}
	// Set our local variables
	pHash = hashString
	// Send our username, hash remainder.
	_, err := c.Send(
		plib.CMD_LOGIN,
		utils.MarshalResponse(registerObject),
	)
	if err != nil {
		app.Stop()
	}
}

// cmdMsgTo - Send a private encrypted message to a particular user
func (c *Client) cmdMsgTo(m *models.Message) {
	// Attempt to send our message
	_, err := c.Send(plib.CMD_MSGTO, utils.MarshalResponse(&models.MsgToRequestModel{
		Message: *m,
	}))
	if err != nil {
		app.Stop()
	}
}

// ||
// Server Responses
// ||

func (c *Client) svrRegister(p []byte) error {
	var regObj models.RegisterResponseModel
	if err := json.Unmarshal(p, &regObj); err != nil {
		app.Stop()
	}
	if !regObj.Success {
		showError(ClientError{
			Message: regObj.Message,
			Button:  "Continue",
			Continue: func() {
				pages.SwitchToPage("login")
			},
		})
		return nil
	}
	return nil
}

func (c *Client) svrLogin(p []byte) {
	var loginObj models.LoginResponseModel
	if err := json.Unmarshal(p, &loginObj); err != nil {
		log.Debug("unable to unmarshal packet")
		return
	}
	// Make sure our response object was successful
	if !loginObj.Success {
		showError(ClientError{
			Message: loginObj.Message,
			Button:  "Continue",
			Continue: func() {
				pages.SwitchToPage("login")
			},
		})
		return
	}
	// Set our logged in user
	lUser = loginObj.Username
	// Load our private key
	b, err := utils.LoadFile(lUser)
	if err != nil {
		showError(ClientError{
			Message: "Login failed. Unable to load private key for " + loginObj.Username + ".",
			Button:  "Continue",
		})
		return
	}
	// Set our private key
	privKey = string(b)
	// Success, switch pages
	pages.SwitchToPage("inbox")
	// get our contacts
	drawContactsList()
}

func (c *Client) svrMsg(p []byte) {
	var msgObj models.MsgResponseModel
	if err := json.Unmarshal(p, &msgObj); err != nil {
		log.Debug("unable to unmarshal packet")
		return
	}
	// Mark our message as being received successfully
	msgObj.Success = true
	// Add our message to our local DB
	if _, err := dbMessageAdd(&msgObj.Message); err != nil {
		panic(err)
	}
	panic("balls")
	// reload our message container
	go loadMessages(msgObj.From, inboxMessageContainer)
}

func (c *Client) svrMsgTo(p []byte) {
	var msgObj models.MsgToResponseModel
	if err := json.Unmarshal(p, &msgObj); err != nil {
		log.Debug("unable to unmarshal packet")
		return
	}
	// Mark our request successful
	if msgObj.Success {
		if err := dbMessageSuccess(msgObj.MsgID, msgObj.To); err != nil {
			panic(err)
		}
		// redraw our messages
		go loadMessages(msgObj.To, inboxMessageContainer)
	}
}

// svrUser - User Object response from network and update our local DB
func (c *Client) svrUser(p []byte) {
	var userObj models.UserResponseModel
	if err := json.Unmarshal(p, &userObj); err != nil {
		log.Debug("unable to unmarshal packet")
		return
	}
	if !userObj.Success {
		showError(ClientError{
			Message:  userObj.Message,
			Button:   "Continue",
			Continue: nil,
		})
		return
	}
	// Insert our user into our local DB
	dbUserAdd(userObj.User)
	// Reset UI
	inboxToField.SetText("")
	app.SetFocus(userListContainer)
	app.Draw() // force draw to speed up the changes
}

func (c *Client) Close() {}
