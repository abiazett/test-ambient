package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"flag"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
)

var cfgFile string
var kubeconfig string
var namespace string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "odh",
	Short: "OpenShift AI CLI",
	Long: `Command Line Interface for OpenShift AI.
This CLI allows you to manage training jobs and models in OpenShift AI.`,
	// No run function - the root command itself does nothing
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.odh/config.yaml)")

	// Kubeconfig flag
	defaultKubeconfig := ""
	if home := homedir.HomeDir(); home != "" {
		defaultKubeconfig = filepath.Join(home, ".kube", "config")
	}
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", defaultKubeconfig, "(optional) absolute path to the kubeconfig file")

	// Namespace flag
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "namespace to use")

	// Add verbose flag
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home := homedir.HomeDir()
		if home == "" {
			fmt.Println("Error: unable to locate home directory")
			os.Exit(1)
		}

		// Search config in home directory with name ".odh" (without extension).
		viper.AddConfigPath(filepath.Join(home, ".odh"))
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Set up logging based on verbose flag
	verbose, err := rootCmd.Flags().GetBool("verbose")
	if err == nil && verbose {
		klog.InitFlags(nil)
		_ = flag.Set("v", "4")
	}

	// Try to load kubeconfig
	if _, err := clientcmd.LoadFromFile(kubeconfig); err != nil {
		klog.V(4).Infof("Warning: Could not load kubeconfig file: %v", err)
	}
}