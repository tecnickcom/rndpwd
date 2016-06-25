// Copyright (c) 2015-2016 Nicola Asuni - Tecnick.com LTD
// Command-line and Web-service Random Password Generator
package main

func main() {
	rootCmd, err := cli()
	if err != nil {
		critLog.Fatalf("unable to start the service: %v", err)
	}
	// execute the root command and log errors (if any)
	if err = rootCmd.Execute(); err != nil {
		critLog.Fatalf("unable to start the service: %v", err)
	}
}
