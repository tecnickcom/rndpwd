package main

// ProgramName defines this application name
const ProgramName = "rndpwd"

// ProgramVersion set this application version
// This is supposed to be automatically populated by the Makefile using the value from the VERSION file
// (-ldflags '-X main.ProgramVersion=${VERSION}')
var ProgramVersion = "0.0.0"

// ProgramRelease contains this program release number (or build number)
// This is automatically populated by the Makefile using the value from the RELEASE file
var ProgramRelease = "0"

// ValidCharset is a string containing the valid characters for a password
const ValidCharset = "!#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnopqrstuvwxyz{|}~"

// NumPasswords is the default number of passwords to return
const NumPasswords = 10

// PasswordLength is the default length of each password (number of characters or bytes)
const PasswordLength = 16

// ServerMode is the default HTTP server mode (on = true)
const ServerMode = false

// ServerAddress is the default HTTP API URL (ip:port) or just (:port)
const ServerAddress = ":8081"

// ConfigPath list the paths where to look for configuration files (in order)
var ConfigPath = [...]string{
	"../resources/test/etc/" + ProgramName + "/",
	"./",
	"config/",
	"$HOME/." + ProgramName + "/",
	"/etc/" + ProgramName + "/",
}

// LogLevel defines the default log level: NONE, EMERGENCY, ALERT, CRITICAL, ERROR, WARNING, NOTICE, INFO, DEBUG
const LogLevel = "info"

// RemoteConfigProvider is the remote configuration source ("consul", "etcd")
const RemoteConfigProvider = ""

// RemoteConfigEndpoint is the remote configuration URL (ip:port)
const RemoteConfigEndpoint = ""

// RemoteConfigPath is the remote configuration path where to search fo the configuration file ("/config/rndpwd")
const RemoteConfigPath = ""

// RemoteConfigSecretKeyring is the path to the openpgp secret keyring used to decript the remote configuration data ("/etc/rndpwd/configkey.gpg")
const RemoteConfigSecretKeyring = ""
