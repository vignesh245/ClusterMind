package builtin

import (
	"context"
	"fmt"
	"strings"

	"github.com/vignesh245/ClusterMind/internal/model"
	corev1 "k8s.io/api/core/v1"
)

// PodAnalyzer implements the diagnostics.Analyzer interface for Pods.
type PodAnalyzer struct{}

func (a *PodAnalyzer) Name() string {
	return "builtin-pod-analyzer"
}

func (a *PodAnalyzer) AnalyzePod(ctx context.Context, pod corev1.Pod, logs string, events []corev1.Event) ([]model.Finding, error) {
	var findings []model.Finding

	// Check CrashLoopBackOff
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.State.Waiting != nil && cs.State.Waiting.Reason == "CrashLoopBackOff" {
			findings = append(findings, model.Finding{
				Severity:       model.SeverityCritical,
				Category:       model.CategoryRuntime,
				Title:          "Container CrashLoopBackOff",
				Detail:         "Container " + cs.Name + " is crashlooping. It has restarted " + fmt.Sprintf("%d", cs.RestartCount) + " times.",
				SourceAnalyzer: a.Name(),
			})
		}
	}

	// Check OOMKilled
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.LastTerminationState.Terminated != nil && cs.LastTerminationState.Terminated.Reason == "OOMKilled" {
			findings = append(findings, model.Finding{
				Severity:       model.SeverityCritical,
				Category:       model.CategoryResource,
				Title:          "Container OOMKilled",
				Detail:         "Container " + cs.Name + " was terminated because it exceeded its memory limit.",
				SourceAnalyzer: a.Name(),
			})
		}
	}

	// Check ImagePullBackOff
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.State.Waiting != nil && (cs.State.Waiting.Reason == "ImagePullBackOff" || cs.State.Waiting.Reason == "ErrImagePull") {
			findings = append(findings, model.Finding{
				Severity:       model.SeverityWarning,
				Category:       model.CategoryConfig,
				Title:          "Image Pull Failure",
				Detail:         "Failed to pull image for container " + cs.Name + ": " + cs.State.Waiting.Message,
				SourceAnalyzer: a.Name(),
			})
		}
	}

	return findings, nil
}
