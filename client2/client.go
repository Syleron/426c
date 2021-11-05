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
	"github.com/syleron/426c/common/security"
	"github.com/syleron/426c/common/utils"
	"net"
	"os"
	"strings"
	"sync"
)

type Client struct {
	// Connection reader
	Reader *bufio.Reader
	// Connection writer
	Writer *bufio.Writer
	// Connection object
	Conn   net.Conn
	// Request cache
	Cache *Cache
}

type Cache struct {
	// TODO: Consider Cache to only accept string values
	store map[string]interface{}
	sync.Mutex
}

func (c *Cache) Add(key string, value interface{}) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.store[key]; ok {
		return errors.New("cache add failed. key already exists")
	}
	c.store[key] = value
	return nil
}

func (c *Cache) Update(key string, value interface{}) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.store[key]; ok {
		c.store[key] = value
		return nil
	}
	return errors.New("cache update failed. key " + key + " doesn't exist ")
}

func (c *Cache) Remove(key string) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.store[key]; !ok {
		return errors.New("cache remove failed. key " + key + " doesn't exist")
	}
	delete(c.store, key)
	return nil
}

func (c *Cache) Get(key string) (interface{}, error) {
	if _, ok := c.store[key]; !ok {
		return nil, errors.New("cache get failed. key " + key + " doesn't exist")
	}
	return c.store[key], nil
}

func setupClient(address string) (*Client, error) {
	// Setup our listener
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatal("unable to load cert keys")
		os.Exit(0)
	}
	config := tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	config.Rand = rand.Reader
	// connect to this socket
	// TODO This should be a client command rather done automagically.
	conn, err := tls.Dial("tcp", address, &config)
	if err != nil {
		return &Client{}, errors.New("unable to connect to host")
	}
	// Put our handlers into a go routine
	c := &Client{
		Writer: bufio.NewWriter(conn),
		Reader: bufio.NewReader(conn),
		Conn:   conn,
		Cache: &Cache{
			store: make(map[string]interface{}),
			Mutex: sync.Mutex{},
		},
	}
	// Set our default cache options
	c.Cache.Add("blocks", 0)
	c.Cache.Add("msgCost", 0)
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
			os.Exit(0)
		}
		c.commandRouter(p)
	}
}

func (c *Client) commandRouter(p []byte) {
	if len(p) <= 0 {
		return
	}
	switch p[0] {
	case plib.SVR_LOGIN:
		c.svrLogin(p[1:])
	case plib.SVR_USER:
		c.svrUser(p[1:])
	case plib.SVR_MSGTO:
		c.svrMsgTo(p[1:])
	case plib.SVR_MSG:
		c.svrMsg(p[1:])
	case plib.SVR_BLOCK:
		c.svrBlock(p[1:])
	default:
		panic("balls")
	}
}

// ||
// Client Requests
// ||

func (c *Client) cmdRegister(username string, password string) {
	var pgp = gopenpgp.GetGopenPGP()
	// Some validation
	username = strings.ToLower(username)
	// Generate password hash
	hashString := security.SHA512HashEncode(password)
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
	if err != nil {
		os.Exit(0)
	}
	keyRing, err := gopenpgp.ReadArmoredKeyRing(strings.NewReader(rsaKey))
	if err != nil {
		os.Exit(0)
	}
	publicKey, err := keyRing.GetArmoredPublicKey()
	if err != nil {
		os.Exit(0)
	}
	// Encrypt our private RSA key
	encryptedKey, err := security.EncryptRSA([]byte(rsaKey), []byte(hashRemainder), []byte(hashKey))
	if err != nil {
		os.Exit(0)
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
		os.Exit(0)
	}
}

func (c *Client) cmdLogin(username string, password string) {
	// Some validation
	username = strings.ToLower(username)
	// Add our username to our connection cache
	c.Cache.Add("username", username)
	// Generate password hash
	hashString := security.SHA512HashEncode(password)
	// Set our client password hash
	c.Cache.Add("passHash", hashString)
	// Calculate hash remainder
	hashRemainder := hashString[32:48]
	// Create our object to send
	registerObject := &models.LoginRequestModel{
		Username: username,
		Password: hashRemainder,
		Version:  VERSION,
	}
	// Send our username, hash remainder.
	_, err := c.Send(
		plib.CMD_LOGIN,
		utils.MarshalResponse(registerObject),
	)
	if err != nil {
		os.Exit(0)
	}
}

// cmdMsgTo - Send a private encrypted message to a particular user
func (c *Client) cmdMsgTo(m *models.Message) {
	// Retry failed messages
	//if inboxFailedMessageCount > 0 {
	//	inboxRetryFailedMessages(m.To)
	//}
	// Attempt to send our message
	_, err := c.Send(plib.CMD_MSGTO, utils.MarshalResponse(&models.MsgToRequestModel{
		Message: *m,
	}))
	if err != nil {
		os.Exit(0)
	}
}

