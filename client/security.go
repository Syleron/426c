package main

import (
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
	fmt.Println("Now for password Hash...")
	hash := hashPassword("testing")
	hashString := hex.EncodeToString(hash)
	fmt.Println(hashString)
	fmt.Println("Now for hash remainder...")
	hashKey := hashString[32:]
	fmt.Println(hashKey)
	// we want only 16 values
	hashRemainder := hashKey[:16]
	fmt.Println(hashRemainder)
	// TODO: Send the remainder to the server
	//var newHash []byte
	//// RSA
	fmt.Println("Now for RSA Key...")
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
	fmt.Println("Now for encrypted private key...")
	// encode our private key
	//c, err := aes.NewCipher([]byte(initVector))
	//if err != nil {
	//	panic(err)
	//}
	//out := make([]byte, len(rsaKey))
	//c.Encrypt(out, []byte(rsaKey))
	//fmt.Println(hex.EncodeToString(out))
	//fmt.Println("Now for descrypted private key...")
}


func decryptRSA() {

}

func encryptRSA() {

}