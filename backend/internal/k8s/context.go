package k8s

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// ContextManager manages Kubernetes contexts
type ContextManager struct {
	kubeconfig string
	config     *api.Config
}

// NewContextManager creates a new context manager
func NewContextManager(kubeconfig string) (*ContextManager, error) {
	config, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	return &ContextManager{
		kubeconfig: kubeconfig,
		config:     config,
	}, nil
}

// ListContexts returns all available contexts
func (cm *ContextManager) ListContexts() []string {
	contexts := make([]string, 0, len(cm.config.Contexts))
	for name := range cm.config.Contexts {
		contexts = append(contexts, name)
	}
	return contexts
}

// GetCurrentContext returns the current context name
func (cm *ContextManager) GetCurrentContext() string {
	return cm.config.CurrentContext
}

// SwitchContext switches to a different context and returns a new client
func (cm *ContextManager) SwitchContext(contextName string) (*Client, error) {
	if _, exists := cm.config.Contexts[contextName]; !exists {
		return nil, fmt.Errorf("context %s not found", contextName)
	}

	// Build config for the specific context
	configOverrides := &clientcmd.ConfigOverrides{
		CurrentContext: contextName,
	}

	clientConfig := clientcmd.NewNonInteractiveClientConfig(
		*cm.config,
		contextName,
		configOverrides,
		nil,
	)

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create client config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return &Client{
		Clientset: clientset,
		Config:    restConfig, // *rest.Config
	}, nil
}

// NewClientFromContext creates a new client for a specific context
func NewClientFromContext(contextName string) (*Client, error) {
	// Get kubeconfig path
	kubeconfig, err := GetKubeconfigPath()
	if err != nil {
		return nil, err
	}

	cm, err := NewContextManager(kubeconfig)
	if err != nil {
		return nil, err
	}

	return cm.SwitchContext(contextName)
}

// NewClientForCurrentContext creates a client for the current context
func NewClientForCurrentContext() (*Client, error) {
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

// ContextInfo holds information about a Kubernetes context
type ContextInfo struct {
	Name      string `json:"name"`
	Cluster   string `json:"cluster"`
	User      string `json:"user"`
	Namespace string `json:"namespace"`
	IsCurrent bool   `json:"is_current"`
}

// GetContextsInfo returns detailed information about all contexts
func (cm *ContextManager) GetContextsInfo() []ContextInfo {
	contexts := make([]ContextInfo, 0, len(cm.config.Contexts))
	currentContext := cm.config.CurrentContext

	for name, context := range cm.config.Contexts {
		info := ContextInfo{
			Name:      name,
			Cluster:   context.Cluster,
			User:      context.AuthInfo,
			Namespace: context.Namespace,
			IsCurrent: name == currentContext,
		}
		contexts = append(contexts, info)
	}

	return contexts
}
