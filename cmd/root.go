// Package cmd contains all CLI commands used by the application.
package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// default const values for application.
const (
	DefaultLoggingLevel = "warn"
)

const rootCommandLongDesc = "LONG DESCRIPTION GOES HERE."

var (
	cfgFile string

	// rootCmd represents the base command when called without any subcommands.
	rootCmd = &cobra.Command{
		Use:               "pbc",
		Version:           "0.0.0",
		Short:             "pbc (punch-board-calculator) is a CLI application for calculating envelope punch positions when using a 1-2-3 punch board.",
		Long:              rootCommandLongDesc,
		Args:              cobra.NoArgs,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pbc/config)")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

// initConfig sets up Viper and Logging.
func initConfig() {
	log.Trace("initializing configuration and logging")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".pbc/config" (without extension).
		viper.AddConfigPath(home + "/.pbc")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("PBCALC")
	viper.AutomaticEnv() // read in environment variables that match

	viper.SetDefault("logging.level", DefaultLoggingLevel)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	loggingLevel, err := log.ParseLevel(viper.GetString("logging.level"))
	if err != nil {
		log.Warn("error parsing logging level: ", err)
	}

	log.SetLevel(loggingLevel)
	log.WithFields(log.Fields{"level": loggingLevel}).Debug("set log level")
}
