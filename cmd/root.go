package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// This defines the base command when you type 'sentinel' with no arguments
var rootCmd = &cobra.Command{
	Use:   "sentinel",
	Short: "A headless workflow orchestrator",
	Long: `Sentinel is a local event bus that connects your SaaS tools 
to your terminal via webhooks and polling.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Sentinel is installed! Use 'go run main.go start' to run the daemon.")
	},
}

// Execute adds all child commands to the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}