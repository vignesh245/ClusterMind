package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vignesh245/ClusterMind/internal/ai"
	"github.com/vignesh245/ClusterMind/internal/ai/providers"
	"github.com/vignesh245/ClusterMind/internal/kube"
	"github.com/vignesh245/ClusterMind/internal/remediation"
	"github.com/vignesh245/ClusterMind/internal/ui"
)

func main() {
	client, err := kube.NewClient()
	if err != nil {
		fmt.Printf("Error initializing kubernetes client: %v\n", err)
		os.Exit(1)
	}

	ollamaProvider := providers.NewOllamaProvider("", "llama3.2")
	orchestrator := ai.NewOrchestrator(ollamaProvider)
	exec := remediation.NewExecutor(client)

	app := ui.NewApp(client, orchestrator, exec)
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
