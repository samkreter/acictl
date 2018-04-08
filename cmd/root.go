package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/samkreter/acictl/util"
	"github.com/spf13/cobra"
)

var deploymentFile string
var resourceGroup string
var region string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "acictl",
	Short: "acictl provides a simple way to interact with Azure Container Instance.",
	Long:  `TODO`,
}

var convert = &cobra.Command{
	Use:   "convert",
	Short: "Convert a Kubernetes deployment spec into and ACI Template.",
	Long:  `Convert a Kubernetes deployment spec into and ACI Template.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := util.Convert(deploymentFile, resourceGroup, region)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var create = &cobra.Command{
	Use:   "create",
	Short: "Create an Azure Container Instance from a Kubernetes deployment spec.",
	Long:  `Create an Azure Container Instance from a Kubernetes deployment spec.`,
	Run: func(cmd *cobra.Command, args []string) {
		if resourceGroup == "" {
			log.Fatal("Must supply an Azure resource group with the -g flag.")
		}

		err := util.Create(deploymentFile, resourceGroup, region)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var delete = &cobra.Command{
	Use:   "delete",
	Short: "Delete an Azure Container Instance from a Kubernetes deployment spec.",
	Long:  `Delete an Azure Container Instance from a Kubernetes deployment spec.`,
	Run: func(cmd *cobra.Command, args []string) {
		if resourceGroup == "" {
			log.Fatal("Must supply an Azure resource group with the -g flag.")
		}

		err := util.Delete(deploymentFile, resourceGroup)
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
	RootCmd.PersistentFlags().StringVarP(&region, "region", "r", "westus", "region for aci.")
	RootCmd.PersistentFlags().StringVarP(&deploymentFile, "deployment-file", "f", "", "the kubernetes deployment file (required).")
	RootCmd.MarkFlagRequired("deployment-file")

	create.PersistentFlags().StringVarP(&resourceGroup, "resource-group", "g", "", "azure resource group for aci (required).")
	create.MarkFlagRequired("resource-group")
	delete.PersistentFlags().StringVarP(&resourceGroup, "resource-group", "g", "", "azure resource group for aci (required).")
	delete.MarkFlagRequired("resource-group")

	//Add the sub commands
	RootCmd.AddCommand(convert)
	RootCmd.AddCommand(create)
	RootCmd.AddCommand(delete)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if deploymentFile == "" {
		log.Fatal("Must supply a deployment file with the -f flag.")
	}

	//Make westus the default region
	if region == "" {
		region = "westus"
	}
}
