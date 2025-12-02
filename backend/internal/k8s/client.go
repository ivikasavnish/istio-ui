package k8s

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client wraps the Kubernetes clientset
type Client struct {
	Clientset *kubernetes.Clientset
	Config    *rest.Config
}

// NewClient creates a new Kubernetes client
func NewClient() (*Client, error) {
	config, err := getConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		Clientset: clientset,
		Config:    config,
	}, nil
}

// getConfig returns Kubernetes config from in-cluster or kubeconfig
func getConfig() (*rest.Config, error) {
	// Try in-cluster config first
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	// Fall back to kubeconfig
	kubeconfig, err := GetKubeconfigPath()
	if err != nil {
		return nil, err
	}

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// GetKubeconfigPath returns the path to the kubeconfig file
func GetKubeconfigPath() (string, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	return kubeconfig, nil
}
