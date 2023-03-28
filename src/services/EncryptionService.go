package services

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

func EncryptPayload(payload string, sharedSecret string) (string, error) {
	return encryptAES(payload, sharedSecret)
}

func DecryptPayload(payload string, sharedSecret string) (string, error) {
	return decryptAES(payload, sharedSecret)
}

func SignRSASHA(message string, filename string) (string, error) {
	key, err := getPrivateKeyFromFile(filename)
	if err != nil {
		log.Error("Unable to retrieve key from %s", filename)
		return "", err
	}

	messageByte := []byte(message)

	msgHash := sha256.New()
	_, err = msgHash.Write(messageByte)
	if err != nil {
		log.Error("Unable to hash signature message")
		return "", err
	}

	msgHashSum := msgHash.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, msgHashSum)
	if err != nil {
		log.Error("Unable to sign message")
		return "", err
	}

	return hex.EncodeToString(signature), nil
}

func encryptAES(text string, keyStr string) (string, error) {
	key := []byte(keyStr)
	plaintext := []byte(text)
	paddedText := pkcs7Padding(plaintext, aes.BlockSize)

	ciphertext := make([]byte, len(paddedText))
	iv := ciphertext[:aes.BlockSize]

	aesBlock, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Error("Unable to encrypt AES")
		return "", err
	}

	aesCBC := cipher.NewCBCEncrypter(aesBlock, iv)
	aesCBC.CryptBlocks(ciphertext, paddedText)

	encrypted := hex.EncodeToString(ciphertext)

	return encrypted, nil
}

func decryptAES(text string, keyStr string) (string, error) {
	key := []byte(keyStr)
	ciphertext, err := hex.DecodeString(text)
	if err != nil {
		log.Error("Unable to decode ciphertext")
		return "", err
	}

	plaintext := make([]byte, len(ciphertext))
	iv := plaintext[:aes.BlockSize]

	aesBlock, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Error("Unable to decrypt AES")
		return "", err
	}

	aesCBC := cipher.NewCBCDecrypter(aesBlock, iv)
	aesCBC.CryptBlocks(plaintext, ciphertext)

	unpaddedText := pkcs7UnPadding(plaintext)

	decrypted := string(unpaddedText)

	return decrypted, nil
}

func pkcs7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs7UnPadding(plainText []byte) []byte {
	length := len(plainText)
	unpadding := int(plainText[length-1])
	return plainText[:(length - unpadding)]
}

func getPrivateKeyFromFile(filename string) (*rsa.PrivateKey, error) {
	pembytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error("Unable to read key from %s", filename)
		return nil, err
	}

	data, _ := pem.Decode([]byte(pembytes))
	key, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		log.Error("Unable to parse key from %s", filename)
		return nil, err
	}

	return key, nil
}
