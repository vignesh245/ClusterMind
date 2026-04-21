package intent

import (
	"context"

	"github.com/vignesh245/ClusterMind/internal/kube"
	"github.com/vignesh245/ClusterMind/internal/model"
)

type IntentType string

const (
	IntentPodRestarts         IntentType = "pod_restarts"
	IntentDeploymentUnhealthy IntentType = "deployment_unhealthy"
	IntentUnknown             IntentType = "unknown"
)

type IntentFilters struct {
	Namespace     string
	LabelSelector string
	ThresholdInt  int
}

type Intent struct {
	Type     IntentType
	Filters  IntentFilters
	RawQuery string
}

type IntentEngine interface {
	Parse(query string) (*Intent, error)
	Execute(ctx context.Context, intent *Intent, client kube.Client) ([]model.Resource, error)
}
