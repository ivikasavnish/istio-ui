package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivikasavnish/istio-ui/internal/istio"
	networkingclient "istio.io/client-go/pkg/apis/networking/v1beta1"
)

// ListGateways returns all Gateways
func ListGateways(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", "")
		
		gwList, err := client.ListGateways(c.Request.Context(), namespace)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gwList.Items)
	}
}

// GetGateway returns a specific Gateway
func GetGateway(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		gw, err := client.GetGateway(c.Request.Context(), namespace, name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gw)
	}
}

// CreateGateway creates a new Gateway
func CreateGateway(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var gw networkingclient.Gateway
		if err := c.ShouldBindJSON(&gw); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		created, err := client.CreateGateway(c.Request.Context(), &gw)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusCreated, created)
	}
}

// UpdateGateway updates an existing Gateway
func UpdateGateway(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		var gw networkingclient.Gateway
		if err := c.ShouldBindJSON(&gw); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		gw.Namespace = namespace
		gw.Name = name
		
		updated, err := client.UpdateGateway(c.Request.Context(), &gw)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, updated)
	}
}

// DeleteGateway deletes a Gateway
func DeleteGateway(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		err := client.DeleteGateway(c.Request.Context(), namespace, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "Gateway deleted successfully"})
	}
}

// ListServiceEntries returns all ServiceEntries
func ListServiceEntries(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", "")
		
		seList, err := client.ListServiceEntries(c.Request.Context(), namespace)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, seList.Items)
	}
}

// GetServiceEntry returns a specific ServiceEntry
func GetServiceEntry(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		se, err := client.GetServiceEntry(c.Request.Context(), namespace, name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, se)
	}
}

// CreateServiceEntry creates a new ServiceEntry
func CreateServiceEntry(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var se networkingclient.ServiceEntry
		if err := c.ShouldBindJSON(&se); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		created, err := client.CreateServiceEntry(c.Request.Context(), &se)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusCreated, created)
	}
}

// UpdateServiceEntry updates an existing ServiceEntry
func UpdateServiceEntry(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		var se networkingclient.ServiceEntry
		if err := c.ShouldBindJSON(&se); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		se.Namespace = namespace
		se.Name = name
		
		updated, err := client.UpdateServiceEntry(c.Request.Context(), &se)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, updated)
	}
}

// DeleteServiceEntry deletes a ServiceEntry
func DeleteServiceEntry(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		err := client.DeleteServiceEntry(c.Request.Context(), namespace, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "ServiceEntry deleted successfully"})
	}
}
