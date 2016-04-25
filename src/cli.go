// Copyright (c) 2015 Nicola Asuni - Tecnick.com LTD
// Command-line password generator
package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func cli() *cobra.Command {

	// configuration parameters
	cfgParams := getConfigParams()

	// application parameters
	appParams := new(params)

	// set the root command
	rootCmd := new(cobra.Command)

	// overwrites the configuration parameters with the ones specified in the command line (if any)
	rootCmd.Flags().IntVarP(&appParams.quantity, "quantity", "q", cfgParams.quantity, "Number of passwords to generate")
	rootCmd.Flags().IntVarP(&appParams.length, "length", "l", cfgParams.length, "Length of each password (number of characters or bytes)")
	rootCmd.Flags().StringVarP(&appParams.charset, "charset", "c", cfgParams.charset, "Characters to use to generate a password")

	rootCmd.Use = "rndpwd"
	rootCmd.Short = "Generate random passwords"
	rootCmd.Long = `rndpwd is a command-line tool to generate random passwords`
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		// check values
		err := checkParams(appParams)
		if err != nil {
			return err
		}
		// generate and print the passwords
		for i := 0; i < appParams.quantity; i++ {
			fmt.Println(getNewPassword(appParams))
		}
		return nil
	}

	// sub-command to print the version
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print this program version",
		Long:  `print this program version`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(AppVersion)
		},
	}
	rootCmd.AddCommand(versionCmd)

	return rootCmd
}
