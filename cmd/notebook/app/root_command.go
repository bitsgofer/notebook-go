package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	klog "k8s.io/klog/v2"
)

var (
	rootCmd = &cobra.Command{
		Use:          "notebook",
		Short:        "Notebook is a DIY blogging tool for learning",
		SilenceUsage: true,
	}
	configFile string
)

func setupSubcommands() {
	rootCmd.AddCommand(checkDependenciesCmd)
	rootCmd.AddCommand(developCmd)
	rootCmd.AddCommand(renderCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func setupRootCmdFlags() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(
		&configFile,
		"config",
		"",
		"config file (default: [REPO-ROOT]/.notebook)",
	)

	viper.SetDefault("theme", "bitsgofer")
	viper.SetDefault("contentDir", "content")
	viper.SetDefault("outputDir", "_public_html")

	viper.SetDefault("devServer.dataDir", "_public_html")
	viper.SetDefault("devServer.addr", "localhost:8080")

	// include klog's flags into CLI
	fs := flag.NewFlagSet("", flag.PanicOnError)
	klog.InitFlags(fs)
	rootCmd.PersistentFlags().AddGoFlagSet(fs)
}

func initConfig() {
	// config should be load from a YAML file, e.g: ./.notebook
	viper.SetConfigType("yaml")
	viper.SetConfigName(".notebook")

	if configFile != "" { // override config file
		klog.V(3).InfoS("use config file from flag", "config", configFile)

		viper.SetConfigFile(configFile)
	} else {
		klog.V(3).InfoS("search for config file (.notebook) in current dir")

		currentDir, err := os.Getwd()
		cobra.CheckErr(err)

		viper.AddConfigPath(currentDir)
	}

	// load config file
	err := viper.ReadInConfig()
	if err != nil {
		cobra.CheckErr(fmt.Errorf("cannot read config file; err= %w", err))
		return
	}
}
