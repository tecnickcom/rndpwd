# rndpwd

*Web-Service Random Password Generator*

This is an example project built using [gosrvlib](https://github.com/nexmoinc/gosrvlib)

[![Build](https://secure.travis-ci.org/tecnickcom/rndpwd.png?branch=main)](https://travis-ci.org/tecnickcom/rndpwd?branch=main)
[![Coverage](https://coveralls.io/repos/tecnickcom/rndpwd/badge.svg?branch=main&service=github)](https://coveralls.io/github/tecnickcom/rndpwd?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/tecnickcom/rndpwd)](https://goreportcard.com/report/github.com/tecnickcom/rndpwd)

[![Donate via PayPal](https://img.shields.io/badge/donate-paypal-87ceeb.svg)](https://www.paypal.com/cgi-bin/webscr?cmd=_donations&currency_code=GBP&business=paypal@tecnick.com&item_name=donation%20for%20rndpwd%20project)
*Please consider supporting this project by making a donation via [PayPal](https://www.paypal.com/cgi-bin/webscr?cmd=_donations&currency_code=GBP&business=paypal@tecnick.com&item_name=donation%20for%20rndpwd%20project)*

* **category**    Application
* **author**      Nicola Asuni <info@tecnick.com>
* **copyright**   2015-2022 Nicola Asuni - Tecnick.com LTD
* **license**     MIT see [LICENSE](LICENSE)
* **link**        https://github.com/tecnickcom/rndpwd

-----------------------------------------------------------------

## TOC

* [Description](#description)
* [Requirements](#requirements)
* [Quick Start](#quickstart)
* [Running all tests](#runtest)
* [Usage](#usage)
* [Configuration](#configuration)
* [Examples](#examples)
* [Logs](#logs)
* [Metrics](#metrics)
* [Profiling](#profiling)
* [OpenApi](#openapi)
* [Docker](#docker)
* [Links](#links)

-----------------------------------------------------------------

<a name="description"></a>
## Description

Web-Service Random Password Generator - gosrvlib example

-----------------------------------------------------------------

<a name="requirements"></a>
## Requirements

An additional Python program is used to check the validity of the JSON configuration files against a JSON schema:

```
sudo pip install jsonschema
```

-----------------------------------------------------------------

<a name="quickstart"></a>
## Quick Start

This project includes a Makefile that allows you to test and build the project in a Linux-compatible system with simple commands.  
All the artifacts and reports produced using this Makefile are stored in the *target* folder.  

All the packages listed in the *resources/docker/Dockerfile* file are required in order to build and test all the library options in the current environment.
Alternatively, everything can be built inside a [Docker](https://www.docker.com) container using the command "make dbuild".

To see all available options:
```
make help
```

To download all dependencies:
```
make deps
```

To update the mod file:
```
make mod
```

To format the code (please use this command before submitting any pull request):
```
make format
```

To execute all the default test builds and generate reports in the current environment:
```
make qa
```

To build the executable file:
```
make build
```

-----------------------------------------------------------------

<a name="runtest"></a>
## Running all tests

Before committing the code, please check if it passes all tests using
```bash
make qa
```

-----------------------------------------------------------------

<a name="usage"></a>
## Usage

```bash
rndpwd [flags]

Flags:

-c, --configDir  string  Configuration directory to be added on top of the search list
-f, --logFormat  string  Logging format: CONSOLE, JSON
-o, --loglevel   string  Log level: EMERGENCY, ALERT, CRITICAL, ERROR, WARNING, NOTICE, INFO, DEBUG
```

----------------------------------------------------------------

<a name="configuration"></a>
## Configuration

See [CONFIG.md](CONFIG.md).

-----------------------------------------------------------------

<a name="examples"></a>
## Examples

Once the application has being compiled with `make build`, it can be quickly tested:

```bash
target/usr/bin/rndpwd -c resources/test/etc/rndpwd
```

-----------------------------------------------------------------

<a name="logs"></a>
## Logs

This program logs the log messages in JSON format:

```
{
	"level": "info",
	"timestamp": 1595942715776382171,
	"msg": "Request",
	"program": "rndpwd",
	"version": "0.0.0",
	"release": "0",
    "hostname":"myserver",
	"request_id": "c4iah65ldoyw3hqec1rluoj93",
	"request_method": "GET",
	"request_path": "/uid",
	"request_query": "",
	"request_uri": "/uid",
	"request_useragent": "curl/7.69.1",
	"remote_ip": "[::1]:36790",
	"response_code": 200,
	"response_message": "OK",
	"response_status": "success",
	"response_data": "avxkjeyk43av"
}
```

Logs are sent to stderr by default.

The log level can be set either in the configuration or as command argument (`logLevel`).

-----------------------------------------------------------------

<a name="metrics"></a>
## Metrics

This service provides [Prometheus](https://prometheus.io/) metrics at the `/metrics` endpoint.

-----------------------------------------------------------------

<a name="profiling"></a>
## Profiling

This service provides [PPROF](https://github.com/google/pprof) profiling data at the `/pprof` endpoint.

The pprof data can be analyzed and displayed using the pprof tool:

```
go get github.com/google/pprof
```

Example:

```
pprof -seconds 10 -http=localhost:8182 http://INSTANCE_URL:PORT/pprof/profile
```

-----------------------------------------------------------------

<a name="openapi"></a>
## OpenAPI

The numberpools API is specified via the [OpenAPI 3](https://www.openapis.org/) file: `openapi.yaml`.

The openapi file can be edited using the Swagger Editor:

```
docker pull swaggerapi/swagger-editor
docker run -p 8056:8080 swaggerapi/swagger-editor
```

and pointing the Web browser to http://localhost:8056

-----------------------------------------------------------------
<a name="docker"></a>
## Docker

To build a Docker scratch container for the rndpwd executable binary execute the following command:
```
make docker
```

### Useful Docker commands

To manually create the container you can execute:
```
docker build --tag="rndpwd/rndpwddev" .
```

To log into the newly created container:
```
docker run -t -i rndpwd/rndpwddev /bin/bash
```

To get the container ID:
```
CONTAINER_ID=`docker ps -a | grep rndpwd/rndpwddev | cut -c1-12`
```

To delete the newly created docker container:
```
docker rm -f $CONTAINER_ID
```

To delete the docker image:
```
docker rmi -f rndpwd/rndpwddev
```

To delete all containers
```
docker rm $(docker ps -a -q)
```

To delete all images
```
docker rmi $(docker images -q)
```

-----------------------------------------------------------------

