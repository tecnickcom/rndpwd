package main

import (
	"crypto/rand"
	"math/big"
)

// getNewPassword returns a random password containing "length" characters from the "charset" string
// NOTE: the charsetLength must be correctly set, as there are no controls in this function to improve performances
func getNewPassword(length int, charset string) string {
	// count generated passwords
	stats.Increment("passwords")

	// time this function
	defer stats.NewTiming().Send("getNewPassword.time")

	password := make([]byte, length)

	chars := []byte(charset)
	maxValue := new(big.Int).SetInt64(int64(len(charset)))

	for i := 0; i < length; i++ {
		rnd, _ := rand.Int(rand.Reader, maxValue)
		password[i] = chars[rnd.Int64()]
	}

	return string(password)
}

// getAllPassword returns the specified amount of random passwords
func getAllPassword(quantity int, length int, charset string) []string {
	defer stats.NewTiming().Send("getAllPassword.time")
	passwords := make([]string, quantity)
	for i := 0; i < quantity; i++ {
		passwords[i] = getNewPassword(length, charset)
	}
	return passwords
}
