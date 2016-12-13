package main

import (
	"errors"
	"regexp"

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
	log           *LogData   // Log level: EMERGENCY, ALERT, CRITICAL, ERROR, WARNING, NOTICE, INFO, DEBUG.
	serverMode    bool       // set this to true to start a RESTful API server mode
	serverAddress string     // HTTP API address for server mode (ip:port) or just (:port)
	charset       string     // characters to use to generate a password
	charsetLength int        // length of the character set in bytes
	length        int        // length of each password (number of characters or bytes)
	quantity      int        // number of passwords to generate
	stats         *StatsData // StatsD configuration, it is used to collect usage metrics
}

var configDir string
var appParams = new(params)

// getConfigParams returns the configuration parameters
func getConfigParams() (params, error) {
	cfg, rcfg := getLocalConfigParams()
	return getRemoteConfigParams(cfg, rcfg)
}

// getLocalConfigParams returns the local configuration parameters
func getLocalConfigParams() (cfg params, rcfg remoteConfigParams) {

	viper.Reset()

	// set default remote configuration values
	viper.SetDefault("remoteConfigProvider", RemoteConfigProvider)
	viper.SetDefault("remoteConfigEndpoint", RemoteConfigEndpoint)
	viper.SetDefault("remoteConfigPath", RemoteConfigPath)
	viper.SetDefault("remoteConfigSecretKeyring", RemoteConfigSecretKeyring)

	viper.SetDefault("log.level", LogLevel)
	viper.SetDefault("log.network", LogNetwork)
	viper.SetDefault("log.address", LogAddress)

	// set default configuration values
	viper.SetDefault("serverMode", ServerMode)
	viper.SetDefault("serverAddress", ServerAddress)
	viper.SetDefault("charset", ValidCharset)
	viper.SetDefault("length", PasswordLength)
	viper.SetDefault("quantity", NumPasswords)

	viper.SetDefault("stats.prefix", StatsPrefix)
	viper.SetDefault("stats.network", StatsNetwork)
	viper.SetDefault("stats.address", StatsAddress)
	viper.SetDefault("stats.flush_period", StatsFlushPeriod)

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
	viper.ReadInConfig()

	// read configuration parameters
	cfg = getViperParams()

	// support environment variables for the remote configuration
	viper.AutomaticEnv()
	viper.SetEnvPrefix(ProgramName) // will be uppercased automatically
	viper.BindEnv("remoteConfigProvider")
	viper.BindEnv("remoteConfigEndpoint")
	viper.BindEnv("remoteConfigPath")
	viper.BindEnv("remoteConfigSecretKeyring")

	rcfg = remoteConfigParams{
		remoteConfigProvider:      viper.GetString("remoteConfigProvider"),
		remoteConfigEndpoint:      viper.GetString("remoteConfigEndpoint"),
		remoteConfigPath:          viper.GetString("remoteConfigPath"),
		remoteConfigSecretKeyring: viper.GetString("remoteConfigSecretKeyring"),
	}

	return cfg, rcfg
}

// getRemoteConfigParams returns the remote configuration parameters
func getRemoteConfigParams(cfg params, rcfg remoteConfigParams) (params, error) {

	if rcfg.isEmpty() {
		return cfg, nil
	}

	viper.Reset()

	viper.SetDefault("log.level", cfg.log.Level)
	viper.SetDefault("log.network", cfg.log.Network)
	viper.SetDefault("log.address", cfg.log.Address)

	// set default configuration values
	viper.SetDefault("serverMode", cfg.serverMode)
	viper.SetDefault("serverAddress", cfg.serverAddress)
	viper.SetDefault("charset", cfg.charset)
	viper.SetDefault("length", cfg.length)
	viper.SetDefault("quantity", cfg.quantity)

	viper.SetDefault("stats.prefix", cfg.stats.Prefix)
	viper.SetDefault("stats.network", cfg.stats.Network)
	viper.SetDefault("stats.address", cfg.stats.Address)
	viper.SetDefault("stats.flush_period", cfg.stats.FlushPeriod)

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

		serverMode:    viper.GetBool("serverMode"),
		serverAddress: viper.GetString("serverAddress"),
		charset:       viper.GetString("charset"),
		length:        viper.GetInt("length"),
		quantity:      viper.GetInt("quantity"),

		stats: &StatsData{
			Prefix:      viper.GetString("stats.prefix"),
			Network:     viper.GetString("stats.network"),
			Address:     viper.GetString("stats.address"),
			FlushPeriod: viper.GetInt("stats.flush_period"),
		},
	}
}

// checkParams cheks if the configuration parameters are valid
func checkParams(prm *params) error {
	// Log
	if prm.log.Level == "" {
		return errors.New("log Level is empty")
	}
	err := prm.log.setLog()
	if err != nil {
		return err
	}

	// Server
	if prm.serverMode && prm.serverAddress == "" {
		return errors.New("The Server address is empty")
	}

	// StatsD
	if prm.stats.Prefix == "" {
		return errors.New("The stats Prefix is empty")
	}
	if prm.stats.Network != "udp" && prm.stats.Network != "tcp" {
		return errors.New("The stats Network must be udp or tcp")
	}
	if prm.stats.FlushPeriod < 0 {
		return errors.New("The stats FlushPeriod must be >= 0")
	}

	prm.charsetLength = len(prm.charset)
	if prm.charsetLength < 2 || prm.charsetLength > 92 {
		return errors.New("The charset string must contain between 2 and 92 ASCII characters")
	}
	validChr := regexp.MustCompile("[^" + regexp.QuoteMeta(ValidCharset) + "]")
	if validChr.MatchString(prm.charset) {
		return errors.New("The charset string contains invalid characters")
	}
	if prm.length < 2 {
		return errors.New("The length of the passwords to generate must be at least 2")
	}
	if prm.quantity < 1 {
		return errors.New("The number of passwords to generate must be at least 1")
	}

	return nil
}
