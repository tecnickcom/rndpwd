package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestCheckParams(t *testing.T) {
	err := checkParams(&params{quantity: 1, length: 2, charset: "abc", statsNetwork: "udp", statsFlushPeriod: 100, logLevel: "info"})
	if err != nil {
		t.Error(fmt.Errorf("No errors are expected"))
	}
}

func TestCheckParamsErrorsServer(t *testing.T) {
	err := checkParams(&params{serverMode: true, serverAddress: ""})
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

func TestCheckParamsErrorsLogLevelEmpty(t *testing.T) {
	err := checkParams(&params{quantity: 1, length: 2, charset: "abc", logLevel: ""})
	if err == nil {
		t.Error(fmt.Errorf("An error was expected because the logLevel is empty"))
	}
}

func TestCheckParamsErrorsLogLevelInvalid(t *testing.T) {
	err := checkParams(&params{quantity: 1, length: 2, charset: "abc", logLevel: "INVALID"})
	if err == nil {
		t.Error(fmt.Errorf("An error was expected because the logLevel is not valid"))
	}
}

func TestGetConfigParams(t *testing.T) {
	prm, err := getConfigParams()
	if err != nil {
		t.Error(fmt.Errorf("An error was not expected: %v", err))
	}
	if prm.serverMode {
		t.Error(fmt.Errorf("Found different server mode than expected"))
	}
	if prm.serverAddress != ":8080" {
		t.Error(fmt.Errorf("Found different server address than expected"))
	}
	if prm.quantity != 10 {
		t.Error(fmt.Errorf("The expected quantity is 10, found %d", prm.quantity))
	}
	if prm.length != 32 {
		t.Error(fmt.Errorf("The expected length is 32, found %d", prm.quantity))
	}
	if prm.charset != "!#$%&()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ^_abcdefghijklmnopqrstuvwxyz~" {
		t.Error(fmt.Errorf("Fond different charset than expected"))
	}
	if prm.logLevel != "debug" {
		t.Error(fmt.Errorf("Found different logLevel than expected"))
	}
}

func TestGetLocalConfigParams(t *testing.T) {

	// test environment variables
	defer unsetRemoteConfigEnv()
	os.Setenv("RNDPWD_REMOTECONFIGPROVIDER", "consul")
	os.Setenv("RNDPWD_REMOTECONFIGENDPOINT", "127.0.0.1:98765")
	os.Setenv("RNDPWD_REMOTECONFIGPATH", "/config/rndpwd")
	os.Setenv("RNDPWD_REMOTECONFIGSECRETKEYRING", "")

	prm, rprm := getLocalConfigParams()

	if prm.serverMode {
		t.Error(fmt.Errorf("Found different server mode than expected"))
	}
	if prm.serverAddress != ":8080" {
		t.Error(fmt.Errorf("Found different server address than expected"))
	}
	if prm.quantity != 10 {
		t.Error(fmt.Errorf("The expected quantity is 10, found %d", prm.quantity))
	}
	if prm.length != 32 {
		t.Error(fmt.Errorf("The expected length is 32, found %d", prm.quantity))
	}
	if prm.charset != "!#$%&()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ^_abcdefghijklmnopqrstuvwxyz~" {
		t.Error(fmt.Errorf("Fond different charset than expected"))
	}
	if prm.logLevel != "debug" {
		t.Error(fmt.Errorf("Found different logLevel than expected"))
	}
	if rprm.remoteConfigProvider != "consul" {
		t.Error(fmt.Errorf("Found different remoteConfigProvider than expected"))
	}
	if rprm.remoteConfigEndpoint != "127.0.0.1:98765" {
		t.Error(fmt.Errorf("Found different remoteConfigEndpoint than expected"))
	}
	if rprm.remoteConfigPath != "/config/rndpwd" {
		t.Error(fmt.Errorf("Found different remoteConfigPath than expected"))
	}
	if rprm.remoteConfigSecretKeyring != "" {
		t.Error(fmt.Errorf("Found different remoteConfigSecretKeyring than expected"))
	}

	_, err := getRemoteConfigParams(prm, rprm)
	if err == nil {
		t.Error(fmt.Errorf("A remote configuration error was expected"))
	}

	rprm.remoteConfigSecretKeyring = "/etc/rndpwd/cfgkey.gpg"
	_, err = getRemoteConfigParams(prm, rprm)
	if err == nil {
		t.Error(fmt.Errorf("A remote configuration error was expected"))
	}
}

// Test real Consul provider
// To activate this define the environmental variable RNDPWD_LIVECONSUL
func TestGetConfigParamsRemote(t *testing.T) {

	enable := os.Getenv("RNDPWD_LIVECONSUL")
	if enable == "" {
		return
	}

	// test environment variables
	defer unsetRemoteConfigEnv()
	os.Setenv("RNDPWD_REMOTECONFIGPROVIDER", "consul")
	os.Setenv("RNDPWD_REMOTECONFIGENDPOINT", "127.0.0.1:8500")
	os.Setenv("RNDPWD_REMOTECONFIGPATH", "/config/rndpwd")
	os.Setenv("RNDPWD_REMOTECONFIGSECRETKEYRING", "")

	// load a specific config file just for testing
	oldCfg := ConfigPath
	viper.Reset()
	ConfigPath[0] = "wrong/path/0/"
	ConfigPath[1] = "wrong/path/1/"
	ConfigPath[2] = "wrong/path/2/"
	ConfigPath[3] = "wrong/path/3/"
	ConfigPath[4] = "wrong/path/4/"
	defer func() { ConfigPath = oldCfg }()

	prm, err := getConfigParams()
	if err != nil {
		t.Error(fmt.Errorf("An error was not expected: %v", err))
	}
	if prm.serverMode {
		t.Error(fmt.Errorf("Found different server mode than expected"))
	}
	if prm.serverAddress != ":8123" {
		t.Error(fmt.Errorf("Found different server address than expected"))
	}
	if prm.quantity != 10 {
		t.Error(fmt.Errorf("The expected quantity is 11, found %d", prm.quantity))
	}
	if prm.length != 32 {
		t.Error(fmt.Errorf("The expected length is 33, found %d", prm.quantity))
	}
	if prm.charset != "0123456789abcdefghijklmnopqrstuvwxyz" {
		t.Error(fmt.Errorf("Fond different charset than expected"))
	}
	if prm.logLevel != "debug" {
		t.Error(fmt.Errorf("Found different logLevel than expected"))
	}
}

// unsetRemoteConfigEnv clear the environmental variables used to set the remote configuration
func unsetRemoteConfigEnv() {
	os.Setenv("RNDPWD_REMOTECONFIGPROVIDER", "")
	os.Setenv("RNDPWD_REMOTECONFIGENDPOINT", "")
	os.Setenv("RNDPWD_REMOTECONFIGPATH", "")
	os.Setenv("RNDPWD_REMOTECONFIGSECRETKEYRING", "")
}
