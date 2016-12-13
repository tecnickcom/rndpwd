package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func cli() (*cobra.Command, error) {

	// parse the configDir argument
	cfgCmd := new(cobra.Command)
	cfgCmd.Flags().StringVarP(&configDir, "configDir", "d", "", "Configuration directory to be added on top of the search list")
	cfgCmd.ParseFlags(os.Args)

	// configuration parameters
	cfgParams, err := getConfigParams()
	if err != nil {
		return nil, err
	}

	// overwrites the configuration parameters with the ones specified in the command line (if any)
	appParams = &cfgParams
	rootCmd := new(cobra.Command)

	rootCmd.Flags().StringVarP(&configDir, "configDir", "d", "", "Configuration directory to be added on top of the search list")
	rootCmd.Flags().StringVarP(&appParams.log.Level, "logLevel", "o", cfgParams.log.Level, "Log level: EMERGENCY, ALERT, CRITICAL, ERROR, WARNING, NOTICE, INFO, DEBUG")
	rootCmd.Flags().BoolVarP(&appParams.serverMode, "serverMode", "s", cfgParams.serverMode, "Start an HTTP RESTful API server")
	rootCmd.Flags().StringVarP(&appParams.serverAddress, "serverAddress", "u", cfgParams.serverAddress, "HTTP API address (ip:port) or just (:port)")
	rootCmd.Flags().StringVarP(&appParams.charset, "charset", "c", cfgParams.charset, "Characters to use to generate a password")
	rootCmd.Flags().IntVarP(&appParams.length, "length", "l", cfgParams.length, "Length of each password (number of characters or bytes)")
	rootCmd.Flags().IntVarP(&appParams.quantity, "quantity", "q", cfgParams.quantity, "Number of passwords to generate")
	rootCmd.Flags().StringVarP(&appParams.stats.Prefix, "statsPrefix", "p", cfgParams.stats.Prefix, "StatsD bucket prefix name")
	rootCmd.Flags().StringVarP(&appParams.stats.Network, "statsNetwork", "k", cfgParams.stats.Network, "StatsD client network type (udp or tcp)")
	rootCmd.Flags().StringVarP(&appParams.stats.Address, "statsAddress", "m", cfgParams.stats.Address, "StatsD daemon address (ip:port) or just (:port)")
	rootCmd.Flags().IntVarP(&appParams.stats.FlushPeriod, "statsFlushPeriod", "r", cfgParams.stats.FlushPeriod, "StatsD client flush period in milliseconds")

	rootCmd.Use = "rndpwd"
	rootCmd.Short = "Command-line and Web-service Random Password Generator"
	rootCmd.Long = `rndpwd is a Command-line and Web-service Random Password Generator`
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		// check values
		err := checkParams(appParams)
		if err != nil {
			return err
		}

		// initialize StatsD client (ignore errors)
		initStats(appParams.stats)
		defer stats.Close()

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
			fmt.Println(ProgramVersion)
		},
	}
	rootCmd.AddCommand(versionCmd)

	return rootCmd, nil
}