func (c *Client) cmdUser(username string) error {
	if username == "" {
		return errors.New("please specify a username")
	}
	// Attempt to send our message
	//_, err := client.Send(plib.CMD_USER, utils.MarshalResponse(&models.UserRequestModel{
	//	Username: username,
	//}))
	//if err != nil {
	//	os.Exit(0)
	//}
	return nil
}

// ||
// Server Responses
// ||
func (c *Client) svrBlock(p []byte) {
	var blockObj models.BlockResponseModel
	if err := json.Unmarshal(p, &blockObj); err != nil {
		os.Exit(0)
	}
	// Set our available blocks
	if err := c.Cache.Update("blocks", blockObj.Blocks); err != nil {
		log.Fatal(err)
	}
	// Set our message cost
	if err := c.Cache.Update("msgCost", blockObj.MsgCost); err != nil {
		log.Fatal(err)
	}
}

func (c *Client) svrRegister(p []byte) {
	var regObj models.RegisterResponseModel
	if err := json.Unmarshal(p, &regObj); err != nil {
		os.Exit(0)
	}
	if !regObj.Success {
		showError(ClientError{
			Message: regObj.Message,
			Button:  "Continue",
			Continue: func() {
				//pages.SwitchToPage("login")
			},
		})
		return
	}
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
				//pages.SwitchToPage("login")
			},
		})
		return
	}
	// Make sure we have our encrypted private key
	if loginObj.EncPrivKey == "" {
		showError(ClientError{
			Message: "Missing private key. Internal error",
			Button:  "Continue",
			Continue: func() {
				//pages.SwitchToPage("login")
			},
		})
		return
	}
	// Get our pass hash
	passHash, err := c.Cache.Get("passHash")
	if err != nil {
		log.Fatal(err)
	}
	// Calculate hash key
	hashKey := passHash.(string)[:32]
	// Calculate hash remainder
	hashRemainder := passHash.(string)[32:48]
	// Decrypt private key
	pKey, err := security.DecryptRSA(loginObj.EncPrivKey, []byte(hashRemainder), []byte(hashKey))
	if err := c.Cache.Add("pKey", pKey); err != nil {
		log.Fatal(err)
	}
	// Set our available blocks
	if err := c.Cache.Update("blocks", loginObj.Blocks); err != nil {
		log.Fatal(err)
	}
	// Set our message cost
	if err := c.Cache.Update("msgCost", loginObj.MsgCost); err != nil {
		log.Fatal(err)
	}
	// Success, switch pages
	//pages.SwitchToPage("inbox")
	// Populate our user list
	//userList.PopulateFromDB()
}

func (c *Client) svrMsg(p []byte) {
	var msgObj models.MsgResponseModel
	if err := json.Unmarshal(p, &msgObj); err != nil {
		log.Debug("unable to unmarshal packet")
		return
	}
	// Check if we have the user as a contact
	_, err := dbUserGet(msgObj.From)
	if err != nil {
		c.cmdUser(msgObj.From)
	}
	// Mark our message as being received successfully
	msgObj.Success = true
	// Set our user online if they are current marked offline
	//userList.SetUserOnline(msgObj.From)
	// Add our message to our local DB
	if _, err := dbMessageAdd(&msgObj.Message, msgObj.From); err != nil {
		panic(err)
	}
	// Make sure we are viewing the messages for the incoming message
	//if inboxSelectedUsername == msgObj.From {
	//	// reload our message container
	//	go messageLoad(msgObj.From, inboxMessageContainer)
	//}
	// Increment our userlist indicator
	//userList.NewMessage(msgObj.From)
}

func (c *Client) svrMsgTo(p []byte) {
	var msgObj models.MsgToResponseModel
	if err := json.Unmarshal(p, &msgObj); err != nil {
		log.Debug("unable to unmarshal packet")
		return
	}
	// Mark our request successful
	if msgObj.Success {
		// Update our message in our DB
		dbMessageSuccess(msgObj.MsgID, msgObj.To)
		// Update our user list
		//userList.SetUserOnline(msgObj.To)
		// Set our available blocks
		if err := c.Cache.Update("blocks", msgObj.Blocks); err != nil {
			log.Fatal(err)
		}
	} else {
		// Update our message in our DB
		dbMessageFail(msgObj.MsgID, msgObj.To)
		// Update our user list
		if msgObj.Blocks > 0 {
			//userList.SetUserOffline(msgObj.To)
		}
		// Add to the number of failed messages we have.
		//inboxFailedMessageCount++ // TODO: Consider moving this to userlist
	}
	// redraw our messages
	//go messageLoad(msgObj.To, inboxMessageContainer)
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
	// add user to our user list
	//userList.AddUser(userObj.User.Username)
	// Reset UI
	//inboxToField.SetText("")
	//app.SetFocus(userListContainer)
}
