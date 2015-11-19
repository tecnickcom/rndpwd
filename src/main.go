// Copyright (c) 2015 Nicola Asuni - Tecnick.com LTD
// Command-line password generator
package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {

	var numPasswords int
	var passwordLength int
	var charset string
	validCharset := "!#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnopqrstuvwxyz{|}~"

	// name of the configuration file without extension
	viper.SetConfigName("config")

	// search paths for the configuration file ...
	viper.AddConfigPath("../../../resources/etc/" + appName + "/")
	viper.AddConfigPath("$HOME/." + appName + "/")
	viper.AddConfigPath("/etc/" + appName + "/")

	// set default configuration values
	viper.SetDefault("numPasswords", 1)
	viper.SetDefault("passwordLength", 16)
	viper.SetDefault("charset", validCharset)

	// read the configuration file
	viper.ReadInConfig()

	// read configuration parameters
	numPasswordsCfg := viper.GetInt("numPasswords")
	passwordLengthCfg := viper.GetInt("passwordLength")
	charsetCfg := viper.GetString("charset")

	// set the root command
	var rootCmd = &cobra.Command{
		Use:   "rndpwd",
		Short: "Generate random passwords",
		Long:  `rndpwd is a command-line tool to generate random passwords`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// check values
			if numPasswords < 1 {
				return errors.New("The number of passwords to generate must be at least 1")
			}
			if passwordLength < 2 {
				return errors.New("The length of the passwords to generate must be at least 2")
			}
			charsetLength := len(charset)
			if charsetLength < 2 || charsetLength > 92 {
				return errors.New("The charset string must contain between 2 and 92 ASCII characters")
			}
			validChr := regexp.MustCompile("[^" + regexp.QuoteMeta(validCharset) + "]")
			if validChr.MatchString(charset) {
				return errors.New("The charset string contains invalid characters")
			}

			// generate and print the passwords
			for i := 0; i < numPasswords; i++ {
				fmt.Println(getNewPassword(passwordLength, charset, charsetLength))
			}

			return nil
		},
	}

	// overwrites the configuration parameters with the ones specified in the command line (if any)
	rootCmd.Flags().IntVarP(&passwordLength, "length", "l", passwordLengthCfg, "Password length (number of characters)")
	rootCmd.Flags().IntVarP(&numPasswords, "number", "n", numPasswordsCfg, "Number of passwords to generate")
	rootCmd.Flags().StringVarP(&charset, "charset", "c", charsetCfg, "String of valid characters for a password")

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

	// execute the root command and log errors (if any)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
