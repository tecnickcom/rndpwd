package main

import (
	"fmt"
	"testing"
)

var charset = "!#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnopqrstuvwxyz{|}~"
var charsetLength = len(charset)

func TestGetNewPassword(t *testing.T) {
	password := ""
	length := 0
	for i := 8; i < 64; i += 8 {
		password = getNewPassword(i, "ab", 2)
		if length = len(password); length != i {
			t.Error(fmt.Errorf("The expected password length is %d, found %d", i, length))
		}
	}
}

func BenchmarkGetNewPassword(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getNewPassword(32, charset, charsetLength)
	}
}
