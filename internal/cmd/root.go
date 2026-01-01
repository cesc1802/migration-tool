package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile            string
	envName            string
	version, commit, date string
)

var rootCmd = &cobra.Command{
	Use:   "migrate-tool",
	Short: "Database migration CLI tool",
	Long:  `Cross-platform database migration tool with single-file up/down support.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

// SetVersionInfo sets version information for the CLI
func SetVersionInfo(v, c, d string) {
	version, commit, date = v, c, d
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ./migrate-tool.yaml)")
	rootCmd.PersistentFlags().StringVar(&envName, "env", "dev", "environment name")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("migrate-tool")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()
}
