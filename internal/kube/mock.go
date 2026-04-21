package kube

import (
	"context"

	"github.com/vignesh245/ClusterMind/internal/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type mockClient struct{}

// NewMockClient creates a dummy client for UI testing when no cluster is available.
func NewMockClient() Client {
	return &mockClient{}
}

func (m *mockClient) ListPods(ctx context.Context, namespace string) ([]corev1.Pod, error) {
	return []corev1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "demo-frontend-pod", Namespace: "default"},
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
				ContainerStatuses: []corev1.ContainerStatus{
					{
						Name:         "nginx",
						RestartCount: 5,
						State: corev1.ContainerState{
							Waiting: &corev1.ContainerStateWaiting{Reason: "CrashLoopBackOff"},
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "demo-backend-pod", Namespace: "default"},
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
				ContainerStatuses: []corev1.ContainerStatus{
					{Name: "app", RestartCount: 0},
				},
			},
		},
	}, nil
}

func (m *mockClient) ListDeployments(ctx context.Context, namespace string) ([]appsv1.Deployment, error) {
	replicas := int32(3)
	return []appsv1.Deployment{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "demo-frontend", Namespace: "default"},
			Spec: appsv1.DeploymentSpec{
				Replicas: &replicas,
			},
			Status: appsv1.DeploymentStatus{
				Replicas:          3,
				AvailableReplicas: 1,
			},
		},
	}, nil
}

func (m *mockClient) ListReplicaSets(ctx context.Context, namespace string) ([]appsv1.ReplicaSet, error) {
	return nil, nil
}
func (m *mockClient) ListEvents(ctx context.Context, namespace string) ([]corev1.Event, error) {
	return []corev1.Event{
		{
			Reason:  "BackOff",
			Message: "Back-off restarting failed container",
			Type:    "Warning",
			InvolvedObject: corev1.ObjectReference{
				Kind: "Pod",
				Name: "demo-frontend-pod",
			},
		},
	}, nil
}
func (m *mockClient) ListNodes(ctx context.Context) ([]corev1.Node, error) { return nil, nil }
func (m *mockClient) ListNamespaces(ctx context.Context) ([]corev1.Namespace, error) {
	return nil, nil
}
func (m *mockClient) GetPodLogs(ctx context.Context, namespace, name, container string, lines int) (string, error) {
	if name == "demo-frontend-pod" {
		return "panic: connection refused to backend\ngoroutine 1 [running]:\nmain.main()\n", nil
	}
	return "Listening on :8080...", nil
}
func (m *mockClient) GetMetrics(ctx context.Context, namespace string) (*model.MetricsSummary, error) {
	return nil, nil
}
func (m *mockClient) Watch(ctx context.Context, namespace string) (<-chan WatchEvent, error) {
	return make(chan WatchEvent), nil
}
func (m *mockClient) ApplyPatch(ctx context.Context, resource model.Resource, patch []byte, patchType types.PatchType) error {
	return nil
}
