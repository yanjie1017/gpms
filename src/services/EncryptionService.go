package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"strings"
)

/*
* 1. Hash sharedEncryptionKey with encryptionToken as salt to generate encryptionKey
* 2. Encrypt payload with encryptionKey
 */
func EncryptPayload(payload string, sharedEncryptionKey string, encryptionToken string) string {
	// https://gist.github.com/kkirsche/e28da6754c39d5e7ea10
	var encryptionKey []byte = getEncryptionKey(sharedEncryptionKey, encryptionToken)
	aesBlock, _ := aes.NewCipher(encryptionKey)
	gcmInstance, _ := cipher.NewGCM(aesBlock)

	nonce := make([]byte, gcmInstance.NonceSize())
	_, _ = io.ReadFull(rand.Reader, nonce)

	var ciphertext []byte = gcmInstance.Seal(nil, nonce, []byte(payload), nil)
	return string(ciphertext)
}

func getEncryptionKey(sharedEncryptionKey string, encryptionToken string) []byte {
	var concatenatedString strings.Builder
	concatenatedString.WriteString(sharedEncryptionKey)
	concatenatedString.WriteString(encryptionToken)

	sha_256 := sha256.New()
	sha_256.Write([]byte(concatenatedString.String()))
	return sha_256.Sum(nil)
}

/*
* 1.
 */
func DecryptPayload() {
	// TODO
	// API key
	// shared encryption key
}
