package main

import (
	"crypto/sha512"
	"github.com/ProtonMail/gopenpgp/crypto"
)

func hashPassword(password string) []byte {
	hasher := sha512.New()
	hasher.Write([]byte(password))
	return hasher.Sum(nil)
}

func generateKeys() {
	hash := hashPassword("testing")
	var newHash []byte
	// RSA
	rsaKey, err := crypto.GenerateKey(localPart, domain, passphrase, "rsa", rsaBits)
}


func decryptRSA() {

}

func encryptRSA() {

}