/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"maxwolffe.com/recipeUploader/v2/service"
)

// processImagesCmd represents the processImages command
var processImagesCmd = &cobra.Command{
	Use:   "processImages",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("processImages called")
		butService := service.NewButlerService()
		path := args[0]
		fileInfo, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}

		if fileInfo.IsDir() {
			butService.ProcessRecipesInDir(path)
		} else {
			butService.ProcessSingleImage(path)
		}
	},
}

func init() {
	rootCmd.AddCommand(processImagesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// processImagesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// processImagesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
