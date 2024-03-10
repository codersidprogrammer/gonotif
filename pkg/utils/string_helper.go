package utils

import (
	"crypto/rand"
	"math/big"
	"regexp"
	"strings"
)

func ReplaceWhitespaceAndRemoveSpecialChars(input string) string {
	// Replace whitespace with underscores
	input = strings.ReplaceAll(input, " ", "_")

	// Remove special characters using regular expressions
	reg := regexp.MustCompile("[^a-zA-Z0-9_]")
	input = reg.ReplaceAllString(input, "")

	return input
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result strings.Builder
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		result.WriteByte(charset[randomIndex.Int64()])
	}
	return result.String()
}
