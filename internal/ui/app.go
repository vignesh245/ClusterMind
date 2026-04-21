package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vignesh245/ClusterMind/internal/ai"
	"github.com/vignesh245/ClusterMind/internal/intent"
	"github.com/vignesh245/ClusterMind/internal/kube"
	"github.com/vignesh245/ClusterMind/internal/remediation"
	"github.com/vignesh245/ClusterMind/internal/ui/panes"
)

var (
	baseStyle = lipgloss.NewStyle().Margin(1, 2)

	paneStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2).
		Width(45).
		Height(15)

	queryBarStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Width(94) // matches roughly 45*2 + borders
)

// App is the root Bubble Tea model for ClusterMind.
type App struct {
	client       kube.Client
	orchestrator ai.Orchestrator
	intentEngine intent.IntentEngine
	resourcePane      *panes.ResourcePane
	detailPane        *panes.DetailPane
	explainPane       *panes.ExplainPane
	queryBar          *panes.QueryBar
	remediationPrompt *panes.RemediationPrompt
}

// NewApp creates a new Bubble Tea App model.
func NewApp(client kube.Client, orchestrator ai.Orchestrator, exec remediation.Executor) *App {
	ie := intent.NewIntentEngine()
	return &App{
		client:            client,
		orchestrator:      orchestrator,
		intentEngine:      ie,
		resourcePane:      panes.NewResourcePane(client),
		detailPane:        panes.NewDetailPane(),
		explainPane:       panes.NewExplainPane(orchestrator),
		queryBar:          panes.NewQueryBar(ie, client),
		remediationPrompt: panes.NewRemediationPrompt(orchestrator, exec),
	}
}

// Init is the first function that will be called. It returns an optional
// initial command.
func (a *App) Init() tea.Cmd {
	return tea.Batch(
		a.resourcePane.Init(),
		a.detailPane.Init(),
		a.explainPane.Init(),
		a.queryBar.Init(),
		a.remediationPrompt.Init(),
	)
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		}
	}
	
	// Delegate update to panes
	var cmd tea.Cmd
	_, cmd = a.resourcePane.Update(msg)
	cmds = append(cmds, cmd)
	_, cmd = a.detailPane.Update(msg)
	cmds = append(cmds, cmd)
	_, cmd = a.explainPane.Update(msg)
	cmds = append(cmds, cmd)
	_, cmd = a.queryBar.Update(msg)
	cmds = append(cmds, cmd)
	_, cmd = a.remediationPrompt.Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

// View renders the program's UI, which is just a string.
func (a *App) View() string {
	topRow := lipgloss.JoinHorizontal(lipgloss.Top,
		paneStyle.Render(a.resourcePane.View()),
		paneStyle.Render(a.detailPane.View()),
	)

	middleRow := lipgloss.JoinHorizontal(lipgloss.Top,
		paneStyle.Width(94).Render(a.explainPane.View()),
	)

	bottomRow := queryBarStyle.Render(a.queryBar.View())
	
	// Add Remediation overlay if active
	remediationView := a.remediationPrompt.View()
	if remediationView != "" {
		remediationBox := lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("196")).
			Padding(1, 2).
			Render(remediationView)
		
		middleRow = lipgloss.JoinHorizontal(lipgloss.Top, remediationBox)
	}

	layout := lipgloss.JoinVertical(lipgloss.Left,
		"ClusterMind TUI - Press 'q' to quit | Press ':' to query",
		topRow,
		middleRow,
		bottomRow,
	)

	return baseStyle.Render(layout)
}
