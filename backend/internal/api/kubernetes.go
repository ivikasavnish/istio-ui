package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivikasavnish/istio-ui/internal/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListNamespaces returns all namespaces
func ListNamespaces(client *k8s.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespaces, err := client.Clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		names := make([]string, 0, len(namespaces.Items))
		for _, ns := range namespaces.Items {
			names = append(names, ns.Name)
		}
		
		c.JSON(http.StatusOK, names)
	}
}

// ListServices returns all services
func ListServices(client *k8s.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", "")
		
		services, err := client.Clientset.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, services.Items)
	}
}

// GetService returns a specific service
func GetService(client *k8s.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		service, err := client.Clientset.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, service)
	}
}
