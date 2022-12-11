/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"maxwolffe.com/recipeUploader/v2/service"
)

// getExtractionResultsCmd represents the getExtractionResults command
var getExtractionResultsCmd = &cobra.Command{
	Use:   "getExtractionResults",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		butService := service.NewButlerService()

		uploadId := args[0]
		butService.GetExtractionResults(uploadId)
	},
}

func init() {
	rootCmd.AddCommand(getExtractionResultsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getExtractionResultsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getExtractionResultsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
