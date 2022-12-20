/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/maxwolffe/butler-cli/v2/data"
	"github.com/maxwolffe/butler-cli/v2/service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// processImagesCmd represents the processImages command
var processImagesCmd = &cobra.Command{
	Use:   "processImages",
	Short: "Given a path to an image or directory of images, upload them to your Butler Queue and return the processing results. ",
	Long: `Given a path to an image or directory of images, upload them to your Butler Queue and return the processing results.

	$ butler-cli processImages /Users/maxwolffe/Desktop/ThymeChurros.png
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

		sugaredLogger.Infoln("Processing Images...")
		butService := service.NewButlerService(sugaredLogger)
		path := args[0]
		fileInfo, err := os.Stat(path)
		if err != nil {
			sugaredLogger.Fatal(err)
		}

		var documents []data.Document

		if fileInfo.IsDir() {
			documents, _ = butService.ProcessRecipesInDir(path)
		} else {
			documents, _ = butService.ProcessSingleImage(path)
		}
		if csvOutputFilePath != "" {
			service.GenerateCsv(documents, csvOutputFilePath)
		}
		sugaredLogger.Infoln("Done.")
	},
}

func init() {
	rootCmd.AddCommand(processImagesCmd)
}
