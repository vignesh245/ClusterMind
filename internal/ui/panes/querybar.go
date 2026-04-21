package panes

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vignesh245/ClusterMind/internal/intent"
	"github.com/vignesh245/ClusterMind/internal/kube"
	"github.com/vignesh245/ClusterMind/internal/model"
)

// QueryBar handles the natural language intent input.
type QueryBar struct {
	engine intent.IntentEngine
	client kube.Client
	input  string
	active bool
}

// IntentResultsMsg is emitted when an intent query finishes.
type IntentResultsMsg struct {
	Resources []model.Resource
	Intent    *intent.Intent
	Err       error
}

// NewQueryBar creates a new QueryBar.
func NewQueryBar(engine intent.IntentEngine, client kube.Client) *QueryBar {
	return &QueryBar{
		engine: engine,
		client: client,
	}
}

func (q *QueryBar) Init() tea.Cmd {
	return nil
}

func (q *QueryBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !q.active {
			if msg.String() == ":" {
				q.active = true
				q.input = ""
				return q, nil
			}
			return q, nil
		}

		switch msg.String() {
		case "esc":
			q.active = false
			q.input = ""
			return q, nil
		case "enter":
			q.active = false
			return q, q.executeIntent(q.input)
		case "backspace", "delete":
			if len(q.input) > 0 {
				q.input = q.input[:len(q.input)-1]
			}
		default:
			// basic single-char string collection
			if len(msg.String()) == 1 {
				q.input += msg.String()
			}
		}
	}
	return q, nil
}

func (q *QueryBar) executeIntent(query string) tea.Cmd {
	return func() tea.Msg {
		intentObj, err := q.engine.Parse(query)
		if err != nil {
			return IntentResultsMsg{Err: err}
		}

		if intentObj.Type == intent.IntentUnknown {
			return IntentResultsMsg{Err: fmt.Errorf("unrecognized query intent")}
		}

		resources, err := q.engine.Execute(context.Background(), intentObj, q.client)
		return IntentResultsMsg{
			Resources: resources,
			Intent:    intentObj,
			Err:       err,
		}
	}
}

func (q *QueryBar) View() string {
	if q.active {
		return ": " + q.input + "█"
	}
	return "[:] Query Bar"
}
