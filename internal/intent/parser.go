package intent

import "strings"

type engine struct{}

// NewIntentEngine creates a new deterministic intent engine.
func NewIntentEngine() IntentEngine {
	return &engine{}
}

func (e *engine) Parse(query string) (*Intent, error) {
	lower := strings.ToLower(strings.TrimSpace(query))

	// Basic static keyword matching for V1
	if strings.Contains(lower, "restart") || strings.Contains(lower, "crash") {
		return &Intent{
			Type:     IntentPodRestarts,
			RawQuery: query,
			Filters: IntentFilters{
				ThresholdInt: 1, // Any restarts
			},
		}, nil
	}

	if strings.Contains(lower, "unhealthy") || strings.Contains(lower, "deploy") {
		return &Intent{
			Type:     IntentDeploymentUnhealthy,
			RawQuery: query,
		}, nil
	}

	return &Intent{Type: IntentUnknown, RawQuery: query}, nil
}
