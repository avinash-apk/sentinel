package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/avinash-apk/sentinel/pkg/bus"
	"github.com/avinash-apk/sentinel/pkg/ingest"
	"github.com/avinash-apk/sentinel/pkg/postmaster" // Import Postmaster
	"github.com/avinash-apk/sentinel/pkg/tui"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the Sentinel Command Center",
	Run: func(cmd *cobra.Command, args []string) {
		
		token := os.Getenv("DISCORD_TOKEN")
		if token == "" {
			fmt.Println("Error: DISCORD_TOKEN is missing")
			return
		}

		// 1. Setup Bus
		sentinelBus := bus.NewEventBus()
		uiChan := make(chan bus.Event)
		sentinelBus.Subscribe("discord:message", uiChan)

		// 2. Setup Postmaster (Sender)
		// We need this so the TUI can actually send messages
		discordSender, err := postmaster.NewDiscordSender(token)
		if err != nil {
			panic(err)
		}

		// 3. Setup Ingestor (Listener)
		discordListener, err := ingest.NewDiscordIngestor(token, sentinelBus)
		if err != nil {
			panic(err)
		}
		// Start listener
		err = discordListener.Start()
		if err != nil {
			panic(err)
		}

		// 4. Start TUI
		// We pass the sender into the TUI so it can "Reply Straight Up"
		p := tea.NewProgram(tui.InitialModel(uiChan, discordSender))
		if _, err := p.Run(); err != nil {
			fmt.Println("Error starting TUI:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}