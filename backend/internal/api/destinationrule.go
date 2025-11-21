package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivikasavnish/istio-ui/internal/istio"
	networkingclient "istio.io/client-go/pkg/apis/networking/v1beta1"
)

// ListDestinationRules returns all DestinationRules
func ListDestinationRules(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", "")
		
		drList, err := client.ListDestinationRules(c.Request.Context(), namespace)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, drList.Items)
	}
}

// GetDestinationRule returns a specific DestinationRule
func GetDestinationRule(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		dr, err := client.GetDestinationRule(c.Request.Context(), namespace, name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, dr)
	}
}

// CreateDestinationRule creates a new DestinationRule
func CreateDestinationRule(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dr networkingclient.DestinationRule
		if err := c.ShouldBindJSON(&dr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		created, err := client.CreateDestinationRule(c.Request.Context(), &dr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusCreated, created)
	}
}

// UpdateDestinationRule updates an existing DestinationRule
func UpdateDestinationRule(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		var dr networkingclient.DestinationRule
		if err := c.ShouldBindJSON(&dr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		dr.Namespace = namespace
		dr.Name = name
		
		updated, err := client.UpdateDestinationRule(c.Request.Context(), &dr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, updated)
	}
}

// DeleteDestinationRule deletes a DestinationRule
func DeleteDestinationRule(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		err := client.DeleteDestinationRule(c.Request.Context(), namespace, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "DestinationRule deleted successfully"})
	}
}
