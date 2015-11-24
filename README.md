# rndpwd

*Please consider supporting this project by making a donation to <paypal@tecnick.com>*

* **category**    Application
* **author**      Nicola Asuni <info@tecnick.com>
* **copyright**   2015-2015 Nicola Asuni - Tecnick.com LTD
* **license**     MIT (see LICENSE)
* **link**        https://github.com/tecnickcom/rndpwd

## Status
* **MASTER**: [![Build Status](https://secure.travis-ci.org/tecnickcom/rndpwd.png?branch=master)](https://travis-ci.org/tecnickcom/rndpwd?branch=master)
[![Coverage Status](https://coveralls.io/repos/tecnickcom/rndpwd/badge.svg?branch=master&service=github)](https://coveralls.io/github/tecnickcom/rndpwd?branch=master)
* **DEVELOP**: [![Build Status](https://secure.travis-ci.org/tecnickcom/rndpwd.png?branch=develop)](https://travis-ci.org/tecnickcom/rndpwd?branch=develop)
[![Coverage Status](https://coveralls.io/repos/tecnickcom/rndpwd/badge.svg?branch=develop&service=github)](https://coveralls.io/github/tecnickcom/rndpwd?branch=develop)


## Description

Command-line Random Password Generator.

This is an example of GO language project using a Makefile that integrates targets for common QA tasks and packaging, including RPM and Debian. 

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
-c, --charset="!#.0123456789@ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz": Characters to use to generate a password
-l, --length=16: Lenght of each password (number of characters or bytes)
-q, --quantity=1: Number of passwords to generate
```

## Examples

Once the application has being compiled with `make build`, it can be quickly tested:

Generate 10 passwords with 32 characters:
```bash
target/usr/bin/rndpwd --quantity=10 --length=32
```

Generate 10 passwords with 8 characters using only numbers:
```bash
target/usr/bin/rndpwd --quantity=10 --length=8 --charset="0123456789"
```

## Developer(s) Contact

* Nicola Asuni <info@tecnick.com>
