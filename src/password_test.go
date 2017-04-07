package main

import (
	"fmt"
	"testing"
)

var charset = "!#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnopqrstuvwxyz{|}~"

func TestGetNewPassword(t *testing.T) {
	var password string
	for i := 8; i < 64; i += 8 {
		password = getNewPassword(i, charset)

		if length := len(password); length != i {
			t.Error(fmt.Errorf("The expected password length is %d, found %d", i, length))
		}
	}
}

func TestGetAllPassword(t *testing.T) {
	quantity := 13
	passwords := getAllPassword(quantity, 17, charset)
	numPasswords := len(passwords)

	if numPasswords != quantity {
		t.Error(fmt.Errorf("The expected number of password is %d, found %d", quantity, numPasswords))
	}
}

func BenchmarkGetNewPassword(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getNewPassword(32, charset)
	}
}
