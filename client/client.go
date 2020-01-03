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
	"os"
	"strings"
)

type Client struct {
	Reader *bufio.Reader
	Writer *bufio.Writer
	Conn   net.Conn
}

func setupClient() (*Client, error) {
	// Setup our listener
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		panic(err)
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
			os.Exit(1)
		}
		c.commandRouter(p)
	}
}

func (c *Client) commandRouter(p []byte) {
	switch p[0] {
	case plib.SVR_LOGIN:
		c.svrLogin(p[1:])
	default:
	}
}

func (c *Client) msgRegister(username string, password string) {
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
	if err != nil {
		panic(err)
	}
	keyRing, err := gopenpgp.ReadArmoredKeyRing(strings.NewReader(rsaKey))
	if err != nil {
		panic(err)
	}
	publicKey, err := keyRing.GetArmoredPublicKey()
	if err != nil {
		panic(err)
	}
	// Encrypt our private RSA key
	encryptedKey, err := encryptRSA([]byte(rsaKey), []byte(hashRemainder), []byte(hashKey))
	if err != nil {
		panic(err)
	}
	// Create our object to send
	registerObject := &models.RegisterRequestModel{
		Username:   username,
		PassHash:    hashRemainder,
		EncPrivKey: encryptedKey,
		PubKey:     publicKey,
	}
	// Send our username, hash remainder, encrypted private key, and readable public key.
	_, err = c.Send(
		plib.CMD_REGISTER,
		utils.MarshalResponse(registerObject),
	)
	if err != nil {
		panic(err)
	}
}

func (c *Client) svrRegister() error {
	return nil
}

func (c *Client) msgLogin(username string, password string) {
	// Generate password hash
	hashString := hashPassword(password)
	// Calculate hash remainder
	hashRemainder := hashString[32:48]
	// Create our object to send
	registerObject := &models.LoginRequestModel{
		Username:   username,
		Password: hashRemainder,
		Version: VERSION,
	}
	// Send our username, hash remainder.
	_, err := c.Send(
		plib.CMD_LOGIN,
		utils.MarshalResponse(registerObject),
	)
	if err != nil {
		panic(err)
	}
}

func (c *Client) svrLogin(p []byte) {
	var loginObj models.LoginResponseModel
	if err := json.Unmarshal(p, &loginObj); err != nil {
		log.Debug("unable to unmarshal packet")
		return
	}
	if !loginObj.Success {
		showError(ClientError{
			Message:  loginObj.Message,
			Button:   "Continue",
			Continue: func() {
				pages.SwitchToPage("login")
			},
		})
		return
	}
	pages.SwitchToPage("inbox")
}

func (c *Client) msgSearch(username string) error {

	return nil
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
