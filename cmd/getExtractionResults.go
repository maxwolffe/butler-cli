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
	
	$ ./butler-cli getExtractionResults 5cc7d32e-b5cc-48b2-a125-90b3acfe0e1c
	`,
	Run: func(cmd *cobra.Command, args []string) {
		butService := service.NewButlerService()

		uploadId := args[0]
		butService.GetExtractionResults(uploadId)
	},
}

func init() {
	rootCmd.AddCommand(getExtractionResultsCmd)
}
