// Copyright (c) 2015 Nicola Asuni - Tecnick.com LTD
// Command-line password generator
package main

import (
	"errors"
	"regexp"

	"github.com/spf13/viper"
)

// params struct contains the application parameters
type params struct {
	quantity      int    // number of passwords to generate
	length        int    // lenght of each password (number of characters or bytes)
	charset       string // characters to use to generate a password
	charsetLength int    // length of the character set in bytes
}

// getConfigParams returns the configuration parameters
func getConfigParams() params {

	// name of the configuration file without extension
	viper.SetConfigName("config")

	// search paths for the configuration file ...
	viper.AddConfigPath("../resources/etc/" + AppName + "/")
	viper.AddConfigPath("./")
	viper.AddConfigPath("$HOME/." + AppName + "/")
	viper.AddConfigPath("/etc/" + AppName + "/")

	// set default configuration values
	viper.SetDefault("quantity", 1)
	viper.SetDefault("length", 16)
	viper.SetDefault("charset", ValidCharset)

	// Find and read the config file (if any)
	viper.ReadInConfig()

	// read configuration parameters
	return params{
		quantity: viper.GetInt("quantity"),
		length:   viper.GetInt("length"),
		charset:  viper.GetString("charset"),
	}
}

// checkParams cheks if the configuration parameters are valid
func checkParams(appParams *params) error {
	if appParams.quantity < 1 {
		return errors.New("The number of passwords to generate must be at least 1")
	}
	if appParams.length < 2 {
		return errors.New("The length of the passwords to generate must be at least 2")
	}
	appParams.charsetLength = len(appParams.charset)
	if appParams.charsetLength < 2 || appParams.charsetLength > 92 {
		return errors.New("The charset string must contain between 2 and 92 ASCII characters")
	}
	validChr := regexp.MustCompile("[^" + regexp.QuoteMeta(ValidCharset) + "]")
	if validChr.MatchString(appParams.charset) {
		return errors.New("The charset string contains invalid characters")
	}
	return nil
}
