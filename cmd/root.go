package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chat-demo",
	Short: "A little demo of using Amazon Bedrock to chat with LLMs",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("region", "r", "us-east-1", "set the AWS region")
	rootCmd.PersistentFlags().Float64("temperature", 1, "temperature setting")
	rootCmd.PersistentFlags().Float64("topP", 0.999, "topP setting")
	rootCmd.PersistentFlags().Int("topK", 250, "topK setting")
	rootCmd.PersistentFlags().Int("max-tokens", 500, "max tokens to sample")

}
