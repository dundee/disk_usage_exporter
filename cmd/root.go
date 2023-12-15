package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/dundee/disk_usage_exporter/build"
	"github.com/dundee/disk_usage_exporter/exporter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "disk_usage_exporter",
	Short: "Prometheus exporter for detailed disk usage info",
	Long: `Prometheus exporter analysing disk usage of the filesystem
and reporting which directories consume what space.`,
	Run: func(cmd *cobra.Command, args []string) {
		printHeader()

		paths := transformMultipaths(viper.GetStringMapString("multi-paths"))
		if len(paths) == 0 {
			paths[filepath.Clean(viper.GetString("analyzed-path"))] = viper.GetInt("dir-level")
		}

		e := exporter.NewExporter(
			(paths),
			viper.GetBool("follow-symlinks"),
		)
		e.SetIgnoreDirPaths(viper.GetStringSlice("ignore-dirs"))

		if viper.GetString("mode") == "file" {
			e.WriteToTextfile(viper.GetString("output-file"))
			log.Info("Done - exiting.")
		} else {
			e.RunServer(viper.GetString("bind-address"))
		}
	},
}

// Execute runs the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	flags := rootCmd.PersistentFlags()
	flags.StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.disk_usage_exporter.yaml)")
	flags.StringP("mode", "m", "http", "Expose method - either 'file' or 'http'")
	flags.StringP("bind-address", "b", "0.0.0.0:9995", "Address to bind to")
	flags.StringP("output-file", "f", "./disk-usage-exporter.prom", "Target file to store metrics in")
	flags.StringP("analyzed-path", "p", "/", "Path where to analyze disk usage")
	flags.IntP("dir-level", "l", 2, "Directory nesting level to show (0 = only selected dir)")
	flags.StringSliceP("ignore-dirs", "i", []string{"/proc", "/dev", "/sys", "/run", "/var/cache/rsnapshot"}, "Absolute paths to ignore (separated by comma)")
	flags.BoolP(
		"follow-symlinks", "L", false,
		"Follow symlinks for files, i.e. show the size of the file to which symlink points to (symlinks to directories are not followed)",
	)
	flags.StringToString("multi-paths", map[string]string{}, "Multiple paths where to analyze disk usage, in format /path1=level1,/path2=level2,...")

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

		// Search config in home directory with name ".disk_usage_exporter" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".disk_usage_exporter")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func printHeader() {
	log.Printf("Disk Usage Prometheus Exporter %s	build date: %s	sha1: %s	Go: %s	GOOS: %s	GOARCH: %s",
		build.BuildVersion,
		build.BuildDate,
		build.BuildCommitSha,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}

func transformMultipaths(multiPaths map[string]string) map[string]int {
	paths := make(map[string]int, len(multiPaths))
	for path, level := range multiPaths {
		l, err := strconv.Atoi(level)
		if err != nil {
			log.Fatalf("Invalid level for path %s: %s", path, level)
		}
		paths[filepath.Clean(path)] = l
	}
	return paths
}
