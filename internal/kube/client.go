package kube

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/vignesh245/ClusterMind/internal/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// WatchEvent is emitted when a resource changes in the cluster.
type WatchEvent struct {
	Type   watch.EventType
	Object runtime.Object
}

// Client defines the interface for interacting with the Kubernetes cluster.
type Client interface {
	// Read operations
	ListPods(ctx context.Context, namespace string) ([]corev1.Pod, error)
	ListDeployments(ctx context.Context, namespace string) ([]appsv1.Deployment, error)
	ListReplicaSets(ctx context.Context, namespace string) ([]appsv1.ReplicaSet, error)
	ListEvents(ctx context.Context, namespace string) ([]corev1.Event, error)
	ListNodes(ctx context.Context) ([]corev1.Node, error)
	ListNamespaces(ctx context.Context) ([]corev1.Namespace, error)

	GetPodLogs(ctx context.Context, namespace, name, container string, lines int) (string, error)
	GetMetrics(ctx context.Context, namespace string) (*model.MetricsSummary, error)

	// Watch
	Watch(ctx context.Context, namespace string) (<-chan WatchEvent, error)

	// Write (Remediation)
	ApplyPatch(ctx context.Context, resource model.Resource, patch []byte, patchType types.PatchType) error
}

type client struct {
	clientset        *kubernetes.Clientset
	metricsAvailable bool
}

// NewClient creates a new Kubernetes client using the local kubeconfig or in-cluster config.
func NewClient() (Client, error) {
	var config *rest.Config
	var err error

	// Try in-cluster first
	config, err = rest.InClusterConfig()
	if err != nil {
		// Fallback to kubeconfig
		kubeconfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			// Instead of failing entirely, return a mock client for local UI testing
			return NewMockClient(), nil
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return &client{
		clientset:        clientset,
		metricsAvailable: false,
	}, nil
}

func (c *client) ListPods(ctx context.Context, namespace string) ([]corev1.Pod, error) {
	list, err := c.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (c *client) ListDeployments(ctx context.Context, namespace string) ([]appsv1.Deployment, error) {
	list, err := c.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (c *client) ListReplicaSets(ctx context.Context, namespace string) ([]appsv1.ReplicaSet, error) {
	list, err := c.clientset.AppsV1().ReplicaSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (c *client) ListEvents(ctx context.Context, namespace string) ([]corev1.Event, error) {
	list, err := c.clientset.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (c *client) ListNodes(ctx context.Context) ([]corev1.Node, error) {
	list, err := c.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (c *client) ListNamespaces(ctx context.Context) ([]corev1.Namespace, error) {
	list, err := c.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (c *client) GetPodLogs(ctx context.Context, namespace, name, container string, lines int) (string, error) {
	tailLines := int64(lines)
	req := c.clientset.CoreV1().Pods(namespace).GetLogs(name, &corev1.PodLogOptions{
		Container: container,
		TailLines: &tailLines,
	})
	
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return "", err
	}
	defer podLogs.Close()

	buf := new(strings.Builder)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *client) GetMetrics(ctx context.Context, namespace string) (*model.MetricsSummary, error) {
	if !c.metricsAvailable {
		return nil, nil
	}
	// TODO: implement metrics fetching using metrics API clientset
	return nil, nil
}

func (c *client) Watch(ctx context.Context, namespace string) (<-chan WatchEvent, error) {
	ch := make(chan WatchEvent)
	// TODO: setup informers or direct watchers for Pods, Deployments, Events
	// For now, this is a placeholder to satisfy the interface.
	return ch, nil
}

func (c *client) ApplyPatch(ctx context.Context, resource model.Resource, patch []byte, patchType types.PatchType) error {
	// TODO: use dynamic client or clientset to apply the patch based on resource.Kind
	return fmt.Errorf("ApplyPatch not fully implemented in skeleton")
}
