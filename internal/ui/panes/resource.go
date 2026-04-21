package panes

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vignesh245/ClusterMind/internal/kube"
	"github.com/vignesh245/ClusterMind/internal/model"
)

type PodsMsg []model.Resource
type ErrMsg error

// ResourcePane handles the display of the Kubernetes resources list.
type ResourcePane struct {
	client    kube.Client
	resources []model.Resource
	active    bool
}

// NewResourcePane creates a new ResourcePane.
func NewResourcePane(client kube.Client) *ResourcePane {
	return &ResourcePane{
		client:    client,
		resources: []model.Resource{},
		active:    true,
	}
}

// Init initializes the pane.
func (p *ResourcePane) Init() tea.Cmd {
	return func() tea.Msg {
		pods, err := p.client.ListPods(context.Background(), "")
		if err != nil {
			return ErrMsg(err)
		}
		var resources []model.Resource
		for _, pod := range pods {
			resources = append(resources, model.Resource{
				Kind:      "Pod",
				Name:      pod.Name,
				Namespace: pod.Namespace,
			})
		}
		return PodsMsg(resources)
	}
}

// Update handles messages.
func (p *ResourcePane) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case PodsMsg:
		p.resources = msg
		return p, nil
	case IntentResultsMsg:
		if msg.Err == nil {
			p.resources = msg.Resources
		}
		return p, nil
	case ErrMsg:
		// For now just ignore error handling explicitly
		return p, nil
	}
	return p, nil
}

// View renders the pane.
func (p *ResourcePane) View() string {
	if !p.active {
		return "Resource Pane (Inactive)"
	}
	if len(p.resources) == 0 {
		return "Resource Pane (Active)\n[Loading resources or none found...]"
	}
	var b strings.Builder
	b.WriteString("Pods:\n")
	for _, res := range p.resources {
		b.WriteString(fmt.Sprintf("  - [%s] %s\n", res.Namespace, res.Name))
	}
	return b.String()
}
