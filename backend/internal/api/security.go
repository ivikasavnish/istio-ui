package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivikasavnish/istio-ui/internal/istio"
	securityclient "istio.io/client-go/pkg/apis/security/v1beta1"
)

// ListAuthorizationPolicies returns all AuthorizationPolicies
func ListAuthorizationPolicies(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", "")
		
		apList, err := client.ListAuthorizationPolicies(c.Request.Context(), namespace)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, apList.Items)
	}
}

// GetAuthorizationPolicy returns a specific AuthorizationPolicy
func GetAuthorizationPolicy(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		ap, err := client.GetAuthorizationPolicy(c.Request.Context(), namespace, name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, ap)
	}
}

// CreateAuthorizationPolicy creates a new AuthorizationPolicy
func CreateAuthorizationPolicy(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ap securityclient.AuthorizationPolicy
		if err := c.ShouldBindJSON(&ap); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		created, err := client.CreateAuthorizationPolicy(c.Request.Context(), &ap)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusCreated, created)
	}
}

// UpdateAuthorizationPolicy updates an existing AuthorizationPolicy
func UpdateAuthorizationPolicy(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		var ap securityclient.AuthorizationPolicy
		if err := c.ShouldBindJSON(&ap); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		ap.Namespace = namespace
		ap.Name = name
		
		updated, err := client.UpdateAuthorizationPolicy(c.Request.Context(), &ap)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, updated)
	}
}

// DeleteAuthorizationPolicy deletes an AuthorizationPolicy
func DeleteAuthorizationPolicy(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		err := client.DeleteAuthorizationPolicy(c.Request.Context(), namespace, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "AuthorizationPolicy deleted successfully"})
	}
}

// ListPeerAuthentications returns all PeerAuthentications
func ListPeerAuthentications(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", "")
		
		paList, err := client.ListPeerAuthentications(c.Request.Context(), namespace)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, paList.Items)
	}
}

// GetPeerAuthentication returns a specific PeerAuthentication
func GetPeerAuthentication(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		pa, err := client.GetPeerAuthentication(c.Request.Context(), namespace, name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, pa)
	}
}

// CreatePeerAuthentication creates a new PeerAuthentication
func CreatePeerAuthentication(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var pa securityclient.PeerAuthentication
		if err := c.ShouldBindJSON(&pa); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		created, err := client.CreatePeerAuthentication(c.Request.Context(), &pa)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusCreated, created)
	}
}

// UpdatePeerAuthentication updates an existing PeerAuthentication
func UpdatePeerAuthentication(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		var pa securityclient.PeerAuthentication
		if err := c.ShouldBindJSON(&pa); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		pa.Namespace = namespace
		pa.Name = name
		
		updated, err := client.UpdatePeerAuthentication(c.Request.Context(), &pa)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, updated)
	}
}

// DeletePeerAuthentication deletes a PeerAuthentication
func DeletePeerAuthentication(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		err := client.DeletePeerAuthentication(c.Request.Context(), namespace, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "PeerAuthentication deleted successfully"})
	}
}

// ListRequestAuthentications returns all RequestAuthentications
func ListRequestAuthentications(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", "")
		
		raList, err := client.ListRequestAuthentications(c.Request.Context(), namespace)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, raList.Items)
	}
}

// GetRequestAuthentication returns a specific RequestAuthentication
func GetRequestAuthentication(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		ra, err := client.GetRequestAuthentication(c.Request.Context(), namespace, name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, ra)
	}
}

// CreateRequestAuthentication creates a new RequestAuthentication
func CreateRequestAuthentication(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ra securityclient.RequestAuthentication
		if err := c.ShouldBindJSON(&ra); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		created, err := client.CreateRequestAuthentication(c.Request.Context(), &ra)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusCreated, created)
	}
}

// UpdateRequestAuthentication updates an existing RequestAuthentication
func UpdateRequestAuthentication(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		var ra securityclient.RequestAuthentication
		if err := c.ShouldBindJSON(&ra); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		ra.Namespace = namespace
		ra.Name = name
		
		updated, err := client.UpdateRequestAuthentication(c.Request.Context(), &ra)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, updated)
	}
}

// DeleteRequestAuthentication deletes a RequestAuthentication
func DeleteRequestAuthentication(client *istio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		
		err := client.DeleteRequestAuthentication(c.Request.Context(), namespace, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "RequestAuthentication deleted successfully"})
	}
}
