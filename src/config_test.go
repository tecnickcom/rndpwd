package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func getTestCfgParams() *params {
	return &params{
		logLevel: "INFO",
		quantity: 1,
		length:   2,
		charset:  "abc",
		stats: &StatsData{
			Prefix:      "dummy-test",
			Network:     "udp",
			Address:     ":8125",
			FlushPeriod: 100,
		},
	}
}

func TestCheckParams(t *testing.T) {
	err := checkParams(getTestCfgParams())
	if err != nil {
		t.Error(fmt.Errorf("No errors are expected: %v", err))
	}
}

func TestCheckParamsErrorsLogLevelEmpty(t *testing.T) {
	cfg := getTestCfgParams()
	cfg.logLevel = ""
	err := checkParams(cfg)
	if err == nil {
		t.Error(fmt.Errorf("An error was expected because logLevel is empty"))
	}
}

func TestCheckParamsErrorsLogLevelInvalid(t *testing.T) {
	cfg := getTestCfgParams()
	cfg.logLevel = "INVALID"
	err := checkParams(cfg)
	if err == nil {
		t.Error(fmt.Errorf("An error was expected because logLevel is invalid"))
	}
}

func TestCheckParamsErrorsStatsPrefix(t *testing.T) {
	cfg := getTestCfgParams()
	cfg.stats.Prefix = ""
	err := checkParams(cfg)
	if err == nil {
		t.Error(fmt.Errorf("An error was expected because the stats Prefix is empty"))
	}
}

func TestCheckParamsErrorsStatsNetwork(t *testing.T) {
	cfg := getTestCfgParams()
	cfg.stats.Network = ""
	err := checkParams(cfg)
	if err == nil {
		t.Error(fmt.Errorf("An error was expected because the stats Network is empty"))
	}
}

func TestCheckParamsErrorsStatsFlushPeriod(t *testing.T) {
	cfg := getTestCfgParams()
	cfg.stats.FlushPeriod = -1
	err := checkParams(cfg)
	if err == nil {
		t.Error(fmt.Errorf("An error was expected because the stats FlushPeriod is negative"))
	}
}

func TestCheckParamsErrorsServer(t *testing.T) {
	cfg := getTestCfgParams()
	cfg.serverMode = true
	cfg.serverAddress = ""
	err := checkParams(cfg)
	if err == nil {
		t.Error(fmt.Errorf("An error was expected because the server address is empty"))
	}
}

func TestCheckParamsErrorsQuantity(t *testing.T) {
	cfg := getTestCfgParams()
	cfg.quantity = 0
	err := checkParams(cfg)
	if err == nil {
		t.Error(fmt.Errorf("An error was expected because the quantity is <= 0"))
	}
}

func TestCheckParamsErrorsLength(t *testing.T) {
	cfg := getTestCfgParams()
	cfg.length = 0
	err := checkParams(cfg)
	if err == nil {
		t.Error(fmt.Errorf("An error was expected because the length is <= 0"))
	}
}

func TestCheckParamsErrorsCharsetLength(t *testing.T) {
	cfg := getTestCfgParams()
	cfg.charset = ""
	err := checkParams(cfg)
	if err == nil {
		t.Error(fmt.Errorf("An error was expected because the charset length is < 2"))
	}
}

func TestCheckParamsErrorsCharsetLengthMax(t *testing.T) {
	cfg := getTestCfgParams()
	cfg.charset = "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"
	err := checkParams(cfg)
	if err == nil {
		t.Error(fmt.Errorf("An error was expected because the charset length is > 92"))
	}
}

func TestCheckParamsErrorsValidCharset(t *testing.T) {
	cfg := getTestCfgParams()
	cfg.charset = "ab cd"
	err := checkParams(cfg)
	if err == nil {
		t.Error(fmt.Errorf("An error was expected because the charset contains an invalid character"))
	}
}

