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
			// Make sure these match your actual GitHub details
			sender = postmaster.NewGitHubSender(token, "YOUR_USERNAME", "sentinel")

		case "discord":
			token := os.Getenv("DISCORD_TOKEN")
			ds, err := postmaster.NewDiscordSender(token)
			if err != nil {
				fmt.Printf("error creating discord client: %v\n", err)
				return
			}
			sender = ds

		case "slack":
			token := os.Getenv("SLACK_TOKEN")
			if token == "" {
				fmt.Println("error: SLACK_TOKEN env var not set")
				return
			}
			sender = postmaster.NewSlackSender(token)

		case "email":
			user := os.Getenv("EMAIL_USER")
			pass := os.Getenv("EMAIL_PASS")
			if user == "" || pass == "" {
				fmt.Println("error: EMAIL_USER or EMAIL_PASS not set")
				return
			}
			sender = postmaster.NewEmailSender(user, pass)

		default:
			fmt.Println("unknown platform. use 'gh', 'discord', 'slack', or 'email'")
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