// Copyright (c) 2015 Tecnick.com
// Command-line and Web-service Random Password Generator
package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	rootCmd, err := cli()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("UNABLE TO START THE PROGRAM")
	}
	// execute the root command and log errors (if any)
	if err = rootCmd.Execute(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("UNABLE TO RUN THE COMMAND")
	}
}
