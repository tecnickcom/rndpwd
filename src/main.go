// Copyright (c) 2015-2016 Nicola Asuni - Tecnick.com LTD
// Command-line and Web-service Random Password Generator
package main

import (
	"log"
)

func main() {
	rootCmd := cli()
	// execute the root command and log errors (if any)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
