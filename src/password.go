package main

import (
	"crypto/rand"
	"math/big"
)

// getNewPassword returns a random password containing "length" characters from the "charset" string
// NOTE: the charsetLength must be correctly set, as there are no controls in this function to improve performances
func getNewPassword(appParams *params) string {
	// count generated passwords
	stats.Increment("passwords")

	// time this function
	defer stats.NewTiming().Send("getNewPassword.time")

	password := make([]byte, appParams.length)
	// #nosec
	chars := []byte(appParams.charset)
	maxValue := new(big.Int).SetInt64(int64(appParams.charsetLength))

	for i := 0; i < appParams.length; i++ {
		rnd, _ := rand.Int(rand.Reader, maxValue)
		password[i] = chars[rnd.Int64()]
	}

	return string(password)
}

// getAllPassword returns the specified amount of random passwords
func getAllPassword(appParams *params) []string {
	defer stats.NewTiming().Send("getAllPassword.time")
	passwords := make([]string, appParams.quantity)
	for i := 0; i < appParams.quantity; i++ {
		passwords[i] = getNewPassword(appParams)
	}
	return passwords
}
