package main

import (
	"errors"
	"regexp"
	"strings"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

// params struct contains the application parameters
type remoteConfigParams struct {
	remoteConfigProvider      string // remote configuration source ("consul", "etcd")
	remoteConfigEndpoint      string // remote configuration URL (ip:port)
	remoteConfigPath          string // remote configuration path where to search fo the configuration file ("/config/rndpwd")
	remoteConfigSecretKeyring string // path to the openpgp secret keyring used to decript the remote configuration data ("/etc/rndpwd/configkey.gpg")
}

// isEmpty returns true if all the fields are empty strings
func (rcfg remoteConfigParams) isEmpty() bool {
	return rcfg.remoteConfigProvider == "" && rcfg.remoteConfigEndpoint == "" && rcfg.remoteConfigPath == "" && rcfg.remoteConfigSecretKeyring == ""
}

// params struct contains the application parameters
type params struct {
	log           *LogData   // Log level: EMERGENCY, ALERT, CRITICAL, ERROR, WARNING, NOTICE, INFO, DEBUG
	stats         *StatsData // StatsD configuration, it is used to collect usage metrics
	serverAddress string     // HTTP address (ip:port) or just (:port)
	charset       string     // characters to use to generate a password
	length        int        // length of each password (number of characters or bytes)
	quantity      int        // number of passwords to generate
}

var configDir string
var appParams = &params{}

// getConfigParams returns the configuration parameters
func getConfigParams() (par params, err error) {
	cfg, rcfg, err := getLocalConfigParams()
	if err != nil {
		return par, err
	}
	return getRemoteConfigParams(cfg, rcfg)
}

// getLocalConfigParams returns the local configuration parameters
func getLocalConfigParams() (cfg params, rcfg remoteConfigParams, err error) {

	viper.Reset()

	// set default remote configuration values
	viper.SetDefault("remoteConfigProvider", RemoteConfigProvider)
	viper.SetDefault("remoteConfigEndpoint", RemoteConfigEndpoint)
	viper.SetDefault("remoteConfigPath", RemoteConfigPath)
	viper.SetDefault("remoteConfigSecretKeyring", RemoteConfigSecretKeyring)

	// set default configuration values

	viper.SetDefault("log.level", LogLevel)
	viper.SetDefault("log.network", LogNetwork)
	viper.SetDefault("log.address", LogAddress)

	viper.SetDefault("stats.prefix", StatsPrefix)
	viper.SetDefault("stats.network", StatsNetwork)
	viper.SetDefault("stats.address", StatsAddress)
	viper.SetDefault("stats.flush_period", StatsFlushPeriod)

	viper.SetDefault("serverAddress", ServerAddress)
	viper.SetDefault("charset", ValidCharset)
	viper.SetDefault("length", PasswordLength)
	viper.SetDefault("quantity", NumPasswords)

	// name of the configuration file without extension
	viper.SetConfigName("config")

	// configuration type
	viper.SetConfigType("json")

	if configDir != "" {
		viper.AddConfigPath(configDir)
	}

	// add local configuration paths
	for _, cpath := range ConfigPath {
		viper.AddConfigPath(cpath)
	}

	// Find and read the local configuration file (if any)
	err = viper.ReadInConfig()
	if err != nil {
		return cfg, rcfg, err
	}

	// read configuration parameters
	cfg = getViperParams()

	// support environment variables for the remote configuration
	viper.AutomaticEnv()
	viper.SetEnvPrefix(strings.Replace(ProgramName, "-", "_", -1)) // will be uppercased automatically
	envVar := []string{
		"remoteConfigProvider",
		"remoteConfigEndpoint",
		"remoteConfigPath",
		"remoteConfigSecretKeyring",
	}
	for _, ev := range envVar {
		err = viper.BindEnv(ev)
		if err != nil {
			return cfg, rcfg, err
		}
	}

	rcfg = remoteConfigParams{
		remoteConfigProvider:      viper.GetString("remoteConfigProvider"),
		remoteConfigEndpoint:      viper.GetString("remoteConfigEndpoint"),
		remoteConfigPath:          viper.GetString("remoteConfigPath"),
		remoteConfigSecretKeyring: viper.GetString("remoteConfigSecretKeyring"),
	}

	return cfg, rcfg, nil
}

// getRemoteConfigParams returns the remote configuration parameters
func getRemoteConfigParams(cfg params, rcfg remoteConfigParams) (params, error) {

	if rcfg.isEmpty() {
		return cfg, nil
	}

	viper.Reset()

	// set default configuration values

	viper.SetDefault("log.level", cfg.log.Level)
	viper.SetDefault("log.network", cfg.log.Network)
	viper.SetDefault("log.address", cfg.log.Address)

	viper.SetDefault("stats.prefix", cfg.stats.Prefix)
	viper.SetDefault("stats.network", cfg.stats.Network)
	viper.SetDefault("stats.address", cfg.stats.Address)
	viper.SetDefault("stats.flush_period", cfg.stats.FlushPeriod)

	viper.SetDefault("serverAddress", cfg.serverAddress)
	viper.SetDefault("charset", cfg.charset)
	viper.SetDefault("length", cfg.length)
	viper.SetDefault("quantity", cfg.quantity)

	// configuration type
	viper.SetConfigType("json")

	// add remote configuration provider
	var err error
	if rcfg.remoteConfigSecretKeyring == "" {
		err = viper.AddRemoteProvider(rcfg.remoteConfigProvider, rcfg.remoteConfigEndpoint, rcfg.remoteConfigPath)
	} else {
		err = viper.AddSecureRemoteProvider(rcfg.remoteConfigProvider, rcfg.remoteConfigEndpoint, rcfg.remoteConfigPath, rcfg.remoteConfigSecretKeyring)
	}
	if err == nil {
		// try to read the remote configuration (if any)
		err = viper.ReadRemoteConfig()
	}
	if err != nil {
		return cfg, err
	}

	// read configuration parameters
	return getViperParams(), nil
}

// getViperParams reads the config params via Viper
func getViperParams() params {
	return params{

		log: &LogData{
			Level:   viper.GetString("log.level"),
			Network: viper.GetString("log.network"),
			Address: viper.GetString("log.address"),
		},

		stats: &StatsData{
			Prefix:      viper.GetString("stats.prefix"),
			Network:     viper.GetString("stats.network"),
			Address:     viper.GetString("stats.address"),
			FlushPeriod: viper.GetInt("stats.flush_period"),
		},

		serverAddress: viper.GetString("serverAddress"),
		charset:       viper.GetString("charset"),
		length:        viper.GetInt("length"),
		quantity:      viper.GetInt("quantity"),
	}
}

// checkParams cheks if the configuration parameters are valid
func checkParams(prm *params) error {
	// Log
	if prm.log.Level == "" {
		return errors.New("log.level is empty")
	}
	err := prm.log.setLog()
	if err != nil {
		return err
	}

	// StatsD
	if prm.stats.Prefix == "" {
		return errors.New("stats prefix is empty")
	}
	if prm.stats.Network != "udp" && prm.stats.Network != "tcp" {
		return errors.New("stats.network must be udp or tcp")
	}
	if prm.stats.FlushPeriod < 0 {
		return errors.New("stats.flush_period must be >= 0")
	}

	// Server
	if prm.serverAddress == "" {
		return errors.New("the Server address is empty")
	}

	return checkPasswordParams(prm.quantity, prm.length, prm.charset)
}

// checkPasswordParams cheks if the password configuration parameters are valid
func checkPasswordParams(quantity int, length int, charset string) error {

	charsetLength := len(charset)
	if charsetLength < 2 || charsetLength > 92 {
		return errors.New("the charset string must contain between 2 and 92 ASCII characters")
	}
	validChr := regexp.MustCompile("[^" + regexp.QuoteMeta(ValidCharset) + "]")
	if validChr.MatchString(charset) {
		return errors.New("the charset string contains invalid characters")
	}
	if length < 2 {
		return errors.New("the length of the passwords to generate must be at least 2")
	}
	if quantity < 1 {
		return errors.New("the number of passwords to generate must be at least 1")
	}

	return nil
}
