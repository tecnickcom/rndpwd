package main

import (
	"crypto/rand"
	"log"
	"math/big"
)

// GetNewPassword returns a random password containing "length" characters from the "charset" string
func getNewPassword(length int, charset string, charsetLength int) string {

	password := make([]byte, length)
	chars := []byte(charset)
	maxValue := new(big.Int).SetInt64(int64(charsetLength))

	for i := 0; i < length; i++ {
		rnd, err := rand.Int(rand.Reader, maxValue)
		if err != nil {
			log.Fatal(err)
		}
		password[i] = chars[rnd.Int64()]
	}

	return string(password)
}
