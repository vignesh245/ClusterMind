package context

import (
	"context"
	"fmt"

	"github.com/vignesh245/ClusterMind/internal/kube"
	"github.com/vignesh245/ClusterMind/internal/model"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EvidenceBuilder aggregates K8s resource state into an EvidencePackage.
type EvidenceBuilder struct {
	client kube.Client
}

// NewEvidenceBuilder creates a new EvidenceBuilder.
func NewEvidenceBuilder(client kube.Client) *EvidenceBuilder {
	return &EvidenceBuilder{
		client: client,
	}
}

// BuildForPod constructs an EvidencePackage for a specific Pod.
func (b *EvidenceBuilder) BuildForPod(ctx context.Context, namespace, name string) (*model.EvidencePackage, error) {
	// 1. Fetch Pod
	pods, err := b.client.ListPods(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	var targetPod *corev1.Pod
	for i := range pods {
		if pods[i].Name == name {
			targetPod = &pods[i]
			break
		}
	}
	if targetPod == nil {
		return nil, fmt.Errorf("pod %s not found in namespace %s", name, namespace)
	}

	// 2. Fetch Events
	events, err := b.client.ListEvents(ctx, namespace)
	if err != nil {
		// Log warning, continue
		events = []corev1.Event{}
	}

	// Filter events for this pod
	var recentEvents []model.EventSummary
	for _, e := range events {
		if e.InvolvedObject.Name == targetPod.Name && e.InvolvedObject.Kind == "Pod" {
			recentEvents = append(recentEvents, model.EventSummary{
				Reason:    e.Reason,
				Message:   e.Message,
				Count:     e.Count,
				FirstSeen: e.FirstTimestamp.Time,
				LastSeen:  e.LastTimestamp.Time,
				Type:      e.Type,
			})
		}
	}

	// 3. Status Conditions
	var conditions []model.Condition
	for _, c := range targetPod.Status.Conditions {
		conditions = append(conditions, model.Condition{
			Type:               string(c.Type),
			Status:             string(c.Status),
			Reason:             c.Reason,
			Message:            c.Message,
			LastTransitionTime: c.LastTransitionTime.Time,
		})
	}

	// 4. Restart History
	var restarts []model.ContainerRestart
	for _, cs := range targetPod.Status.ContainerStatuses {
		if cs.RestartCount > 0 {
			var reason string
			var exitCode int32
			var finishedAt metav1.Time
			if cs.LastTerminationState.Terminated != nil {
				reason = cs.LastTerminationState.Terminated.Reason
				exitCode = cs.LastTerminationState.Terminated.ExitCode
				finishedAt = cs.LastTerminationState.Terminated.FinishedAt
			}
			restarts = append(restarts, model.ContainerRestart{
				Container:  cs.Name,
				Count:      cs.RestartCount,
				Reason:     reason,
				ExitCode:   exitCode,
				FinishedAt: finishedAt.Time,
			})
		}
	}

	// 5. Fetch Logs (try failing containers first, or just the first container)
	logExcerpt := ""
	if len(targetPod.Spec.Containers) > 0 {
		containerName := targetPod.Spec.Containers[0].Name
		logs, err := b.client.GetPodLogs(ctx, namespace, name, containerName, 50)
		if err == nil {
			logExcerpt = logs
		}
	}

	pkg := &model.EvidencePackage{
		ResourceKind:     "Pod",
		ResourceName:     name,
		Namespace:        namespace,
		StatusConditions: conditions,
		RecentEvents:     recentEvents,
		LogExcerpt:       logExcerpt,
		RestartHistory:   restarts,
		// ProbeConfig, OwnerChain, Metrics, AnalyzerFindings omitted for brevity in skeleton
	}

	return pkg, nil
}
