package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivikasavnish/istio-ui/internal/istio"
	networkingclient "istio.io/client-go/pkg/apis/networking/v1beta1"
)

// ListVirtualServices returns all VirtualServices
func ListVirtualServices(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", "")
		
		vsList, err := client.ListVirtualServices(c.Request.Context(), namespace)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, vsList.Items)
	}
}

// GetVirtualService returns a specific VirtualService
func GetVirtualService(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		vs, err := client.GetVirtualService(c.Request.Context(), namespace, name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, vs)
	}
}

// CreateVirtualService creates a new VirtualService
func CreateVirtualService(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var vs networkingclient.VirtualService
		if err := c.ShouldBindJSON(&vs); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		created, err := client.CreateVirtualService(c.Request.Context(), &vs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusCreated, created)
	}
}

// UpdateVirtualService updates an existing VirtualService
func UpdateVirtualService(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		var vs networkingclient.VirtualService
		if err := c.ShouldBindJSON(&vs); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		vs.Namespace = namespace
		vs.Name = name
		
		updated, err := client.UpdateVirtualService(c.Request.Context(), &vs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, updated)
	}
}

// DeleteVirtualService deletes a VirtualService
func DeleteVirtualService(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		err := client.DeleteVirtualService(c.Request.Context(), namespace, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "VirtualService deleted successfully"})
	}
}
