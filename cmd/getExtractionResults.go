/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/maxwolffe/butler-cli/v2/service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// getExtractionResultsCmd represents the getExtractionResults command
var getExtractionResultsCmd = &cobra.Command{
	Use:   "getExtractionResults",
	Short: "Given an uploadID return Butler results.",
	Long: `Given an uploadID (available from the website or upload response), return butler results.
	
	$ ./butler-cli getExtractionResults 5cc7d32e-b5cc-48b2-a125-90b3acfe0e1c
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO - avoid setting up logger in multiple places
		options := make([]zap.Option, 0)

		if !verbose {
			options = append(options, zap.IncreaseLevel(zap.InfoLevel))
		}

		// Set up logger
		logger, _ := zap.NewDevelopment(options...)
		defer logger.Sync() // flushes buffer, if any
		sugaredLogger = logger.Sugar()

		sugaredLogger.Infoln("Getting Extraction Results...")
		butService := service.NewButlerService(sugaredLogger)

		uploadId := args[0]
		documents, _ := butService.GetExtractionResults(uploadId)
		if csvOutputFilePath != "" {
			service.GenerateCsv(documents.Response.Items, csvOutputFilePath)
		}
		sugaredLogger.Infoln("Done.")
	},
}

func init() {
	rootCmd.AddCommand(getExtractionResultsCmd)
}
