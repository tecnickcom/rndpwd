package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestCli(t *testing.T) {
	os.Args = []string{"rndpwd", "--quantity=0"}
	cmd := cli()
	if cmdtype := reflect.TypeOf(cmd).String(); cmdtype != "*cobra.Command" {
		t.Error(fmt.Errorf("The expected type is '*cobra.Command', found: '%s'", cmdtype))
	}

	old := os.Stderr // keep backup of the real stdout
	defer func() { os.Stderr = old }()
	os.Stderr = nil

	// execute the main function
	if err := cmd.Execute(); err == nil {
		t.Error(fmt.Errorf("An error was expected"))
	}
}
