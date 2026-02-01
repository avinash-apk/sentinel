package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/avinash-apk/sentinel/pkg/bus"
	"github.com/avinash-apk/sentinel/pkg/ingest"
	"github.com/avinash-apk/sentinel/pkg/postmaster"
	"github.com/avinash-apk/sentinel/pkg/tui"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the Sentinel Command Center",
	Run: func(cmd *cobra.Command, args []string) {

		// --- LOAD SECRETS ---
		discordToken := os.Getenv("DISCORD_TOKEN")
		slackAppToken := os.Getenv("SLACK_APP_TOKEN")
		slackBotToken := os.Getenv("SLACK_BOT_TOKEN")

		if discordToken == "" {
			fmt.Println("Warning: DISCORD_TOKEN is missing (Discord features disabled)")
		}

		// --- SETUP BUS ---
		sentinelBus := bus.NewEventBus()
		uiChan := make(chan bus.Event)
		sentinelBus.Subscribe("discord:message", uiChan)
		sentinelBus.Subscribe("slack:message", uiChan)

		// --- SETUP SERVICES ---
		var discordSender *postmaster.DiscordSender
		var slackSender *postmaster.SlackSender

		// 1. Initialize Discord
		if discordToken != "" {
			// Sender
			ds, err := postmaster.NewDiscordSender(discordToken)
			if err != nil {
				fmt.Printf("Error starting Discord Sender: %v\n", err)
			} else {
				discordSender = ds
			}

			// Listener
			dl, err := ingest.NewDiscordIngestor(discordToken, sentinelBus)
			if err != nil {
				fmt.Printf("Error starting Discord Listener: %v\n", err)
			} else {
				if err := dl.Start(); err != nil {
					fmt.Printf("Failed to connect to Discord Gateway: %v\n", err)
				}
			}
		}

		// 2. Initialize Slack
		if slackAppToken != "" && slackBotToken != "" {
			// Sender
			slackSender = postmaster.NewSlackSender(slackBotToken)

			// Listener
			sl := ingest.NewSlackIngestor(slackAppToken, slackBotToken, sentinelBus)
			go sl.Start()
		} else {
			fmt.Println("Warning: SLACK_APP_TOKEN or SLACK_BOT_TOKEN missing (Slack features disabled)")
		}

		// --- START TUI ---
		p := tea.NewProgram(tui.InitialModel(uiChan, discordSender, slackSender))
		if _, err := p.Run(); err != nil {
			fmt.Println("Error starting TUI:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}