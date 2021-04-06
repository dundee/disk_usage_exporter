package cmd

import (
	"fmt"
	"os"

	exporter "github.com/dundee/gdu-prometheus-exporter/exporter"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	address    string
	ignoreDirs []string
)

var rootCmd = &cobra.Command{
	Use:   "gdu-prometheus-exporter",
	Short: "Prometheus exporter for detailed disk usage info",
	Long: `Prometheus exporter analysing disk usage of the filesystem
and reporting which directories consume what space.`,
	Run: func(cmd *cobra.Command, args []string) {
		go exporter.RunAnalysis(ignoreDirs)
		exporter.RunServer(address)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gdu-prometheus-exporter.yaml)")
	rootCmd.PersistentFlags().StringVarP(&address, "bind-address", "b", "0.0.0.0:9108", "Address to bind to")
	rootCmd.PersistentFlags().StringSliceVarP(&ignoreDirs, "ignore-dirs", "i", []string{"/proc", "/dev", "/sys", "/run"}, "Absolute paths to ignore (separated by comma)")
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

		// Search config in home directory with name ".gdu-prometheus-exporter" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gdu-prometheus-exporter")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
