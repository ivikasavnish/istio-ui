package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivikasavnish/istio-ui/internal/istio"
	"github.com/ivikasavnish/istio-ui/internal/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TopologyNode represents a node in the service mesh topology
type TopologyNode struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Type      string            `json:"type"` // service, virtualservice, gateway, etc.
	Labels    map[string]string `json:"labels,omitempty"`
}

// TopologyEdge represents a connection between nodes
type TopologyEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"` // http, tcp, grpc, etc.
}

// TopologyResponse contains the complete topology graph
type TopologyResponse struct {
	Nodes []TopologyNode `json:"nodes"`
	Edges []TopologyEdge `json:"edges"`
}

// GetTopology returns the service mesh topology
func GetTopology(istioClient *istio.Client, k8sClient *k8s.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", "")
		ctx := context.Background()
		
		topology := TopologyResponse{
			Nodes: []TopologyNode{},
			Edges: []TopologyEdge{},
		}
		
		// Get all services
		services, err := k8sClient.Clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
		if err == nil {
			for _, svc := range services.Items {
				topology.Nodes = append(topology.Nodes, TopologyNode{
					Name:      svc.Name,
					Namespace: svc.Namespace,
					Type:      "service",
					Labels:    svc.Labels,
				})
			}
		}
		
		// Get all VirtualServices
		virtualServices, err := istioClient.ListVirtualServices(ctx, namespace)
		if err == nil {
			for _, vs := range virtualServices.Items {
				topology.Nodes = append(topology.Nodes, TopologyNode{
					Name:      vs.Name,
					Namespace: vs.Namespace,
					Type:      "virtualservice",
					Labels:    vs.Labels,
				})
				
				// Create edges from VirtualService to destination services
				if vs.Spec.Http != nil {
					for _, httpRoute := range vs.Spec.Http {
						if httpRoute.Route != nil {
							for _, dest := range httpRoute.Route {
								if dest.Destination != nil && dest.Destination.Host != "" {
									topology.Edges = append(topology.Edges, TopologyEdge{
										Source: vs.Namespace + "/" + vs.Name,
										Target: vs.Namespace + "/" + dest.Destination.Host,
										Type:   "http",
									})
								}
							}
						}
					}
				}
			}
		}
		
		// Get all Gateways
		gateways, err := istioClient.ListGateways(ctx, namespace)
		if err == nil {
			for _, gw := range gateways.Items {
				topology.Nodes = append(topology.Nodes, TopologyNode{
					Name:      gw.Name,
					Namespace: gw.Namespace,
					Type:      "gateway",
					Labels:    gw.Labels,
				})
			}
		}
		
		// Get all DestinationRules
		destinationRules, err := istioClient.ListDestinationRules(ctx, namespace)
		if err == nil {
			for _, dr := range destinationRules.Items {
				topology.Nodes = append(topology.Nodes, TopologyNode{
					Name:      dr.Name,
					Namespace: dr.Namespace,
					Type:      "destinationrule",
					Labels:    dr.Labels,
				})
			}
		}
		
		c.JSON(http.StatusOK, topology)
	}
}
