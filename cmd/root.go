package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var logger *zap.Logger

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xo",
	Short: "eXecution Orchestrator",
	Long:  `Connective tooling for bundle provisioning.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("A subcommand is required")
	},
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.xo.yaml)")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debugging logs")
	setupLogging()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".xo" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".xo")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func setupLogging() {
	logger, _ = zap.NewProduction()
	defer logger.Sync()
}
