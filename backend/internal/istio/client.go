package istio

import (
	"context"

	"github.com/ivikasavnish/istio-ui/internal/k8s"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	networkingv1beta1 "istio.io/api/networking/v1beta1"
	securityv1beta1 "istio.io/api/security/v1beta1"
	networkingclient "istio.io/client-go/pkg/apis/networking/v1beta1"
	securityclient "istio.io/client-go/pkg/apis/security/v1beta1"
)

// Client wraps the Istio clientset
type Client struct {
	Clientset *versionedclient.Clientset
}

// NewClient creates a new Istio client
func NewClient(k8sClient *k8s.Client) (*Client, error) {
	clientset, err := versionedclient.NewForConfig(k8sClient.Config)
	if err != nil {
		return nil, err
	}

	return &Client{
		Clientset: clientset,
	}, nil
}

// VirtualService operations
func (c *Client) ListVirtualServices(ctx context.Context, namespace string) (*networkingclient.VirtualServiceList, error) {
	return c.Clientset.NetworkingV1beta1().VirtualServices(namespace).List(ctx, metav1.ListOptions{})
}

func (c *Client) GetVirtualService(ctx context.Context, namespace, name string) (*networkingclient.VirtualService, error) {
	return c.Clientset.NetworkingV1beta1().VirtualServices(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (c *Client) CreateVirtualService(ctx context.Context, vs *networkingclient.VirtualService) (*networkingclient.VirtualService, error) {
	return c.Clientset.NetworkingV1beta1().VirtualServices(vs.Namespace).Create(ctx, vs, metav1.CreateOptions{})
}

func (c *Client) UpdateVirtualService(ctx context.Context, vs *networkingclient.VirtualService) (*networkingclient.VirtualService, error) {
	return c.Clientset.NetworkingV1beta1().VirtualServices(vs.Namespace).Update(ctx, vs, metav1.UpdateOptions{})
}

func (c *Client) DeleteVirtualService(ctx context.Context, namespace, name string) error {
	return c.Clientset.NetworkingV1beta1().VirtualServices(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// DestinationRule operations
func (c *Client) ListDestinationRules(ctx context.Context, namespace string) (*networkingclient.DestinationRuleList, error) {
	return c.Clientset.NetworkingV1beta1().DestinationRules(namespace).List(ctx, metav1.ListOptions{})
}

func (c *Client) GetDestinationRule(ctx context.Context, namespace, name string) (*networkingclient.DestinationRule, error) {
	return c.Clientset.NetworkingV1beta1().DestinationRules(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (c *Client) CreateDestinationRule(ctx context.Context, dr *networkingclient.DestinationRule) (*networkingclient.DestinationRule, error) {
	return c.Clientset.NetworkingV1beta1().DestinationRules(dr.Namespace).Create(ctx, dr, metav1.CreateOptions{})
}

func (c *Client) UpdateDestinationRule(ctx context.Context, dr *networkingclient.DestinationRule) (*networkingclient.DestinationRule, error) {
	return c.Clientset.NetworkingV1beta1().DestinationRules(dr.Namespace).Update(ctx, dr, metav1.UpdateOptions{})
}

func (c *Client) DeleteDestinationRule(ctx context.Context, namespace, name string) error {
	return c.Clientset.NetworkingV1beta1().DestinationRules(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// Gateway operations
func (c *Client) ListGateways(ctx context.Context, namespace string) (*networkingclient.GatewayList, error) {
	return c.Clientset.NetworkingV1beta1().Gateways(namespace).List(ctx, metav1.ListOptions{})
}

func (c *Client) GetGateway(ctx context.Context, namespace, name string) (*networkingclient.Gateway, error) {
	return c.Clientset.NetworkingV1beta1().Gateways(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (c *Client) CreateGateway(ctx context.Context, gw *networkingclient.Gateway) (*networkingclient.Gateway, error) {
	return c.Clientset.NetworkingV1beta1().Gateways(gw.Namespace).Create(ctx, gw, metav1.CreateOptions{})
}

func (c *Client) UpdateGateway(ctx context.Context, gw *networkingclient.Gateway) (*networkingclient.Gateway, error) {
	return c.Clientset.NetworkingV1beta1().Gateways(gw.Namespace).Update(ctx, gw, metav1.UpdateOptions{})
}

func (c *Client) DeleteGateway(ctx context.Context, namespace, name string) error {
	return c.Clientset.NetworkingV1beta1().Gateways(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// ServiceEntry operations
func (c *Client) ListServiceEntries(ctx context.Context, namespace string) (*networkingclient.ServiceEntryList, error) {
	return c.Clientset.NetworkingV1beta1().ServiceEntries(namespace).List(ctx, metav1.ListOptions{})
}

func (c *Client) GetServiceEntry(ctx context.Context, namespace, name string) (*networkingclient.ServiceEntry, error) {
	return c.Clientset.NetworkingV1beta1().ServiceEntries(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (c *Client) CreateServiceEntry(ctx context.Context, se *networkingclient.ServiceEntry) (*networkingclient.ServiceEntry, error) {
	return c.Clientset.NetworkingV1beta1().ServiceEntries(se.Namespace).Create(ctx, se, metav1.CreateOptions{})
}

func (c *Client) UpdateServiceEntry(ctx context.Context, se *networkingclient.ServiceEntry) (*networkingclient.ServiceEntry, error) {
	return c.Clientset.NetworkingV1beta1().ServiceEntries(se.Namespace).Update(ctx, se, metav1.UpdateOptions{})
}

func (c *Client) DeleteServiceEntry(ctx context.Context, namespace, name string) error {
	return c.Clientset.NetworkingV1beta1().ServiceEntries(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// AuthorizationPolicy operations
func (c *Client) ListAuthorizationPolicies(ctx context.Context, namespace string) (*securityclient.AuthorizationPolicyList, error) {
	return c.Clientset.SecurityV1beta1().AuthorizationPolicies(namespace).List(ctx, metav1.ListOptions{})
}

func (c *Client) GetAuthorizationPolicy(ctx context.Context, namespace, name string) (*securityclient.AuthorizationPolicy, error) {
	return c.Clientset.SecurityV1beta1().AuthorizationPolicies(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (c *Client) CreateAuthorizationPolicy(ctx context.Context, ap *securityclient.AuthorizationPolicy) (*securityclient.AuthorizationPolicy, error) {
	return c.Clientset.SecurityV1beta1().AuthorizationPolicies(ap.Namespace).Create(ctx, ap, metav1.CreateOptions{})
}

func (c *Client) UpdateAuthorizationPolicy(ctx context.Context, ap *securityclient.AuthorizationPolicy) (*securityclient.AuthorizationPolicy, error) {
	return c.Clientset.SecurityV1beta1().AuthorizationPolicies(ap.Namespace).Update(ctx, ap, metav1.UpdateOptions{})
}

func (c *Client) DeleteAuthorizationPolicy(ctx context.Context, namespace, name string) error {
	return c.Clientset.SecurityV1beta1().AuthorizationPolicies(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// PeerAuthentication operations
func (c *Client) ListPeerAuthentications(ctx context.Context, namespace string) (*securityclient.PeerAuthenticationList, error) {
	return c.Clientset.SecurityV1beta1().PeerAuthentications(namespace).List(ctx, metav1.ListOptions{})
}

func (c *Client) GetPeerAuthentication(ctx context.Context, namespace, name string) (*securityclient.PeerAuthentication, error) {
	return c.Clientset.SecurityV1beta1().PeerAuthentications(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (c *Client) CreatePeerAuthentication(ctx context.Context, pa *securityclient.PeerAuthentication) (*securityclient.PeerAuthentication, error) {
	return c.Clientset.SecurityV1beta1().PeerAuthentications(pa.Namespace).Create(ctx, pa, metav1.CreateOptions{})
}

func (c *Client) UpdatePeerAuthentication(ctx context.Context, pa *securityclient.PeerAuthentication) (*securityclient.PeerAuthentication, error) {
	return c.Clientset.SecurityV1beta1().PeerAuthentications(pa.Namespace).Update(ctx, pa, metav1.UpdateOptions{})
}

func (c *Client) DeletePeerAuthentication(ctx context.Context, namespace, name string) error {
	return c.Clientset.SecurityV1beta1().PeerAuthentications(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// RequestAuthentication operations
func (c *Client) ListRequestAuthentications(ctx context.Context, namespace string) (*securityclient.RequestAuthenticationList, error) {
	return c.Clientset.SecurityV1beta1().RequestAuthentications(namespace).List(ctx, metav1.ListOptions{})
}

func (c *Client) GetRequestAuthentication(ctx context.Context, namespace, name string) (*securityclient.RequestAuthentication, error) {
	return c.Clientset.SecurityV1beta1().RequestAuthentications(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (c *Client) CreateRequestAuthentication(ctx context.Context, ra *securityclient.RequestAuthentication) (*securityclient.RequestAuthentication, error) {
	return c.Clientset.SecurityV1beta1().RequestAuthentications(ra.Namespace).Create(ctx, ra, metav1.CreateOptions{})
}

func (c *Client) UpdateRequestAuthentication(ctx context.Context, ra *securityclient.RequestAuthentication) (*securityclient.RequestAuthentication, error) {
	return c.Clientset.SecurityV1beta1().RequestAuthentications(ra.Namespace).Update(ctx, ra, metav1.UpdateOptions{})
}

func (c *Client) DeleteRequestAuthentication(ctx context.Context, namespace, name string) error {
	return c.Clientset.SecurityV1beta1().RequestAuthentications(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// Helper types for unused imports
var (
	_ = networkingv1beta1.VirtualService{}
	_ = securityv1beta1.AuthorizationPolicy{}
)
