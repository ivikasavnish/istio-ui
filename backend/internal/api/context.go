package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivikasavnish/istio-ui/internal/k8s"
)

// GetContexts returns all available Kubernetes contexts
func GetContexts() gin.HandlerFunc {
	return func(c *gin.Context) {
		kubeconfig, err := k8s.GetKubeconfigPath()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get kubeconfig path: " + err.Error()})
			return
		}

		cm, err := k8s.NewContextManager(kubeconfig)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create context manager: " + err.Error()})
			return
		}

		contexts := cm.GetContextsInfo()
		c.JSON(http.StatusOK, contexts)
	}
}

// GetCurrentContext returns the current context
func GetCurrentContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		kubeconfig, err := k8s.GetKubeconfigPath()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get kubeconfig path: " + err.Error()})
			return
		}

		cm, err := k8s.NewContextManager(kubeconfig)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create context manager: " + err.Error()})
			return
		}

		currentContext := cm.GetCurrentContext()
		c.JSON(http.StatusOK, gin.H{"context": currentContext})
	}
}

// SwitchContext switches to a different Kubernetes context
func SwitchContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Context string `json:"context" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		kubeconfig, err := k8s.GetKubeconfigPath()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get kubeconfig path: " + err.Error()})
			return
		}

		cm, err := k8s.NewContextManager(kubeconfig)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create context manager: " + err.Error()})
			return
		}

		_, err = cm.SwitchContext(req.Context)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to switch context: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Context switched successfully",
			"context": req.Context,
		})
	}
}
