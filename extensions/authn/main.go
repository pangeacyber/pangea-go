package main

import (
	"extensions/authn/imports"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

// Execute executes the root command.
func initConfig(cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".pangea")
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

/*
This script help customer to import their existing users to the pangea
  - Should support reading data from csv/json
  - Create users with random password generation
  - Generate output with stats
  - Generate file with temp password for each user
  - Should support a dry-run (dry-run does not call endpoint to create user but rest of workflow is generated)
*/
func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	// Don't forget to close and flush the logs before exiting.
	defer logger.Sync() //nolint:errcheck

	var cfgFile string
	var importFile string
	var token, domain string
	rootCmd := &cobra.Command{
		Use:   "authn",
		Short: "A authn cli tool",
		Long:  `AuthN is a CLI library for Go that empowers customer to import users to pangea.`,
	}
	importCmd := &cobra.Command{
		Use:   "import [filePath]",
		Short: "Import users from given csv or json file",
		Long:  "One time user import to the pangea.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Importing Users")
			if token == "" {
				fmt.Print("Enter Pangea Token: ")
				tokenBytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))
				if err != nil {
					fmt.Printf("failed to read token, err=%v\n", err)
					os.Exit(1)
				}
				token = string(tokenBytes)
			}

			err := imports.ImportUsers(token, domain, importFile)
			if err != nil {
				fmt.Printf("failed to import, err=%v\n", err)
				os.Exit(1)
			}
		},
	}
	// Ideally we should split into different files as per sub commands
	// TODO: Move to different file when we have more commands
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c",
		"", "config file (default is $HOME/.pangea.yaml)")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "pangea token (default is PANGEA_TOKEN env)")
	rootCmd.PersistentFlags().StringVarP(&domain, "domain", "d", "", "pangea domain (default is PANGEA_DOMAIN env)")
	rootCmd.MarkPersistentFlagRequired("domain")
	// Binding flag with viper
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("pangea.token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("pangea.domain", rootCmd.PersistentFlags().Lookup("domain"))
	importCmd.PersistentFlags().StringVarP(&importFile, "importFile", "i", "",
		"import user csv or json file")
	importCmd.MarkPersistentFlagRequired("importFile")
	initConfig(cfgFile)
	rootCmd.AddCommand(importCmd)
	rootCmd.Execute()
}
