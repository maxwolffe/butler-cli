/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var csvOutputFilePath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Short: "A cli for interacting with butlerlabs.ai",
	Long:  `A cli for interacting with butlerlabs.ai.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&csvOutputFilePath, "outputPath", "o", "", "Specifies the path to write the output CSV to. If empty, no CSV is output. (.e.g -o '/Users/maxwolffe/Desktop/output.csv') ")
}
