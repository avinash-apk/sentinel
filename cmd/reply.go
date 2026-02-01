package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	// MAKE SURE THIS PATH MATCHES YOUR GO.MOD
	"github.com/avinash-apk/sentinel/pkg/postmaster" 
)

var replyCmd = &cobra.Command{
	Use:   "reply [platform] [id] [message]",
	Short: "Send a reply to GitHub, Slack, or Discord",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		platform := args[0]
		id := args[1]
		message := args[2]

		var sender postmaster.Messenger
		// We remove 'var err error' from here to avoid the unused variable error

		switch platform {
		case "gh":
			token := os.Getenv("GITHUB_TOKEN")
			// Hardcoded for hackathon speed
			sender = postmaster.NewGitHubSender(token, "YOUR_USERNAME", "YOUR_REPO")

		case "discord":
			token := os.Getenv("DISCORD_TOKEN")
			// Handle error immediately inside the case
			ds, err := postmaster.NewDiscordSender(token)
			if err != nil {
				fmt.Printf("error creating discord client: %v\n", err)
				return
			}
			sender = ds

		default:
			fmt.Println("unknown platform. use 'gh' or 'discord'")
			return
		}

		// Execute the send
		fmt.Printf("sending to %s %s...\n", platform, id)
		if err := sender.Send(id, message); err != nil {
			fmt.Printf("failed: %v\n", err)
		} else {
			fmt.Println("success: message sent")
		}
	},
}

func init() {
	rootCmd.AddCommand(replyCmd)
}