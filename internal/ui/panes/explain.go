package panes

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vignesh245/ClusterMind/internal/ai"
	"github.com/vignesh245/ClusterMind/internal/model"
)

// ExplainPane handles the display of AI-generated Root Cause Analysis.
type ExplainPane struct {
	orchestrator ai.Orchestrator
	result       *model.ExplainResult
	loading      bool
	err          error
	active       bool
}

// NewExplainPane creates a new ExplainPane.
func NewExplainPane(orchestrator ai.Orchestrator) *ExplainPane {
	return &ExplainPane{
		orchestrator: orchestrator,
		active:       false,
	}
}

type ExplainStartMsg struct{}
type ExplainResultMsg struct {
	Result *model.ExplainResult
	Err    error
}

// StartExplain triggers the AI explanation for a given evidence package.
func (p *ExplainPane) StartExplain(pkg *model.EvidencePackage) tea.Cmd {
	p.loading = true
	p.err = nil
	p.result = nil
	
	return func() tea.Msg {
		// In a real application, context should be passed properly.
		// For skeleton, using background context.
		
		res, err := p.orchestrator.Explain(context.Background(), pkg)
		return ExplainResultMsg{
			Result: res,
			Err:    err,
		}
	}
}

func (p *ExplainPane) Init() tea.Cmd {
	return nil
}

func (p *ExplainPane) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ExplainResultMsg:
		p.loading = false
		p.result = msg.Result
		p.err = msg.Err
		return p, nil
	}
	return p, nil
}

func (p *ExplainPane) View() string {
	if !p.active {
		return "Explain Pane (Inactive)"
	}

	if p.loading {
		return "Explain Pane (Active)\n[Generating explanation...]"
	}

	if p.err != nil {
		return fmt.Sprintf("Explain Pane (Active)\n[Error: %v]", p.err)
	}

	if p.result == nil {
		return "Explain Pane (Active)\n[Select a resource and press Shift+E to explain]"
	}

	var b strings.Builder
	b.WriteString("Root Cause Analysis\n")
	b.WriteString(fmt.Sprintf("Confidence: %s\n\n", p.result.Confidence))
	
	b.WriteString("Summary:\n")
	b.WriteString(p.result.Summary + "\n\n")

	b.WriteString("Likely Root Cause:\n")
	b.WriteString(p.result.LikelyRootCause + "\n\n")

	if len(p.result.Evidence) > 0 {
		b.WriteString("Evidence:\n")
		for _, e := range p.result.Evidence {
			b.WriteString(fmt.Sprintf("- %s: %s\n", e.Ref, e.Description))
		}
		b.WriteString("\n")
	}

	if len(p.result.RecommendedActions) > 0 {
		b.WriteString("Recommended Actions:\n")
		for _, a := range p.result.RecommendedActions {
			cmdStr := ""
			if a.Command != nil {
				cmdStr = " (Cmd: " + *a.Command + ")"
			}
			b.WriteString(fmt.Sprintf("- %s%s\n", a.Description, cmdStr))
		}
	}

	return b.String()
}
