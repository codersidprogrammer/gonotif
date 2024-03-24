package utils

import (
	"crypto/rand"
	"errors"
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

func CheckIfHasSpecifiedSuffix(s string, limiter string, suffix string) bool {
	if len(s) == 0 {
		return false
	}

	parts := strings.Split(s, limiter)
	last := parts[len(parts)-1]

	return strings.HasSuffix(last, suffix)
}

func GetItemFromSplitText(splitText string, delimiter string, index int) (string, error) {
	if len(splitText) == 0 {
		return "", errors.New("split text is not good form")
	}
	parts := strings.Split(splitText, delimiter)
	return parts[index], nil
}
