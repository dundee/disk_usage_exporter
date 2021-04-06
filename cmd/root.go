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
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "gdu-prometheus-exporter",
	Short: "Prometheus exporter for detailed disk usage info",
	Long: `Prometheus exporter analysing disk usage of the filesystem
and reporting which directories consume what space.`,
	Run: func(cmd *cobra.Command, args []string) {
		go exporter.RunAnalysis(
			viper.GetString("analyzed-path"),
			viper.GetStringSlice("ignore-dirs"),
			viper.GetInt("dir-level"),
		)
		exporter.RunServer(
			viper.GetString("bind-address"),
		)
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
	flags := rootCmd.PersistentFlags()
	flags.StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.gdu-prometheus-exporter.yaml)")
	flags.StringP("bind-address", "b", "0.0.0.0:9108", "Address to bind to")
	flags.StringP("analyzed-path", "p", "/", "Path where to analyze disk usage")
	flags.IntP("dir-level", "l", 1, "Directory nesting level to show (0 = only selected dir)")
	flags.StringSliceP("ignore-dirs", "i", []string{"/proc", "/dev", "/sys", "/run"}, "Absolute paths to ignore (separated by comma)")

	viper.BindPFlags(flags)
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
