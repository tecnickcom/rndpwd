# Configuration Guide

The rndpwd service can load the configuration either from a local configuration file or remotely via [Consul](https://www.consul.io/), [Etcd](https://github.com/coreos/etcd) or a single Environmental Variable.

The local configuration file is always loaded before the remote configuration, the latter always overwrites any local setting.

If the *configDir* parameter is not specified, then the program searches for a **config.json** file in the following directories (in order of precedence):

* ./
* $HOME/rndpwd/
* /etc/rndpwd/


## Default Configuration

The default configuration file is installed in the **/etc/rndpwd/** folder (**config.json**) along with the JSON schema **config.schema.json**.


## Remote Configuration

This program supports secure remote configuration via Consul, Etcd or single environment variable.
The remote configuration server can be defined either in the local configuration file using the following parameters, or with environment variables:

The configuration fields are:

* **remoteConfigProvider**      : Remote configuration source ("consul", "etcd", "envvar")
* **remoteConfigEndpoint**      : Remote configuration URL (ip:port)
* **remoteConfigPath**          : Remote configuration path in which to search for the configuration file (e.g. "/config/rndpwd")
* **remoteConfigSecretKeyring** : Path to the [OpenPGP](http://openpgp.org/) secret keyring used to decrypt the remote configuration data (e.g. "/etc/rndpwd/configkey.gpg"); if empty a non secure connection will be used instead
* **remoteConfigData**          : Base64 encoded JSON configuration data to be used with the "envvar" provider

The equivalent environment variables are:

* RNDPWD_REMOTECONFIGPROVIDER
* RNDPWD_REMOTECONFIGENDPOINT
* RNDPWD_REMOTECONFIGPATH
* RNDPWD_REMOTECONFIGSECRETKEYRING
* RNDPWD_REMOTECONFIGDATA


## Configuration Format

The configuration format is a single JSON structure with the following fields:

* **remoteConfigProvider**      : Remote configuration source ("consul", "etcd", "envvar")
* **remoteConfigEndpoint**      : Remote configuration URL (ip:port)
* **remoteConfigPath**          : Remote configuration path in which to search for the configuration file (e.g. "/config/rndpwd")
* **remoteConfigSecretKeyring** : Path to the openpgp secret keyring used to decrypt the remote configuration data (e.g. "/etc/rndpwd/configkey.gpg"); if empty a non secure connection will be used instead

* **enabled**: Enable or disable the service

* **log**:  Logging settings
    * **format**:  Logging format: CONSOLE, JSON
    * **level**:   Defines the default log level: EMERGENCY, ALERT, CRITICAL, ERROR, WARNING, NOTICE, INFO, DEBUG
    * **network**: (OPTIONAL) Network type used by the Syslog (i.e. udp or tcp)
    * **address**: (OPTIONAL) Network address of the Syslog daemon (ip:port) or just (:port)

* **servers**: Configuration for exposed servers
    * **monitoring**: Monitoring HTTP server
        * **address**: HTTP address (ip:port) or just (:port)
        * **timeout**: HTTP request timeout [seconds]
    * **public**: *Public HTTP server*
        * **address**: HTTP address (ip:port) or just (:port)
        * **timeout**: HTTP request timeout [seconds]

* **shutdown_timeout**: Time to wait on exit for a graceful shutdown [seconds]

* **clients**: Configuration for external service clients
    * **ipify**:  ipify service client
        * **address**:  Base URL of the service
        * **timeout**:  HTTP client timeout [seconds]

* **random**: *Settings for the random generator*
    * **charset**:  *String containing the valid characters for a password*
    * **length**:   *Length of each password (number of characters or bytes)*
    * **quantity**: *Number of passwords to return*


## Formatting Configuration

All configuration files are formatted and ordered by key using the [jq](https://github.com/stedolan/jq) tool.
For example:

```cat 'resources/etc/rndpwd/config.schema.json' | jq -S .```


## Validating Configuration

The [jv](https://github.com/santhosh-tekuri/jsonschema) program can be used to check the validity of the configuration file against the JSON schema.
It can be installed via:

```
go install github.com/santhosh-tekuri/jsonschema/cmd/jv@latest
```

Example usage:

```
jv resources/etc/rndpwd/config.schema.json resources/etc/rndpwd/config.json
```
