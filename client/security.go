package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	gopenpgp "github.com/ProtonMail/gopenpgp/crypto"
	"io"
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
	fmt.Println("Now for hash key...")
	hashKey := hashString[:32]
	fmt.Println(hashKey)
	fmt.Println("Now for hash remainder...")
	// we want only 16 values
	hashRemainder := hashString[32:48]
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
	encryptedKey, err := encryptRSA([]byte(rsaKey), []byte(hashRemainder), []byte(hashKey))
	if err != nil {
		panic(err)
	}
	fmt.Println(encryptedKey)
	fmt.Println("Now for decrypt private key...")
	decryptedKey, err := decryptRSA(encryptedKey, []byte(hashRemainder), []byte(hashKey))
	if err != nil {
		panic(err)
	}
	fmt.Println(decryptedKey)
}

func encryptRSA(message []byte, addData []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := aesgcm.Seal(nonce, nonce, message, addData)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func decryptRSA(message string, addData []byte, key []byte) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(message)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("invalid nonce size")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, addData)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
