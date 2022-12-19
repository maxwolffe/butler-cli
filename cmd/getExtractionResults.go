/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/maxwolffe/butler-cli/v2/service"
	"github.com/spf13/cobra"
)

// getExtractionResultsCmd represents the getExtractionResults command
var getExtractionResultsCmd = &cobra.Command{
	Use:   "getExtractionResults",
	Short: "Given an uploadID return Butler results.",
	Long: `Given an uploadID (available from the website or upload response), return butler results.
	

	`,
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
