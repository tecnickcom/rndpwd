package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/spf13/cobra"
)

var stopServerChannel chan bool

var emptyParamCases = []string{
	"--serverMode=true --serverAddress=",
	"--charset=",
	"--length=",
	"--quantity=",
	"--logLevel=",
	"--logLevel=INVALID",
	"--statsNetwork=INVALID",
	"--statsFlushPeriod=-1",
}

func TestCliEmptyParamError(t *testing.T) {
	for _, param := range emptyParamCases {
		os.Args = []string{"rndpwd", param}
		cmd, err := cli()
		if err != nil {
			t.Error(fmt.Errorf("An error wasn't expected: %v", err))
			return
		}
		if cmdtype := reflect.TypeOf(cmd).String(); cmdtype != "*cobra.Command" {
			t.Error(fmt.Errorf("The expected type is '*cobra.Command', found: '%s'", cmdtype))
			return
		}

		old := os.Stderr // keep backup of the real stdout
		defer func() { os.Stderr = old }()
		os.Stderr = nil

		// execute the main function
		if err := cmd.Execute(); err == nil {
			t.Error(fmt.Errorf("An error was expected"))
		}
	}
}
func TestCli(t *testing.T) {
	os.Args = []string{
		ProgramName,
		"--configDir=resources/test/etc/rndpwd",
		"--serverMode=true",
		"--serverAddress=:8765",
		"--statsPrefix=rndpwdtest",
		"--statsNetwork=udp",
		"--statsAddress=:8125",
		"--statsFlushPeriod=100",
	}
	cmd, err := cli()
	if err != nil {
		t.Error(fmt.Errorf("An error wasn't expected: %v", err))
		return
	}
	if cmdtype := reflect.TypeOf(cmd).String(); cmdtype != "*cobra.Command" {
		t.Error(fmt.Errorf("The expected type is '*cobra.Command', found: '%s'", cmdtype))
		return
	}

	old := os.Stderr // keep backup of the real stdout
	defer func() { os.Stderr = old }()
	os.Stderr = nil

	// add an endpoint to test the panic handler
	routes = append(routes,
		Route{
			"GET",
			"/panic",
			triggerPanic,
			"TRIGGER PANIC",
		})
	defer func() { routes = routes[:len(routes)-1] }()

	// use two separate channels for server and client testing
	var twg sync.WaitGroup
	startTestServer(t, cmd, &twg)
	startTestClient(t)
	twg.Wait()
}

func startTestServer(t *testing.T, cmd *cobra.Command, twg *sync.WaitGroup) {

	stopServerChannel = make(chan bool)

	twg.Add(1)
	go func() {
		defer twg.Done()

		chp := make(chan error, 1)
		go func() {
			chp <- cmd.Execute()
		}()

		stopped := false
		for {
			select {
			case err, ok := <-chp:
				if ok && !stopped && err != nil {
					stopServerChannel <- true
					t.Error(fmt.Errorf("An error was not expected: %v", err))
				}
				return
			case <-stopServerChannel:
				stopped = true
				if server != nil {
					server.Close() // this triggers the cmd.Execute error
				} else {
					return
				}
			}
		}
	}()

	// wait for the server to start
	time.Sleep(500 * time.Millisecond)
}

func startTestClient(t *testing.T) {

	// check if the server is running
	select {
	case stop, ok := <-stopServerChannel:
		if ok && stop {
			return
		}
	default:
		break
	}

	defer func() { stopServerChannel <- true }()

	testEndPoint(t, "GET", "/", "", 200)
	testEndPoint(t, "GET", "/status", "", 200)

	// error conditions
	testEndPoint(t, "GET", "/INVALID", "", 404) // NotFound
	testEndPoint(t, "DELETE", "/", "", 405)     // MethodNotAllowed
	testEndPoint(t, "GET", "/panic", "", 500)   // PanicHandler

	// valid endpoints
	testEndPoint(t, "GET", "/password", "", 200)
	testEndPoint(t, "GET", "/password?quantity=5", "", 200)
	testEndPoint(t, "GET", "/password?quantity=5&length=13", "", 200)
	testEndPoint(t, "GET", "/password?quantity=5&length=13&charset=abcdef0123456789", "", 200)

	// test query errors
	testEndPoint(t, "GET", "/password?quantity=0", "", 400)
	testEndPoint(t, "GET", "/password?length=0", "", 400)
	testEndPoint(t, "GET", "/password?charset=", "", 400)
}

// triggerPanic triggers a Panic
func triggerPanic(rw http.ResponseWriter, hr *http.Request, ps httprouter.Params) {
	panic("TEST PANIC")
}

// return true if the input is a JSON
func isJSON(s []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(s, &js) == nil
}

func testEndPoint(t *testing.T, method string, path string, data string, code int) {
	var payload = []byte(data)
	req, err := http.NewRequest(method, fmt.Sprintf("http://127.0.0.1:8765%s", path), bytes.NewBuffer(payload))
	if err != nil {
		t.Error(fmt.Errorf("An error was not expected: %v", err))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(fmt.Errorf("An error was not expected: %v", err))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != code {
		t.Error(fmt.Errorf("The espected status code is %d, found %d", code, resp.StatusCode))
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(fmt.Errorf("An error was not expected: %v", err))
		return
	}
	if !isJSON(body) {
		t.Error(fmt.Errorf("The body is not a JSON"))
	}
}
