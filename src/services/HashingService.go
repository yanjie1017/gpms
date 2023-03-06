package services

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"math/rand"
	"strings"

	model "src/models"

	"golang.org/x/crypto/scrypt"
)

func GeneratePassword(request model.PasswordGenerationRequest, passwordInfo model.PasswordGenerationInfo, systemKey string) string {
	var siteInfo string = passwordInfo.Metadata
	var length int = int(passwordInfo.Length)

	var generationToken string = request.GenerationToken
	var userInput = request.UserInput

	var salt []byte = generateSalt(siteInfo, generationToken, systemKey)
	passwordBytes, _ := scrypt.Key([]byte(userInput), salt, 16384, 8, 1, 32)
	var passwordString string = hex.EncodeToString(passwordBytes)
	fmt.Println(passwordString)

	var mappedPassword string = mapPassword(passwordString, length)

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

	return salt
}

func toHashCode(s string) uint64 {
	h := fnv.New64()
	h.Write([]byte(s))
	return h.Sum64()
}

func mapPassword(password string, requiredLength int) string {
	var seed int64 = int64(toHashCode(password))

	passwordBytes := []byte(password)

	// Define the list of possible characters to include in the output string
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	symbol := "~`! @#$%^&*()_-+={[}]|:;<,>.?/"

	// Seed the random number generator with a fixed value to ensure consistent output
	rand.Seed(seed)

	originalLength := requiredLength / 2
	upperLength := (requiredLength - originalLength) / 2
	// symbolLength := requiredLength - originalLength - upperLength

	output := make([]byte, requiredLength)

	// Add characters from the hash to the output string
	for i := 0; i < originalLength; i++ {
		output[i] = passwordBytes[i]
	}

	for i := originalLength; i < originalLength+upperLength; i++ {
		output[i] = upper[rand.Intn(len(upper))]
	}

	for i := originalLength + upperLength; i < requiredLength; i++ {
		output[i] = symbol[rand.Intn(len(symbol))]
	}

	// Shuffle the output string to ensure randomness
	rand.Shuffle(len(output), func(i, j int) {
		output[i], output[j] = output[j], output[i]
	})

	return string(output)
}
