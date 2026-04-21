package diagnostics

import (
	"context"

	"github.com/vignesh245/ClusterMind/internal/model"
	corev1 "k8s.io/api/core/v1"
)

// Analyzer defines the interface for resource diagnostics.
type Analyzer interface {
	Name() string
	AnalyzePod(ctx context.Context, pod corev1.Pod, logs string, events []corev1.Event) ([]model.Finding, error)
}
