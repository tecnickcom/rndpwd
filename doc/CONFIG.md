# Configuration Guide

The rndpwd service can load the configuration either from a local configuration file or from a single environment variable.

The local configuration file is always loaded before the remote configuration, the latter always overwrites any local setting.

If the *configDir* parameter is not specified, then the program searches for a **config.json** file in the following directories (in order of precedence):

* ./
* $HOME/rndpwd/
* /etc/rndpwd/


## Default Configuration

The default configuration file is installed in the **/etc/rndpwd/** folder (**config.json**) along with the JSON schema **config.schema.json**.


## Remote Configuration

This program supports a remote configuration source via a single environment variable (the "envvar" provider), enabling fully file-less deployments.
Arbitrary remote storage backends (e.g. Consul, Etcd, Vault, S3) can also be supported by registering a custom loader with the `config.WithRemoteLoader` option (see below).
The remote configuration source can be defined either in the local configuration file using the following parameters, or with environment variables:

The configuration fields are:

* **remoteConfigProvider**      : Remote configuration source ("envvar" or empty; any other value requires a custom loader registered via config.WithRemoteLoader)
* **remoteConfigEndpoint**      : Remote configuration URL (ip:port), passed to the custom remote loader
* **remoteConfigPath**          : Remote configuration path in which to search for the configuration data (e.g. "/config/rndpwd"), passed to the custom remote loader
* **remoteConfigSecretKeyring** : Path to an optional secret keyring used by the custom remote loader to decrypt the remote configuration data (e.g. "/etc/rndpwd/configkey.gpg")
* **remoteConfigData**          : Base64 encoded JSON configuration data to be used with the "envvar" provider (typically set via the RNDPWD_REMOTECONFIGDATA environment variable, which takes precedence over any value in the configuration file)

The equivalent environment variables are:

* RNDPWD_REMOTECONFIGPROVIDER
* RNDPWD_REMOTECONFIGENDPOINT
* RNDPWD_REMOTECONFIGPATH
* RNDPWD_REMOTECONFIGSECRETKEYRING
* RNDPWD_REMOTECONFIGDATA

### Custom Remote Loader

Providers other than "envvar" are delegated to an application-supplied function registered with the `config.WithRemoteLoader` option, so the backend client dependencies live in the application module and not in the nurago library.

The legacy Viper remote backends (consul, etcd, etcd3, firestore, nats) can be restored by plugging the ready-made `config.ViperRemoteLoader` function into the `config.WithRemoteLoader` option; the application must also register the backends manually with a blank import of the viper remote package (this is where the backend client dependencies come from):

```go
import (
	"github.com/tecnickcom/nurago/pkg/config"

	_ "github.com/spf13/viper/remote" // registers the legacy remote backends
)

err := config.Load(cmdName, configDir, envPrefix, cfg, config.WithRemoteLoader(config.ViperRemoteLoader))
```

If the blank import is missing, loading from a remote provider fails with `config.ErrViperRemoteNotRegistered`.


## Configuration Format

The configuration format is a single JSON structure with the following fields:

* **remoteConfigProvider**      : Remote configuration source ("envvar" or empty; any other value requires a custom loader registered via config.WithRemoteLoader)
* **remoteConfigEndpoint**      : Remote configuration URL (ip:port), passed to the custom remote loader
* **remoteConfigPath**          : Remote configuration path in which to search for the configuration data (e.g. "/config/rndpwd"), passed to the custom remote loader
* **remoteConfigSecretKeyring** : Path to an optional secret keyring used by the custom remote loader to decrypt the remote configuration data
* **remoteConfigData**          : Base64 encoded JSON configuration data to be used with the "envvar" provider

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

All configuration files are formatted and ordered by key using the [jq](https://github.com/jqlang/jq) tool.
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
