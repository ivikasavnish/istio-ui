package istio

import (
	"context"
	"fmt"
	
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	
	"github.com/ivikasavnish/istio-ui/backend/internal/kube"
	"github.com/ivikasavnish/istio-ui/backend/internal/models"
)

type Manager struct {
	kubeClient *kube.Client
	dynamicClient dynamic.Interface
}

func NewManager(kubeClient *kube.Client) (*Manager, error) {
	dynamicClient, err := dynamic.NewForConfig(kubeClient.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}
	
	return &Manager{
		kubeClient:    kubeClient,
		dynamicClient: dynamicClient,
	}, nil
}

// VirtualService operations
var virtualServiceGVR = schema.GroupVersionResource{
	Group:    "networking.istio.io",
	Version:  "v1beta1",
	Resource: "virtualservices",
}

func (m *Manager) CreateVirtualService(ctx context.Context, spec *models.VirtualServiceSpec) error {
	obj := m.virtualServiceToUnstructured(spec)
	_, err := m.dynamicClient.Resource(virtualServiceGVR).
		Namespace(spec.Namespace).
		Create(ctx, obj, metav1.CreateOptions{})
	return err
}

func (m *Manager) UpdateVirtualService(ctx context.Context, spec *models.VirtualServiceSpec) error {
	obj := m.virtualServiceToUnstructured(spec)
	_, err := m.dynamicClient.Resource(virtualServiceGVR).
		Namespace(spec.Namespace).
		Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func (m *Manager) GetVirtualService(ctx context.Context, namespace, name string) (*unstructured.Unstructured, error) {
	return m.dynamicClient.Resource(virtualServiceGVR).
		Namespace(namespace).
		Get(ctx, name, metav1.GetOptions{})
}

func (m *Manager) ListVirtualServices(ctx context.Context, namespace string) (*unstructured.UnstructuredList, error) {
	return m.dynamicClient.Resource(virtualServiceGVR).
		Namespace(namespace).
		List(ctx, metav1.ListOptions{})
}

func (m *Manager) DeleteVirtualService(ctx context.Context, namespace, name string) error {
	return m.dynamicClient.Resource(virtualServiceGVR).
		Namespace(namespace).
		Delete(ctx, name, metav1.DeleteOptions{})
}

// DestinationRule operations
var destinationRuleGVR = schema.GroupVersionResource{
	Group:    "networking.istio.io",
	Version:  "v1beta1",
	Resource: "destinationrules",
}

func (m *Manager) CreateDestinationRule(ctx context.Context, spec *models.DestinationRuleSpec) error {
	obj := m.destinationRuleToUnstructured(spec)
	_, err := m.dynamicClient.Resource(destinationRuleGVR).
		Namespace(spec.Namespace).
		Create(ctx, obj, metav1.CreateOptions{})
	return err
}

func (m *Manager) UpdateDestinationRule(ctx context.Context, spec *models.DestinationRuleSpec) error {
	obj := m.destinationRuleToUnstructured(spec)
	_, err := m.dynamicClient.Resource(destinationRuleGVR).
		Namespace(spec.Namespace).
		Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func (m *Manager) GetDestinationRule(ctx context.Context, namespace, name string) (*unstructured.Unstructured, error) {
	return m.dynamicClient.Resource(destinationRuleGVR).
		Namespace(namespace).
		Get(ctx, name, metav1.GetOptions{})
}

func (m *Manager) ListDestinationRules(ctx context.Context, namespace string) (*unstructured.UnstructuredList, error) {
	return m.dynamicClient.Resource(destinationRuleGVR).
		Namespace(namespace).
		List(ctx, metav1.ListOptions{})
}

func (m *Manager) DeleteDestinationRule(ctx context.Context, namespace, name string) error {
	return m.dynamicClient.Resource(destinationRuleGVR).
		Namespace(namespace).
		Delete(ctx, name, metav1.DeleteOptions{})
}

// Gateway operations
var gatewayGVR = schema.GroupVersionResource{
	Group:    "networking.istio.io",
	Version:  "v1beta1",
	Resource: "gateways",
}

func (m *Manager) CreateGateway(ctx context.Context, spec *models.GatewaySpec) error {
	obj := m.gatewayToUnstructured(spec)
	_, err := m.dynamicClient.Resource(gatewayGVR).
		Namespace(spec.Namespace).
		Create(ctx, obj, metav1.CreateOptions{})
	return err
}

func (m *Manager) ListGateways(ctx context.Context, namespace string) (*unstructured.UnstructuredList, error) {
	return m.dynamicClient.Resource(gatewayGVR).
		Namespace(namespace).
		List(ctx, metav1.ListOptions{})
}

func (m *Manager) DeleteGateway(ctx context.Context, namespace, name string) error {
	return m.dynamicClient.Resource(gatewayGVR).
		Namespace(namespace).
		Delete(ctx, name, metav1.DeleteOptions{})
}

// PeerAuthentication operations
var peerAuthGVR = schema.GroupVersionResource{
	Group:    "security.istio.io",
	Version:  "v1beta1",
	Resource: "peerauthentications",
}

func (m *Manager) CreatePeerAuthentication(ctx context.Context, spec *models.PeerAuthenticationSpec) error {
	obj := m.peerAuthToUnstructured(spec)
	_, err := m.dynamicClient.Resource(peerAuthGVR).
		Namespace(spec.Namespace).
		Create(ctx, obj, metav1.CreateOptions{})
	return err
}

func (m *Manager) ListPeerAuthentications(ctx context.Context, namespace string) (*unstructured.UnstructuredList, error) {
	return m.dynamicClient.Resource(peerAuthGVR).
		Namespace(namespace).
		List(ctx, metav1.ListOptions{})
}

func (m *Manager) DeletePeerAuthentication(ctx context.Context, namespace, name string) error {
	return m.dynamicClient.Resource(peerAuthGVR).
		Namespace(namespace).
		Delete(ctx, name, metav1.DeleteOptions{})
}

// AuthorizationPolicy operations
var authzPolicyGVR = schema.GroupVersionResource{
	Group:    "security.istio.io",
	Version:  "v1beta1",
	Resource: "authorizationpolicies",
}

func (m *Manager) CreateAuthorizationPolicy(ctx context.Context, spec *models.AuthorizationPolicySpec) error {
	obj := m.authzPolicyToUnstructured(spec)
	_, err := m.dynamicClient.Resource(authzPolicyGVR).
		Namespace(spec.Namespace).
		Create(ctx, obj, metav1.CreateOptions{})
	return err
}

func (m *Manager) ListAuthorizationPolicies(ctx context.Context, namespace string) (*unstructured.UnstructuredList, error) {
	return m.dynamicClient.Resource(authzPolicyGVR).
		Namespace(namespace).
		List(ctx, metav1.ListOptions{})
}

func (m *Manager) DeleteAuthorizationPolicy(ctx context.Context, namespace, name string) error {
	return m.dynamicClient.Resource(authzPolicyGVR).
		Namespace(namespace).
		Delete(ctx, name, metav1.DeleteOptions{})
}

// Helper methods to convert models to unstructured objects
func (m *Manager) virtualServiceToUnstructured(spec *models.VirtualServiceSpec) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "networking.istio.io/v1beta1",
			"kind":       "VirtualService",
			"metadata": map[string]interface{}{
				"name":      spec.Name,
				"namespace": spec.Namespace,
			},
			"spec": map[string]interface{}{
				"hosts": spec.Hosts,
			},
		},
	}
	
	if len(spec.Gateways) > 0 {
		obj.Object["spec"].(map[string]interface{})["gateways"] = spec.Gateways
	}
	
	if len(spec.HTTP) > 0 {
		obj.Object["spec"].(map[string]interface{})["http"] = spec.HTTP
	}
	
	return obj
}

