package ingest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/avinash-apk/sentinel/pkg/bus" 
)

type GitHubIngestor struct {
	Bus *bus.EventBus
}

func (g *GitHubIngestor) Start(port string) {
	// handle requests to /webhook
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// decode the incoming json
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// publish to the bus
		// topic is 'github:event', payload is the json data
		g.Bus.Publish("github:event", payload)
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// start the server in a way that doesn't block the main thread immediately
	// but for this simple version, we will run it in a goroutine in start.go
	fmt.Printf("⚡ http server listening on %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}