package services

import (
	"crypto/sha512"
	"encoding/hex"
	"hash/fnv"
	"math/rand"
	"strings"

	model "src/models"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/scrypt"
)

func GeneratePassword(request model.PasswordGenerationRequest, passwordInfo model.PasswordGenerationInfo, systemKey string) (string, error) {
	var siteInfo string = passwordInfo.Metadata
	var length = int(passwordInfo.Length)

	var generationToken string = request.GenerationToken
	var userInput = request.UserInput

	salt, err := generateSalt(siteInfo, generationToken, systemKey)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to generate salt")
		return "", err
	}

	passwordBytes, err := scrypt.Key([]byte(userInput), salt, 16384, 8, 1, 32)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to hash user input")
		return "", err
	}

	var passwordString string = hex.EncodeToString(passwordBytes)

	mappedPassword, err := mapPassword(passwordString, length)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to map password")
		return "", err
	}

	return mappedPassword, nil
}

func generateSalt(siteInfo string, generationToken string, systemKey string) ([]byte, error) {
	var concatenatedString strings.Builder
	concatenatedString.WriteString(siteInfo)
	concatenatedString.WriteString(generationToken)
	concatenatedString.WriteString(systemKey)

	sha_512 := sha512.New()
	_, err := sha_512.Write([]byte(concatenatedString.String()))

	if err != nil {
		return nil, err
	}

	var salt []byte = sha_512.Sum(nil)

	return salt, nil
}

func toHashCode(s string) (uint64, error) {
	var i uint64
	h := fnv.New64()
	_, err := h.Write([]byte(s))
	if err != nil {
		return i, err
	}
	i = h.Sum64()
	return i, nil
}

func mapPassword(password string, requiredLength int) (string, error) {
	seed, err := toHashCode(password)
	if err != nil {
		return "", err
	}

	passwordBytes := []byte(password)

	// Define the list of possible characters to include in the output string
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	symbol := "~`!@#$%^&*()_-+={[}]|:;<,>.?/"

	// Seed the random number generator with a fixed value to ensure consistent output
	rand.Seed(int64(seed))

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

	return string(output), nil
}
