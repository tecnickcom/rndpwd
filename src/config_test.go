package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func getTestCfgParams() *params {
	return &params{
		log: &LogData{
			Level:   "INFO",
			Network: "",
			Address: "",
		},
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

func TestCheckConfigParametersErrors(t *testing.T) {
	var testCases = []struct {
		fcfg  func(cfg *params) *params
		field string
	}{
		{func(cfg *params) *params { cfg.log.Level = ""; return cfg }, "log.Level"},
		{func(cfg *params) *params { cfg.log.Level = "INVALID"; return cfg }, "log.Level"},
		{func(cfg *params) *params { cfg.stats.Prefix = ""; return cfg }, "stats.Prefix"},
		{func(cfg *params) *params { cfg.stats.Network = ""; return cfg }, "stats.Network"},
		{func(cfg *params) *params { cfg.stats.FlushPeriod = -1; return cfg }, "stats.FlushPeriod"},
		{func(cfg *params) *params { cfg.serverMode = true; cfg.serverAddress = ""; return cfg }, "serverAddress"},
		{func(cfg *params) *params { cfg.quantity = 0; return cfg }, "quantity"},
		{func(cfg *params) *params { cfg.length = 0; return cfg }, "length"},
		{func(cfg *params) *params { cfg.charset = ""; return cfg }, "charset"},
		{func(cfg *params) *params {
			cfg.charset = "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"
			return cfg
		}, "charset"},
		{func(cfg *params) *params { cfg.charset = "ab cd"; return cfg }, "charset"},
	}
	for _, tt := range testCases {
		cfg := getTestCfgParams()
		err := checkParams(tt.fcfg(cfg))
		if err == nil {
			t.Error(fmt.Errorf("An error was expected because the %s field is invalid", tt.field))
		}
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
	if prm.log.Level != "DEBUG" {
		t.Error(fmt.Errorf("Found different logLevel than expected, found %s", prm.log.Level))
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
	if prm.log.Level != "DEBUG" {
		t.Error(fmt.Errorf("Found different logLevel than expected, found %s", prm.log.Level))
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
	if prm.log.Level != "DEBUG" {
		t.Error(fmt.Errorf("Found different logLevel than expected, found %s", prm.log.Level))
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
