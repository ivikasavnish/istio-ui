package kube

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	
	istioclient "istio.io/client-go/pkg/clientset/versioned"
)

type Client struct {
	Clientset     *kubernetes.Clientset
	IstioClient   *istioclient.Clientset
	Config        *rest.Config
}

// NewClient creates a new Kubernetes client
func NewClient() (*Client, error) {
	config, err := getKubeConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	istioClient, err := istioclient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create istio client: %w", err)
	}

	return &Client{
		Clientset:   clientset,
		IstioClient: istioClient,
		Config:      config,
	}, nil
}

// getKubeConfig returns the kubeconfig
func getKubeConfig() (*rest.Config, error) {
	// Try in-cluster config first
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	// Fall back to kubeconfig file
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build config from flags: %w", err)
	}

	return config, nil
}

// ListNamespaces lists all namespaces
func (c *Client) ListNamespaces(ctx context.Context) ([]string, error) {
	namespaces, err := c.Clientset.CoreV1().Namespaces().List(ctx, 
		*new(rest.ListOptions))
	if err != nil {
		return nil, err
	}

	var names []string
	for _, ns := range namespaces.Items {
		names = append(names, ns.Name)
	}
	return names, nil
}

// DiscoverServices discovers services in a namespace
func (c *Client) DiscoverServices(ctx context.Context, namespace string) ([]ServiceInfo, error) {
	services, err := c.Clientset.CoreV1().Services(namespace).List(ctx, 
		*new(rest.ListOptions))
	if err != nil {
		return nil, err
	}

	var infos []ServiceInfo
	for _, svc := range services.Items {
		info := ServiceInfo{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			Labels:    svc.Labels,
			Ports:     []int{},
		}
		
		for _, port := range svc.Spec.Ports {
			info.Ports = append(info.Ports, int(port.Port))
		}
		
		infos = append(infos, info)
	}
	
	return infos, nil
}

// ServiceInfo represents service information
type ServiceInfo struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
	Ports     []int             `json:"ports"`
}

// GetPods gets pods in a namespace with optional label selector
func (c *Client) GetPods(ctx context.Context, namespace string, labelSelector string) ([]PodInfo, error) {
	listOpts := rest.ListOptions{}
	if labelSelector != "" {
		listOpts.LabelSelector = labelSelector
	}
	
	pods, err := c.Clientset.CoreV1().Pods(namespace).List(ctx, listOpts)
	if err != nil {
		return nil, err
	}

	var infos []PodInfo
	for _, pod := range pods.Items {
		info := PodInfo{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Labels:    pod.Labels,
			Phase:     string(pod.Status.Phase),
			Ready:     isPodReady(&pod),
		}
		infos = append(infos, info)
	}
	
	return infos, nil
}

// PodInfo represents pod information
type PodInfo struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
	Phase     string            `json:"phase"`
	Ready     bool              `json:"ready"`
}

func isPodReady(pod *rest.Pod) bool {
	for _, cond := range pod.Status.Conditions {
		if cond.Type == rest.PodReady {
			return cond.Status == rest.ConditionTrue
		}
	}
	return false
}

// GVR returns GroupVersionResource for a given resource
func GVR(group, version, resource string) schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}
}
