package intent

import (
	"context"

	"github.com/vignesh245/ClusterMind/internal/kube"
	"github.com/vignesh245/ClusterMind/internal/model"
)

func (e *engine) Execute(ctx context.Context, intent *Intent, client kube.Client) ([]model.Resource, error) {
	var results []model.Resource

	switch intent.Type {
	case IntentPodRestarts:
		pods, err := client.ListPods(ctx, intent.Filters.Namespace)
		if err != nil {
			return nil, err
		}
		for _, pod := range pods {
			restarts := int32(0)
			for _, cs := range pod.Status.ContainerStatuses {
				restarts += cs.RestartCount
			}
			if restarts >= int32(intent.Filters.ThresholdInt) {
				results = append(results, model.Resource{
					Kind:      "Pod",
					Name:      pod.Name,
					Namespace: pod.Namespace,
				})
			}
		}

	case IntentDeploymentUnhealthy:
		deployments, err := client.ListDeployments(ctx, intent.Filters.Namespace)
		if err != nil {
			return nil, err
		}
		for _, deploy := range deployments {
			if deploy.Status.AvailableReplicas < *deploy.Spec.Replicas {
				results = append(results, model.Resource{
					Kind:      "Deployment",
					Name:      deploy.Name,
					Namespace: deploy.Namespace,
				})
			}
		}
	}

	return results, nil
}
