package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"testing"
)

/*func TestRndpwdError(t *testing.T) {
	os.Args = []string{"rndpwd", "--quantity=0"}
	//defer func() {
	//	if err := recover(); err == nil {
	//		t.Error(fmt.Errorf("An error was expected"))
	//	}
	//}()
	main()
}*/

func TestRndpwdVersion(t *testing.T) {
	os.Args = []string{"rndpwd", "version"}
	out := getMainOutput()
	match, _ := regexp.MatchString("^[\\d]+\\.[\\d]+\\.[\\d]+[\\s]*$", out)
	if !match {
		t.Error(fmt.Errorf("The expected version hs not been returned"))
	}
}

func TestRndpwd(t *testing.T) {
	os.Args = []string{"rndpwd"}
	out := getMainOutput()
	if len(out) != 330 {
		t.Error(fmt.Errorf("Expected 330 characters output (10 x 33 chars)"))
	}
}

func TestRndpwdOne(t *testing.T) {
	os.Args = []string{"rndpwd", "--quantity=1", "--length=64"}
	out := getMainOutput()
	if len(out) != 65 {
		t.Error(fmt.Errorf("Expected 1 64 character password + newline"))
	}
}

func TestRndpwdFixed(t *testing.T) {
	os.Args = []string{"rndpwd", "--quantity=1", "--length=4", "--charset=abc"}
	out := getMainOutput()
	match, _ := regexp.MatchString("^[abc]{4}[\\s]*$", out)
	if !match {
		t.Error(fmt.Errorf("Expected 'aaaa' password %s", out))
	}
}

func getMainOutput() string {
	old := os.Stdout // keep backup of the real stdout
	defer func() { os.Stdout = old }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	// execute the main function
	main()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	out := <-outC

	return out
}
