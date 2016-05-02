package main

import (
	"fmt"
	"testing"
)

func TestGetConfigParams(t *testing.T) {
	prm := getConfigParams()
	if prm.quantity != 10 {
		t.Error(fmt.Errorf("The expected quantity is 10, found %d", prm.quantity))
	}
	if prm.length != 32 {
		t.Error(fmt.Errorf("The expected length is 32, found %d", prm.quantity))
	}
	if prm.charset != "!#$%&()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ^_abcdefghijklmnopqrstuvwxyz~" {
		t.Error(fmt.Errorf("Fond differnt charset than expected"))
	}
}

func TestCheckParams(t *testing.T) {
	err := checkParams(&params{quantity: 1, length: 2, charset: "abc"})
	if err != nil {
		t.Error(fmt.Errorf("No errors are expected"))
	}
}

func TestCheckParamsErrorsServer(t *testing.T) {
	err := checkParams(&params{server: true, httpaddr: ""})
	if err == nil {
		t.Error(fmt.Errorf("An error was expected because the server address is empty"))
	}
}

func TestCheckParamsErrorsQuantity(t *testing.T) {
	err := checkParams(&params{quantity: 0, length: 2, charset: "abc"})
	if err == nil {
		t.Error(fmt.Errorf("An error was expected because the quantity is <= 0"))
	}
}

func TestCheckParamsErrorsLength(t *testing.T) {
	err := checkParams(&params{quantity: 1, length: 0, charset: "abc"})
	if err == nil {
		t.Error(fmt.Errorf("An error was expected because the length is <= 0"))
	}
}

func TestCheckParamsErrorsCharsetLength(t *testing.T) {
	err := checkParams(&params{quantity: 1, length: 2, charset: "a"})
	if err == nil {
		t.Error(fmt.Errorf("An error was expected because the charset length is < 2"))
	}

	err = checkParams(&params{quantity: 1, length: 2, charset: "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"})
	if err == nil {
		t.Error(fmt.Errorf("An error was expected because the charset length is > 92"))
	}
}

func TestCheckParamsErrorsValidCharset(t *testing.T) {
	err := checkParams(&params{quantity: 1, length: 2, charset: "ab cd"})
	if err == nil {
		t.Error(fmt.Errorf("An error was expected because the charset contains an invalid character"))
	}
}
