package main

// AppName contains this application name
const AppName = "rndpwd"

// AppVersion contains this application version
// This is supposed to be automatically populated by the Makefile using the value from the VERSION file
var AppVersion = "0.0.0"

// ValidCharset is a string containing the valid characters for a password
const ValidCharset = "!#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnopqrstuvwxyz{|}~"

// ConfigPath list the paths where to look for configuration files (in order)
var ConfigPath = [...]string{
	"../resources/etc/" + AppName + "/",
	"./",
	"$HOME/." + AppName + "/",
	"/etc/" + AppName + "/",
}
