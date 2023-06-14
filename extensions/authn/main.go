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

// Init config to read the configuration from config file as well
// TODO - test it out if we have config with ~/.pangea_token.yaml
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
	// Register as global so I can refer it via zap.L()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	// Flags
	var importFile, mappingFile, cfgFile, token, domain string
	var isDryRun bool
	// Root command
	rootCmd := &cobra.Command{
		Use:   "authn",
		Short: "A authn cli tool",
		Long:  `AuthN is a CLI library for Go that empowers customer to import users to pangea.`,
	}

	// Sub command
	// Ideally we should split into different files as per sub commands
	// TODO: Move to different file when we have more commands
	importCmd := &cobra.Command{
		Use:   "import [filePath]",
		Short: "Import users from given csv or json file",
		Long:  "One time user import to the pangea.",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("Importing Users")
			if token == "" {
				// TODO: Add multiple retries
				fmt.Print("Enter Pangea Token: ")
				tokenBytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))
				if err != nil {
					logger.Fatal("failed to read token", zap.Error(err))
				}
				token = string(tokenBytes)
			}
			err := imports.ImportUsers(token, domain, importFile, mappingFile, isDryRun)
			if err != nil {
				logger.Fatal("failed to import, err=%v\n", zap.Error(err))
			}
		},
	}
	// Root persist flags which sub command inherits
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c",
		"", "config file (default is $HOME/.pangea.yaml)")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "",
		"pangea token (default is PANGEA_TOKEN env)")
	rootCmd.PersistentFlags().StringVarP(&domain, "domain", "d", "",
		"pangea domain (default is PANGEA_DOMAIN env)")
	rootCmd.MarkPersistentFlagRequired("domain")
	// Binding flag with viper config so env variable will work as well
	viper.BindPFlag("pangea.token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("pangea.domain", rootCmd.PersistentFlags().Lookup("domain"))

	// Import cmd flags
	importCmd.PersistentFlags().StringVarP(&importFile, "importFile", "i", "",
		"import user csv or json file")
	importCmd.PersistentFlags().BoolVarP(&isDryRun, "dry-run", "f", false,
		"mimic run import workflow (it does not make api call to create users). Default is false")
	// Flag local to this command
	importCmd.LocalFlags().StringVarP(&mappingFile, "fieldsMapping", "m", "",
		"Fields mapping file to map source provider to pangea")

	// Init config
	initConfig(cfgFile)

	// Register command and run
	rootCmd.AddCommand(importCmd)
	rootCmd.Execute()
}
