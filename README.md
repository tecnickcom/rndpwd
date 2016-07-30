# rndpwd
*Command-line and Web-service Random Password Generator*

[![Master Branch](https://img.shields.io/badge/-master:-gray.svg)](https://github.com/tecnickcom/rndpwd/tree/master)
[![Master Build Status](https://secure.travis-ci.org/tecnickcom/rndpwd.png?branch=master)](https://travis-ci.org/tecnickcom/rndpwd?branch=master)
[![Master Coverage Status](https://coveralls.io/repos/tecnickcom/rndpwd/badge.svg?branch=master&service=github)](https://coveralls.io/github/tecnickcom/rndpwd?branch=master)

[![Develop Branch](https://img.shields.io/badge/-develop:-gray.svg)](https://github.com/tecnickcom/rndpwd/tree/develop)
[![Develop Build Status](https://secure.travis-ci.org/tecnickcom/rndpwd.png?branch=develop)](https://travis-ci.org/tecnickcom/rndpwd?branch=develop)
[![Develop Coverage Status](https://coveralls.io/repos/tecnickcom/rndpwd/badge.svg?branch=develop&service=github)](https://coveralls.io/github/tecnickcom/rndpwd?branch=develop)
[![Go Report Card](https://goreportcard.com/badge/github.com/tecnickcom/rndpwd)](https://goreportcard.com/report/github.com/tecnickcom/rndpwd)

[![Donate via PayPal](https://img.shields.io/badge/donate-paypal-87ceeb.svg)](https://www.paypal.com/cgi-bin/webscr?cmd=_donations&currency_code=GBP&business=paypal@tecnick.com&item_name=donation%20for%20rndpwd%20project)
*Please consider supporting this project by making a donation via [PayPal](https://www.paypal.com/cgi-bin/webscr?cmd=_donations&currency_code=GBP&business=paypal@tecnick.com&item_name=donation%20for%20rndpwd%20project)*

* **category**    Application
* **author**      Nicola Asuni <info@tecnick.com>
* **copyright**   2015-2016 Nicola Asuni - Tecnick.com LTD
* **license**     MIT (see LICENSE)
* **link**        https://github.com/tecnickcom/rndpwd

## Description

Command-line and Web-service Random Password Generator

This is a full example of command-line and Web-service GO language project using a Makefile that integrates targets for common QA tasks and packaging, including RPM, Debian and Docker.

## Getting started

This application is written in GO language, please refere to the guides in https://golang.org for getting started.

This project include a Makefile that allows you to test and build the project with simple commands.
To see all available options:
```bash
make help
```

To buil dthe project

```bash
make build
```

## Running all tests

Before committing the code, please check if it passes all tests using
```bash
make qa
```

Other make options are available install this library globally and build RPM and DEB packages.
Please check all the available options using `make help`.


## Usage

```bash
rndpwd [flags]

Flags:

-s, --serverMode: Set this to true to start an HTTP RESTful API server
-u, --serverAddress="8080": HTTP API address for server mode (ip:port) or just (:port)
-c, --charset="!#.0123456789@ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz": Characters to use to generate a password
-l, --length=16: Length of each password (number of characters or bytes)
-q, --quantity=1: Number of passwords to generate
-o, --loglevel=INFO: Log level: NONE, EMERGENCY, ALERT, CRITICAL, ERROR, WARNING, NOTICE, INFO, DEBUG
```

## Configuration

If no command-line parameters are specified, then the ones in the configuration file (**config.json**) will be used.
The configuration files can be stored in the current directory or in any of the following (in order of precedence):
* ./
* config/
* $HOME/rndpwd/
* /etc/rndpwd/

This service also support secure remote configuration via Consul or Etcd.
The remote configuration server can be defined either in the local configuration file using the following parameters, or with environment variables:

* **remoteConfigProvider** : remote configuration source ("consul", "etcd");
* **remoteConfigEndpoint** : remote configuration URL (ip:port);
* **remoteConfigPath** : remote configuration path where to search fo the configuration file (e.g. "/config/rndpwd");
* **remoteConfigSecretKeyring** : path to the openpgp secret keyring used to decript the remote configuration data (e.g. "/etc/rndpwd/configkey.gpg"); if empty a non secure connection will be used instead;

The equivalent environment variables are:

* RNDPWD_REMOTECONFIGPROVIDER
* RNDPWD_REMOTECONFIGENDPOINT
* RNDPWD_REMOTECONFIGPATH
* RNDPWD_REMOTECONFIGSECRETKEYRING


## Server Mode

When the server mode is enabled a RESTful HTTP JSON API server will listen on the configured **address:port** for the following entry points:

| ENTRY POINT                   | METHOD | DESCRIPTION                                                    |
|:----------------------------- |:------:|:-------------------------------------------------------------- |
|<nobr> /                </nobr>| GET    |<nobr> return a list of available entry points and tests </nobr>|
|<nobr> /status          </nobr>| GET    |<nobr> check the server status                           </nobr>|
|<nobr> /password        </nobr>| GET    |generate new passwords as configured; charset, length and quantity can be specified as query parameters |


## Examples

Once the application has being compiled with `make build`, it can be quickly tested:

Generate 10 passwords with 32 characters:
```bash
target/usr/bin/rndpwd --length=32 --quantity=10
```

Generate 10 passwords with 8 characters using only numbers:
```bash
target/usr/bin/rndpwd --charset="0123456789" --length=8 --quantity=10
```

Start aan HTTP RESTful API server listening on port 8080
```bash
target/usr/bin/rndpwd --server --serverAddress=:8080
```

## Logs

This service logs the log messages in the *Stderr* or *Stdout* using the syslog prefixes:


| PREFIX                            | LEVEL | DESCRIPTION                                                                            | OUTPUT |
|:--------------------------------- |:-----:|:-------------------------------------------------------------------------------------- |:------:|
|<nobr> [EMERGENCY] [rndpwd] </nobr>|   0   |<nobr> **Emergency**: System is unusable                                         </nobr>| Stderr |
|<nobr> [ALERT] [rndpwd]     </nobr>|   1   |<nobr> **Alert**: Should be corrected immediately                                </nobr>| Stderr |
|<nobr> [CRITICAL] [rndpwd]  </nobr>|   2   |<nobr> **Critical**: Critical conditions                                         </nobr>| Stderr |
|<nobr> [ERROR] [rndpwd]     </nobr>|   3   |<nobr> **Error**: Error conditions                                               </nobr>| Stderr |
|<nobr> [WARNING] [rndpwd]   </nobr>|   4   |<nobr> **Warning**: May indicate that an error will occur if action is not taken </nobr>| Stderr |
|<nobr> [NOTICE] [rndpwd]    </nobr>|   5   |<nobr> **Notice**: Events that are unusual, but not error conditions             </nobr>| Stderr |
|<nobr> [INFO] [rndpwd]      </nobr>|   6   |<nobr> **Informational**: Normal operational messages that require no action     </nobr>| Stderr |
|<nobr> [DEBUG] [rndpwd]     </nobr>|   7   |<nobr> **Debug**: Information useful to developers for debugging the application </nobr>| Stderr |



## Developer(s) Contact

* Nicola Asuni <info@tecnick.com>
