package main

import (
	"crypto/aes"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	gopenpgp "github.com/ProtonMail/gopenpgp/crypto"
)

/**
Returns a hash in hex
 */
func hashPassword(password string) []byte {
	hasher := sha512.New()
	hasher.Write([]byte(password))
	return hasher.Sum(nil)
}

func generateKeys() {
	hash := hashPassword("testing")
	fmt.Println(hex.EncodeToString(hash)) // This is the remainder
	// Remove the first 32 bits
	initVector := make([]byte, 16)
	for i, b := range hash {
		if i >= 16 && len(initVector) <= 16 {
			initVector = append(initVector, b)
		}
	}
	fmt.Println(hex.EncodeToString(initVector))
	//var newHash []byte
	//// RSA
	var pgp = gopenpgp.GetGopenPGP()
	rsaKey, err := pgp.GenerateKey(
		"testing",
		"secure.426c.net",
		hex.EncodeToString(hash),
		"rsa",
		4096,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(rsaKey) // This is the key
	// encode our private key
	c, err := aes.NewCipher(hash)
	if err != nil {
		panic(err)
	}
	out := make([]byte, len(rsaKey))
	c.Encrypt(out, []byte(rsaKey))
	fmt.Println(hex.EncodeToString(out))
}


func decryptRSA() {

}

func encryptRSA() {

}