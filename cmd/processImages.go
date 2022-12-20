/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/maxwolffe/butler-cli/v2/data"
	"github.com/maxwolffe/butler-cli/v2/service"
	"github.com/spf13/cobra"
)

// processImagesCmd represents the processImages command
var processImagesCmd = &cobra.Command{
	Use:   "processImages",
	Short: "Given a path to an image or directory of images, upload them to your Butler Queue and return the processing results. ",
	Long: `Given a path to an image or directory of images, upload them to your Butler Queue and return the processing results.

	$ butler-cli processImages /Users/maxwolffe/Desktop/ThymeChurros.png
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("processImages called")
		butService := service.NewButlerService()
		path := args[0]
		fileInfo, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
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
	},
}

func init() {
	// TODO collate responses into csv based on flag, instead of assuming by default.
	rootCmd.AddCommand(processImagesCmd)
}
