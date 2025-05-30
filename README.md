<!-- Space: tecnickcom -->
<!-- Parent: Projects -->
<!-- Title: rndpwd -->
<!-- Attachment: doc/images/logo.png -->

# rndpwd

*Web-Service Random Password Generator*

![rndpwd logo](doc/images/logo.png)

* **category:**    Application
* **copyright:**   2015-2025 tecnickcom
* **license:**     [LICENSE](https://github.com/tecnickcom/rndpwd/blob/main/LICENSE)
* **cvs:**         https://github.com/tecnickcom/rndpwd

[![check](https://github.com/tecnickcom/rndpwd/actions/workflows/check.yaml/badge.svg)](https://github.com/tecnickcom/rndpwd/actions/workflows/check.yaml)
[![Coverage Status](https://coveralls.io/repos/github/tecnickcom/rndpwd/badge.svg?branch=main)](https://coveralls.io/github/tecnickcom/rndpwd?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/tecnickcom/rndpwd)](https://goreportcard.com/report/github.com/tecnickcom/rndpwd)

[![Donate via PayPal](https://img.shields.io/badge/donate-paypal-87ceeb.svg)](https://www.paypal.com/cgi-bin/webscr?cmd=_donations&currency_code=GBP&business=paypal@tecnick.com&item_name=donation%20for%20rndpwd%20project)
*Please consider supporting this project by making a donation via [PayPal](https://www.paypal.com/cgi-bin/webscr?cmd=_donations&currency_code=GBP&business=paypal@tecnick.com&item_name=donation%20for%20rndpwd%20project)*


----------

## TOC
* [Description](#description)
* [Documentation](#documentation)
	* [public](#documentation_public)
		* [General](#documentation_public_general)
* [Development](#development)
* [Deployment](#deployment)

----------

<a name="description"></a>
## Description
Web-Service Random Password Generator


----------



<a name="documentation"></a>
## Documentation
<a name="documentation_public"></a>
* public
	<a name="documentation_public_general"></a>
	* General  
	_General project documentation_
		* [GitHup project page](https://github.com/tecnickcom/rndpwd)


----------



<a name="development"></a>
## Development
### TOC

* [Style and Conventions](#style)
* [Requirements](#requirements)
* [Quick Start](#quickstart)
* [Running all tests](#runtest)
* [Documentation](#gendoc)
* [Usage](#usage)
* [Configuration](CONFIG.md)
* [Examples](#examples)
* [Logs](#logs)
* [Metrics](#metrics)
* [Profiling](#profiling)
* [OpenAPI](#openapi)
* [Docker](#docker)


<a name="style"></a>
## Style and Conventions

For the general style and conventions, please refer to external documents:
https://github.com/uber-go/guide/blob/master/style.md


<a name="quickstart"></a>
## Developers' Quick Start

To quickly get started with this project, follow these steps:

1. Ensure you ahev installed the latest Go version and Python3 for some extra tests.
1. Clone the repository: `git clone https://github.com/tecnickcom/rndpwd.git`.
2. Change into the project directory: `cd rndpwd`.
3. Install the required dependencies and test everything: `DEVMODE=LOCAL make x`.

Now you are ready to start developing with rndpwd!


This project includes a *Makefile* that allows you to test and build the project in a Linux-compatible system with simple commands.  
All the artifacts and reports produced using this *Makefile* are stored in the *target* folder.  

Alternatively, everything can be built inside a [Docker](https://www.docker.com) container using the command `make dbuild` that uses the environment defined at `resources/docker/Dockerfile.dev`.

To see all available options:
```bash
make help
```


<a name="runtest"></a>
## Running all tests

Before committing the code, please check if it passes all tests using
```bash
make x
```

that is an alias for:
```bash
DEVMODE=LOCAL make format clean mod deps generate qa build docker dockertest
```


<a name="gendoc"></a>
## Documentation

The `README.md` and `doc/RUNBOOK.md` documentation files are generated using the source templates in `doc/src` via `make gendoc` command.

To update links and common information edit the file `doc/src/config.yaml` in YAML format.
The schema of the configuration file is defined by the JSON schema: `doc/src/config.schema.json`.
The document templates are defined by the `*.tmpl` files in [gomplate](https://docs.gomplate.ca)-compatible format.

To regenerate the static documentation file:
```bash
make gendoc
```


<a name="usage"></a>
## Usage

```bash
rndpwd [flags]

Flags:

-c, --configDir  string  Configuration directory to be added on top of the search list
-f, --logFormat  string  Logging format: CONSOLE, JSON
-o, --loglevel   string  Log level: EMERGENCY, ALERT, CRITICAL, ERROR, WARNING, NOTICE, INFO, DEBUG
```

<a name="examples"></a>
## Examples

Once the application has being compiled with `make build`, it can be quickly tested:

```bash
target/usr/bin/rndpwd -c resources/test/etc/rndpwd
```


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


<a name="metrics"></a>
## Metrics

This service provides [Prometheus](https://prometheus.io/) metrics at the `/metrics` endpoint.


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


<a name="openapi"></a>
## OpenAPI

The rndpwd API is specified via the [OpenAPI 3](https://www.openapis.org/) file: `openapi.yaml`.

The openapi file can be edited using the Swagger Editor:

```
docker pull swaggerapi/swagger-editor
docker run -p 8056:8080 swaggerapi/swagger-editor
```

and pointing the Web browser to http://localhost:8056


<a name="docker"></a>
## Docker

To build a Docker scratch container for the rndpwd executable binary execute the following command:
```
make docker
```

### Useful Docker commands

To manually create the container you can execute:
```
docker build --tag="tecnickcom/rndpwddev" .
```

To log into the newly created container:
```
docker run -t -i tecnickcom/rndpwddev /bin/bash
```

To get the container ID:
```
CONTAINER_ID=`docker ps -a | grep tecnickcom/rndpwddev | cut -c1-12`
```

To delete the newly created docker container:
```
docker rm -f $CONTAINER_ID
```

To delete the docker image:
```
docker rmi -f tecnickcom/rndpwddev
```

To delete all containers
```
docker rm $(docker ps -a -q)
```

To delete all images
```
docker rmi $(docker images -q)
```


----------

<a name="deployment"></a>
## Deployment
### Deployment in Production

Add here information on how to deploy in production.


----------


