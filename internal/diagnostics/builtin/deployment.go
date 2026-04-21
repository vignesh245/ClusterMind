package builtin

import (
	"context"

	"github.com/vignesh245/ClusterMind/internal/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// DeploymentAnalyzer implements diagnostics for Deployments.
type DeploymentAnalyzer struct{}

func (a *DeploymentAnalyzer) Name() string {
	return "builtin-deployment-analyzer"
}

func (a *DeploymentAnalyzer) AnalyzeDeployment(ctx context.Context, deploy appsv1.Deployment, pods []corev1.Pod, events []corev1.Event) ([]model.Finding, error) {
	var findings []model.Finding

	// Check if deployment is not fully available
	if deploy.Status.AvailableReplicas < deploy.Status.Replicas {
		findings = append(findings, model.Finding{
			Severity:       model.SeverityWarning,
			Category:       model.CategoryScheduling,
			Title:          "Deployment Unhealthy",
			Detail:         "Deployment has unavailable replicas",
			SourceAnalyzer: a.Name(),
		})
	}

	// Stalled rollout
	for _, c := range deploy.Status.Conditions {
		if c.Type == appsv1.DeploymentProgressing && c.Status == corev1.ConditionFalse {
			findings = append(findings, model.Finding{
				Severity:       model.SeverityCritical,
				Category:       model.CategoryRuntime,
				Title:          "Rollout Stalled",
				Detail:         c.Message,
				SourceAnalyzer: a.Name(),
			})
		}
	}

	return findings, nil
}
