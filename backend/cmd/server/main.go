package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ivikasavnish/istio-ui/internal/api"
	"github.com/ivikasavnish/istio-ui/internal/istio"
	"github.com/ivikasavnish/istio-ui/internal/k8s"
	"github.com/ivikasavnish/istio-ui/internal/scheduler"
)

func main() {
	// Initialize Kubernetes client
	k8sClient, err := k8s.NewClient()
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes client: %v", err)
	}

	// Initialize Istio client
	istioClient, err := istio.NewClient(k8sClient)
	if err != nil {
		log.Fatalf("Failed to initialize Istio client: %v", err)
	}

	// Initialize scheduler
	sched := scheduler.NewScheduler(istioClient)
	sched.Start()
	defer sched.Stop()

	// Setup Gin router
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		// VirtualService routes
		v1.GET("/virtualservices", api.ListVirtualServices(istioClient))
		v1.GET("/virtualservices/:namespace/:name", api.GetVirtualService(istioClient))
		v1.POST("/virtualservices", api.CreateVirtualService(istioClient))
		v1.PUT("/virtualservices/:namespace/:name", api.UpdateVirtualService(istioClient))
		v1.DELETE("/virtualservices/:namespace/:name", api.DeleteVirtualService(istioClient))

		// DestinationRule routes
		v1.GET("/destinationrules", api.ListDestinationRules(istioClient))
		v1.GET("/destinationrules/:namespace/:name", api.GetDestinationRule(istioClient))
		v1.POST("/destinationrules", api.CreateDestinationRule(istioClient))
		v1.PUT("/destinationrules/:namespace/:name", api.UpdateDestinationRule(istioClient))
		v1.DELETE("/destinationrules/:namespace/:name", api.DeleteDestinationRule(istioClient))

		// Gateway routes
		v1.GET("/gateways", api.ListGateways(istioClient))
		v1.GET("/gateways/:namespace/:name", api.GetGateway(istioClient))
		v1.POST("/gateways", api.CreateGateway(istioClient))
		v1.PUT("/gateways/:namespace/:name", api.UpdateGateway(istioClient))
		v1.DELETE("/gateways/:namespace/:name", api.DeleteGateway(istioClient))

		// ServiceEntry routes
		v1.GET("/serviceentries", api.ListServiceEntries(istioClient))
		v1.GET("/serviceentries/:namespace/:name", api.GetServiceEntry(istioClient))
		v1.POST("/serviceentries", api.CreateServiceEntry(istioClient))
		v1.PUT("/serviceentries/:namespace/:name", api.UpdateServiceEntry(istioClient))
		v1.DELETE("/serviceentries/:namespace/:name", api.DeleteServiceEntry(istioClient))

		// AuthorizationPolicy routes
		v1.GET("/authorizationpolicies", api.ListAuthorizationPolicies(istioClient))
		v1.GET("/authorizationpolicies/:namespace/:name", api.GetAuthorizationPolicy(istioClient))
		v1.POST("/authorizationpolicies", api.CreateAuthorizationPolicy(istioClient))
		v1.PUT("/authorizationpolicies/:namespace/:name", api.UpdateAuthorizationPolicy(istioClient))
		v1.DELETE("/authorizationpolicies/:namespace/:name", api.DeleteAuthorizationPolicy(istioClient))

		// PeerAuthentication routes
		v1.GET("/peerauthentications", api.ListPeerAuthentications(istioClient))
		v1.GET("/peerauthentications/:namespace/:name", api.GetPeerAuthentication(istioClient))
		v1.POST("/peerauthentications", api.CreatePeerAuthentication(istioClient))
		v1.PUT("/peerauthentications/:namespace/:name", api.UpdatePeerAuthentication(istioClient))
		v1.DELETE("/peerauthentications/:namespace/:name", api.DeletePeerAuthentication(istioClient))

		// RequestAuthentication routes
		v1.GET("/requestauthentications", api.ListRequestAuthentications(istioClient))
		v1.GET("/requestauthentications/:namespace/:name", api.GetRequestAuthentication(istioClient))
		v1.POST("/requestauthentications", api.CreateRequestAuthentication(istioClient))
		v1.PUT("/requestauthentications/:namespace/:name", api.UpdateRequestAuthentication(istioClient))
		v1.DELETE("/requestauthentications/:namespace/:name", api.DeleteRequestAuthentication(istioClient))

		// Scheduled action routes
		v1.GET("/schedules", api.ListSchedules(sched))
		v1.POST("/schedules", api.CreateSchedule(sched))
		v1.DELETE("/schedules/:id", api.DeleteSchedule(sched))

		// Namespace routes
		v1.GET("/namespaces", api.ListNamespaces(k8sClient))

		// Service routes
		v1.GET("/services", api.ListServices(k8sClient))
		v1.GET("/services/:namespace/:name", api.GetService(k8sClient))

		// Topology/Graph routes
		v1.GET("/topology", api.GetTopology(istioClient, k8sClient))
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("MeshControl Center API Server starting on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
