package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "buckmate",
	Short:   "Deploy to S3 buckets with ease",
	Long:    ``,
	Version: "1.0.1",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("env", "e", "", "Specifies which config to apply - directory name that contains environment specific Config.yaml and files to be copied.")
	rootCmd.PersistentFlags().StringP("path", "p", "", "Specifies path to the directory that contains buckmate directory with Deployment.yaml config.")
}
