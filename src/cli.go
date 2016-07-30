package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func cli() (*cobra.Command, error) {

	// configuration parameters
	cfgParams, err := getConfigParams()
	if err != nil {
		return nil, err
	}

	// set the root command
	rootCmd := new(cobra.Command)

	// overwrites the configuration parameters with the ones specified in the command line (if any)
	rootCmd.Flags().BoolVarP(&appParams.serverMode, "serverMode", "s", cfgParams.serverMode, "Start an HTTP RESTful API server")
	rootCmd.Flags().StringVarP(&appParams.serverAddress, "serverAddress", "u", cfgParams.serverAddress, "HTTP API address (ip:port) or just (:port)")
	rootCmd.Flags().StringVarP(&appParams.charset, "charset", "c", cfgParams.charset, "Characters to use to generate a password")
	rootCmd.Flags().IntVarP(&appParams.length, "length", "l", cfgParams.length, "Length of each password (number of characters or bytes)")
	rootCmd.Flags().IntVarP(&appParams.quantity, "quantity", "q", cfgParams.quantity, "Number of passwords to generate")
	rootCmd.Flags().StringVarP(&appParams.logLevel, "logLevel", "o", cfgParams.logLevel, "Log level: NONE, EMERGENCY, ALERT, CRITICAL, ERROR, WARNING, NOTICE, INFO, DEBUG")

	rootCmd.Use = "rndpwd"
	rootCmd.Short = "Command-line and Web-service Random Password Generator"
	rootCmd.Long = `rndpwd is a Command-line and Web-service Random Password Generator`
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		// check values
		err := checkParams(appParams)
		if err != nil {
			return err
		}
		if appParams.serverMode {
			// start the HTTP server
			return startServer(appParams.serverAddress)
		}

		// generate and print the passwords
		for _, psw := range getAllPassword(appParams) {
			fmt.Println(psw)
		}
		return nil
	}

	// sub-command to print the version
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print this program version",
		Long:  `print this program version`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(ServiceVersion)
		},
	}
	rootCmd.AddCommand(versionCmd)

	return rootCmd, nil
}
