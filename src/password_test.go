package main

import (
	"fmt"
	"testing"
)

var charset = "!#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnopqrstuvwxyz{|}~"
var charsetLength = len(charset)

func TestGetNewPassword(t *testing.T) {
	var password string
	length := 0
	for i := 8; i < 64; i += 8 {
		password = getNewPassword(&params{length: i, charset: charset, charsetLength: charsetLength})
		// #nosec
		if length = len(password); length != i {
			t.Error(fmt.Errorf("The expected password length is %d, found %d", i, length))
		}
	}
}

func TestGetAllPassword(t *testing.T) {
	quantity := 13
	passwords := getAllPassword(&params{quantity: quantity, length: 17, charset: charset, charsetLength: charsetLength})
	numPasswords := len(passwords)
	// #nosec
	if numPasswords != quantity {
		t.Error(fmt.Errorf("The expected number of password is %d, found %d", quantity, numPasswords))
	}
}

func BenchmarkGetNewPassword(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getNewPassword(&params{length: 32, charset: charset, charsetLength: charsetLength})
	}
}
