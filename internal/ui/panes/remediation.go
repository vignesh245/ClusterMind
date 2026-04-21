package panes

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vignesh245/ClusterMind/internal/ai"
	"github.com/vignesh245/ClusterMind/internal/model"
	"github.com/vignesh245/ClusterMind/internal/remediation"
)

// RemediationPrompt displays an AI remediation suggestion and prompts for confirmation.
type RemediationPrompt struct {
	orchestrator ai.Orchestrator
	executor     remediation.Executor
	plan         *model.RemediationPlan
	loading      bool
	executing    bool
	done         bool
	err          error
	active       bool
}

type RemediationPlanMsg struct {
	Plan *model.RemediationPlan
	Err  error
}

type RemediationExecMsg struct {
	Err error
}

func NewRemediationPrompt(orchestrator ai.Orchestrator, executor remediation.Executor) *RemediationPrompt {
	return &RemediationPrompt{
		orchestrator: orchestrator,
		executor:     executor,
		active:       false,
	}
}

// StartRemediation asks the AI for a remediation plan based on evidence.
func (p *RemediationPrompt) StartRemediation(pkg *model.EvidencePackage) tea.Cmd {
	p.loading = true
	p.err = nil
	p.plan = nil
	p.active = true
	p.done = false

	return func() tea.Msg {
		// Skeleton uses Background Context
		res, err := p.orchestrator.SuggestRemediation(context.Background(), pkg)
		return RemediationPlanMsg{
			Plan: res,
			Err:  err,
		}
	}
}

func (p *RemediationPrompt) Init() tea.Cmd {
	return nil
}

func (p *RemediationPrompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case RemediationPlanMsg:
		p.loading = false
		p.plan = msg.Plan
		p.err = msg.Err
		return p, nil
		
	case RemediationExecMsg:
		p.executing = false
		p.done = true
		p.err = msg.Err
		return p, nil

	case tea.KeyMsg:
		if !p.active || p.plan == nil || p.loading || p.executing || p.done {
			return p, nil
		}
		
		switch msg.String() {
		case "y", "Y":
			p.executing = true
			p.err = nil
			return p, func() tea.Msg {
				err := p.executor.Execute(context.Background(), p.plan)
				return RemediationExecMsg{Err: err}
			}
		case "n", "N", "esc":
			p.active = false
			p.plan = nil
			return p, nil
		}
	}
	return p, nil
}

func (p *RemediationPrompt) View() string {
	if !p.active {
		return ""
	}

	var b strings.Builder
	b.WriteString("\n--- REMEDIATION PROMPT ---\n")

	if p.loading {
		b.WriteString("[Analyzing evidence for remediation...]\n")
		return b.String()
	}

	if p.err != nil {
		b.WriteString(fmt.Sprintf("[Error: %v]\n", p.err))
		return b.String()
	}
	
	if p.executing {
		b.WriteString("[Executing remediation action...]\n")
		return b.String()
	}

	if p.done {
		b.WriteString("[Remediation executed successfully! Press ESC to dismiss.]\n")
		return b.String()
	}

	if p.plan != nil {
		b.WriteString("Proposed Action:\n")
		b.WriteString(fmt.Sprintf("Risk: %s\n", p.plan.RiskLevel))
		b.WriteString(fmt.Sprintf("Rationale: %s\n\n", p.plan.Rationale))
		
		b.WriteString(fmt.Sprintf("Type: %s\n", p.plan.RemediationType))
		if p.plan.ProposedCommand != "" {
			b.WriteString(fmt.Sprintf("Command: %s\n", p.plan.ProposedCommand))
		}
		if p.plan.ProposedPatch != "" {
			b.WriteString(fmt.Sprintf("Patch: %s\n", p.plan.ProposedPatch))
		}
		
		b.WriteString("\nExecute this action? [y/N]: ")
	}

	return b.String()
}
