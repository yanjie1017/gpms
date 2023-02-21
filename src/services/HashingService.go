package services

import (
	"crypto/sha512"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)

func GeneratePassword(siteInfo string, generationToken string, systemKey string, userInput string, length int) string {

	var salt []byte = generateSalt(siteInfo, generationToken, systemKey)
	password, _ := scrypt.Key([]byte(userInput), salt, 16384, 8, 1, 32)

	var mappedPassword string = mapPassword(password, length)
	fmt.Printf("scrypt:\t%x\n", mappedPassword)
	return mappedPassword
}

func generateSalt(siteInfo string, generationToken string, systemKey string) []byte {
	var concatenatedString strings.Builder
	concatenatedString.WriteString(siteInfo)
	concatenatedString.WriteString(siteInfo)
	concatenatedString.WriteString(siteInfo)

	sha_512 := sha512.New()
	sha_512.Write([]byte(concatenatedString.String()))
	var salt []byte = sha_512.Sum(nil)

	fmt.Printf("sha512:\t%x\n", salt)
	return salt
}

func mapPassword(password []byte, length int) string {
	var mappedPassword string = string(password)
	// TODO
	return mappedPassword
}
