package main

import (
	"errors"
	"regexp"

	"github.com/spf13/viper"
)

// params struct contains the application parameters
type params struct {
	server        bool   // set this to true to start a RESTful API server mode
	httpaddr      string // HTTP API address for server mode (ip:port) or just (:port)
	quantity      int    // number of passwords to generate
	length        int    // length of each password (number of characters or bytes)
	charset       string // characters to use to generate a password
	charsetLength int    // length of the character set in bytes
}

// application parameters
var appParams = new(params)

// getConfigParams returns the configuration parameters
func getConfigParams() params {

	// name of the configuration file without extension
	viper.SetConfigName("config")

	// add configuration paths
	for _, cpath := range ConfigPath {
		viper.AddConfigPath(cpath)
	}

	// set default configuration values
	viper.SetDefault("server", false)
	viper.SetDefault("httpaddr", HttpAddress)
	viper.SetDefault("quantity", NumPasswords)
	viper.SetDefault("length", PasswordLength)
	viper.SetDefault("charset", ValidCharset)

	// Find and read the config file (if any)
	viper.ReadInConfig()

	// read configuration parameters
	return params{
		server:   viper.GetBool("server"),
		httpaddr: viper.GetString("httpaddr"),
		quantity: viper.GetInt("quantity"),
		length:   viper.GetInt("length"),
		charset:  viper.GetString("charset"),
	}
}

// checkParams cheks if the configuration parameters are valid
func checkParams(prm *params) error {
	if prm.server && prm.httpaddr == "" {
		return errors.New("The Server address is empty")
	}
	if prm.quantity < 1 {
		return errors.New("The number of passwords to generate must be at least 1")
	}
	if prm.length < 2 {
		return errors.New("The length of the passwords to generate must be at least 2")
	}
	prm.charsetLength = len(prm.charset)
	if prm.charsetLength < 2 || prm.charsetLength > 92 {
		return errors.New("The charset string must contain between 2 and 92 ASCII characters")
	}
	validChr := regexp.MustCompile("[^" + regexp.QuoteMeta(ValidCharset) + "]")
	if validChr.MatchString(prm.charset) {
		return errors.New("The charset string contains invalid characters")
	}
	return nil
}
