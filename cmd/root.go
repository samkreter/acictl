package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/samkreter/acictl/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configPath string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "acictl",
	Short: "acictl provides a simple way to interact with Azure Container Instance.",
	Long:  `TODO`,
}

var convert = &cobra.Command{
	Use:   "convert",
	Short: "Convert a Kubernetes deployment spec into and ACI Template.",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("Must provide a Kubernetes deployment file.")
		}

		err := util.Convert(args[0])
		if err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//RootCmd.PersistentFlags().StringVar(&configPath, "config", filepath.Join(home, ".dockdev"), "config file (default is $HOME/.dockdev)")

	//Add the sub commands
	RootCmd.AddCommand(convert)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// if config != "" {
	// 	// Use config file from the flag.
	// 	viper.SetConfigFile(config)
	// } else {
	// 	// Search config in home directory with name ".dockdev" (without extension).
	// 	viper.AddConfigPath(home)
	// 	viper.SetConfigName(".dockdev")
	// }

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
