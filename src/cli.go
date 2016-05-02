package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func cli() *cobra.Command {

	// configuration parameters
	cfgParams := getConfigParams()

	// set the root command
	rootCmd := new(cobra.Command)

	// overwrites the configuration parameters with the ones specified in the command line (if any)
	rootCmd.Flags().BoolVarP(&appParams.server, "server", "s", cfgParams.server, "Start an HTTP RESTful API server")
	rootCmd.Flags().StringVarP(&appParams.httpaddr, "httpaddr", "u", cfgParams.httpaddr, "HTTP API address (ip:port) or just (:port)")
	rootCmd.Flags().IntVarP(&appParams.quantity, "quantity", "q", cfgParams.quantity, "Number of passwords to generate")
	rootCmd.Flags().IntVarP(&appParams.length, "length", "l", cfgParams.length, "Length of each password (number of characters or bytes)")
	rootCmd.Flags().StringVarP(&appParams.charset, "charset", "c", cfgParams.charset, "Characters to use to generate a password")

	rootCmd.Use = "rndpwd"
	rootCmd.Short = "Command-line and Web-service Random Password Generator"
	rootCmd.Long = `rndpwd is a Command-line and Web-service Random Password Generator`
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		// check values
		err := checkParams(appParams)
		if err != nil {
			return err
		}
		if appParams.server {
			// start the HTTP server
			startServer(appParams.httpaddr)
		} else {
			// generate and print the passwords
			for _, psw := range getAllPassword(appParams) {
				fmt.Println(psw)
			}
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
