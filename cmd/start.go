package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/avinash-apk/sentinel/pkg/actions" // replace with your module
	"github.com/avinash-apk/sentinel/pkg/bus"     // replace with your module
	"github.com/avinash-apk/sentinel/pkg/engine"  // replace with your module
	"github.com/avinash-apk/sentinel/pkg/ingest"  // replace with your module
	"github.com/avinash-apk/sentinel/pkg/tui"     // replace with your module
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the Sentinel event bus",
	Run: func(cmd *cobra.Command, args []string) {
		
		// 1. setup bus
		sentinelBus := bus.NewEventBus()

		// 2. setup ui channel
		// the ui needs its own subscription to display events
		uiChan := make(chan bus.Event)
		sentinelBus.Subscribe("github:event", uiChan)

		// 3. setup engine (logic)
		myRules := []engine.Rule{
			{
				Topic: "github:event",
				Action: &actions.ShellAction{
					Command: "echo 'processed event'",
				},
			},
		}
		eng := &engine.Engine{
			Bus:   sentinelBus,
			Rules: myRules,
		}
		eng.Start()

		// 4. setup ingestor (inputs)
		gh := &ingest.GitHubIngestor{Bus: sentinelBus}
		go gh.Start(":8080")

		// 5. start the tui
		// this blocks until the user quits
		p := tea.NewProgram(tui.InitialModel(uiChan))
		if _, err := p.Run(); err != nil {
			fmt.Printf("error starting tui: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}