func (m *Manager) destinationRuleToUnstructured(spec *models.DestinationRuleSpec) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "networking.istio.io/v1beta1",
			"kind":       "DestinationRule",
			"metadata": map[string]interface{}{
				"name":      spec.Name,
				"namespace": spec.Namespace,
			},
			"spec": map[string]interface{}{
				"host": spec.Host,
			},
		},
	}
	
	if spec.TrafficPolicy != nil {
		obj.Object["spec"].(map[string]interface{})["trafficPolicy"] = spec.TrafficPolicy
	}
	
	if len(spec.Subsets) > 0 {
		obj.Object["spec"].(map[string]interface{})["subsets"] = spec.Subsets
	}
	
	return obj
}

func (m *Manager) gatewayToUnstructured(spec *models.GatewaySpec) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "networking.istio.io/v1beta1",
			"kind":       "Gateway",
			"metadata": map[string]interface{}{
				"name":      spec.Name,
				"namespace": spec.Namespace,
			},
			"spec": map[string]interface{}{
				"selector": spec.Selector,
				"servers":  spec.Servers,
			},
		},
	}
}

func (m *Manager) peerAuthToUnstructured(spec *models.PeerAuthenticationSpec) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "security.istio.io/v1beta1",
			"kind":       "PeerAuthentication",
			"metadata": map[string]interface{}{
				"name":      spec.Name,
				"namespace": spec.Namespace,
			},
			"spec": map[string]interface{}{
				"mtls": map[string]interface{}{
					"mode": spec.MtlsMode,
				},
			},
		},
	}
	
	if len(spec.Selector) > 0 {
		obj.Object["spec"].(map[string]interface{})["selector"] = map[string]interface{}{
			"matchLabels": spec.Selector,
		}
	}
	
	return obj
}

func (m *Manager) authzPolicyToUnstructured(spec *models.AuthorizationPolicySpec) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "security.istio.io/v1beta1",
			"kind":       "AuthorizationPolicy",
			"metadata": map[string]interface{}{
				"name":      spec.Name,
				"namespace": spec.Namespace,
			},
			"spec": map[string]interface{}{
				"action": spec.Action,
			},
		},
	}
	
	if len(spec.Selector) > 0 {
		obj.Object["spec"].(map[string]interface{})["selector"] = map[string]interface{}{
			"matchLabels": spec.Selector,
		}
	}
	
	if len(spec.Rules) > 0 {
		obj.Object["spec"].(map[string]interface{})["rules"] = spec.Rules
	}
	
	return obj
}

// GenerateYAML generates YAML from an unstructured object
func GenerateYAML(obj *unstructured.Unstructured) (string, error) {
	data, err := yaml.Marshal(obj.Object)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
