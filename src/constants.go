package main

// AppName defines this application name
const AppName = "rndpwd"

// AppVersion set this application version
// This is supposed to be automatically populated by the Makefile using the value from the VERSION file
// (-ldflags '-X main.ServiceVersion=${VERSION}')
var AppVersion = "0.0.0"

// ValidCharset is a string containing the valid characters for a password
const ValidCharset = "!#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnopqrstuvwxyz{|}~"

// NumPasswords is the default number of passwords to return
const NumPasswords = 10

// PasswordLength is the default length of each password (number of characters or bytes)
const PasswordLength = 16

// ServerAddress is the default HTTP API URL (ip:port) or just (:port)
const HttpAddress = ":8081"

// ConfigPath list the paths where to look for configuration files (in order)
var ConfigPath = [...]string{
	"../resources/etc/" + AppName + "/",
	"./",
	"$HOME/." + AppName + "/",
	"/etc/" + AppName + "/",
}