func TestGetConfigParams(t *testing.T) {
	prm, err := getConfigParams()
	if err != nil {
		t.Error(fmt.Errorf("An error was not expected: %v", err))
	}
	if prm.serverMode {
		t.Error(fmt.Errorf("Found different server mode than expected, found %v", prm.serverMode))
	}
	if prm.serverAddress != ":8080" {
		t.Error(fmt.Errorf("Found different server address than expected, found %s", prm.serverAddress))
	}
	if prm.quantity != 10 {
		t.Error(fmt.Errorf("The expected quantity is 10, found %d", prm.quantity))
	}
	if prm.length != 32 {
		t.Error(fmt.Errorf("The expected length is 32, found %d", prm.length))
	}
	if prm.charset != "!#$%&()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ^_abcdefghijklmnopqrstuvwxyz~" {
		t.Error(fmt.Errorf("Fond different charset than expected, found %s", prm.charset))
	}
	if prm.logLevel != "debug" {
		t.Error(fmt.Errorf("Found different logLevel than expected, found %s", prm.logLevel))
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
		t.Error(fmt.Errorf("Found different server mode than expected, found %v", prm.serverMode))
	}
	if prm.serverAddress != ":8080" {
		t.Error(fmt.Errorf("Found different server address than expected, found %s", prm.serverAddress))
	}
	if prm.quantity != 10 {
		t.Error(fmt.Errorf("The expected quantity is 10, found %d", prm.quantity))
	}
	if prm.length != 32 {
		t.Error(fmt.Errorf("The expected length is 32, found %d", prm.length))
	}
	if prm.charset != "!#$%&()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ^_abcdefghijklmnopqrstuvwxyz~" {
		t.Error(fmt.Errorf("Fond different charset than expected, found %s", prm.charset))
	}
	if prm.logLevel != "debug" {
		t.Error(fmt.Errorf("Found different logLevel than expected, found %s", prm.logLevel))
	}
	if rprm.remoteConfigProvider != "consul" {
		t.Error(fmt.Errorf("Found different remoteConfigProvider than expected, found %s", rprm.remoteConfigProvider))
	}
	if rprm.remoteConfigEndpoint != "127.0.0.1:98765" {
		t.Error(fmt.Errorf("Found different remoteConfigEndpoint than expected, found %s", rprm.remoteConfigEndpoint))
	}
	if rprm.remoteConfigPath != "/config/rndpwd" {
		t.Error(fmt.Errorf("Found different remoteConfigPath than expected, found %s", rprm.remoteConfigPath))
	}
	if rprm.remoteConfigSecretKeyring != "" {
		t.Error(fmt.Errorf("Found different remoteConfigSecretKeyring than expected, found %s", rprm.remoteConfigSecretKeyring))
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
		t.Error(fmt.Errorf("Found different server mode than expected, found %v", prm.serverMode))
	}
	if prm.serverAddress != ":8123" {
		t.Error(fmt.Errorf("Found different server address than expected, found %s", prm.serverAddress))
	}
	if prm.quantity != 10 {
		t.Error(fmt.Errorf("The expected quantity is 11, found %d", prm.quantity))
	}
	if prm.length != 32 {
		t.Error(fmt.Errorf("The expected length is 33, found %d", prm.length))
	}
	if prm.charset != "0123456789abcdefghijklmnopqrstuvwxyz" {
		t.Error(fmt.Errorf("Found different charset than expected, found %s", prm.charset))
	}
	if prm.logLevel != "debug" {
		t.Error(fmt.Errorf("Found different logLevel than expected, found %s", prm.logLevel))
	}
}

func TestCliWrongConfigError(t *testing.T) {

	// test environment variables
	defer unsetRemoteConfigEnv()
	os.Setenv("RNDPWD_REMOTECONFIGPROVIDER", "consul")
	os.Setenv("RNDPWD_REMOTECONFIGENDPOINT", "127.0.0.1:999999")
	os.Setenv("RNDPWD_REMOTECONFIGPATH", "/config/wrong")
	os.Setenv("RNDPWD_REMOTECONFIGSECRETKEYRING", "")

	// load a specific config file just for testing
	oldCfg := ConfigPath
	viper.Reset()
	for k := range ConfigPath {
		ConfigPath[k] = "wrong/path/"
	}
	defer func() { ConfigPath = oldCfg }()

	_, err := cli()
	if err == nil {
		t.Error(fmt.Errorf("An error was expected"))
		return
	}
}

// unsetRemoteConfigEnv clear the environmental variables used to set the remote configuration
func unsetRemoteConfigEnv() {
	os.Setenv("RNDPWD_REMOTECONFIGPROVIDER", "")
	os.Setenv("RNDPWD_REMOTECONFIGENDPOINT", "")
	os.Setenv("RNDPWD_REMOTECONFIGPATH", "")
	os.Setenv("RNDPWD_REMOTECONFIGSECRETKEYRING", "")
}
