// Copyright (c) 2015-2016 Nicola Asuni - Tecnick.com LTD
// Command-line and Web-service Random Password Generator
package main

import (
	log "github.com/Sirupsen/logrus"
)

func main() {
	rootCmd, err := cli()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("unable to start the service")
	}
	// execute the root command and log errors (if any)
	if err = rootCmd.Execute(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("unable to start the service")
	}
}
