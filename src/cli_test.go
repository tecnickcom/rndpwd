package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCli(t *testing.T) {
	cmd := cli()
	if cmdtype := reflect.TypeOf(cmd).String(); cmdtype != "*cobra.Command" {
		t.Error(fmt.Errorf("The expected type is '*cobra.Command', found: '%s'", cmdtype))
	}
}
