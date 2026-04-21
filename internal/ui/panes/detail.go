package panes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vignesh245/ClusterMind/internal/model"
)

// DetailPane handles the display of a selected resource's details.
type DetailPane struct {
	resource *model.Resource
	active   bool
}

// NewDetailPane creates a new DetailPane.
func NewDetailPane() *DetailPane {
	return &DetailPane{
		active: false,
	}
}

// SetResource updates the currently displayed resource.
func (p *DetailPane) SetResource(res *model.Resource) {
	p.resource = res
}

// Init initializes the pane.
func (p *DetailPane) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (p *DetailPane) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return p, nil
}

// View renders the pane.
func (p *DetailPane) View() string {
	if p.resource == nil {
		return "Detail Pane\n[No resource selected]"
	}
	return "Detail Pane\nName: " + p.resource.Name
}
