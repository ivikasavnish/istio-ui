package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivikasavnish/istio-ui/internal/helm"
	"github.com/ivikasavnish/istio-ui/internal/k8s"
)

// ListHelmReleases returns all Helm releases
func ListHelmReleases(k8sClient *k8s.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", "default")

		helmManager := helm.NewManager(k8sClient.Config, namespace)
		releases, err := helmManager.ListReleases(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, releases)
	}
}

// GetHelmRelease returns details of a specific Helm release
func GetHelmRelease(k8sClient *k8s.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Query("namespace")
		if namespace == "" {
			namespace = "default"
		}
		name := c.Param("name")

		helmManager := helm.NewManager(k8sClient.Config, namespace)
		release, err := helmManager.GetRelease(c.Request.Context(), name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, release)
	}
}

// InstallHelmChart installs a Helm chart
func InstallHelmChart(k8sClient *k8s.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name      string                 `json:"name" binding:"required"`
			Namespace string                 `json:"namespace"`
			ChartPath string                 `json:"chart_path" binding:"required"`
			Values    map[string]interface{} `json:"values"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Namespace == "" {
			req.Namespace = "default"
		}

		helmManager := helm.NewManager(k8sClient.Config, req.Namespace)
		release, err := helmManager.InstallChart(c.Request.Context(), req.Name, req.ChartPath, req.Values)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Chart installed successfully",
			"release": release.Name,
		})
	}
}

// UpgradeHelmRelease upgrades a Helm release
func UpgradeHelmRelease(k8sClient *k8s.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")

		var req struct {
			Namespace string                 `json:"namespace"`
			ChartPath string                 `json:"chart_path" binding:"required"`
			Values    map[string]interface{} `json:"values"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Namespace == "" {
			req.Namespace = "default"
		}

		helmManager := helm.NewManager(k8sClient.Config, req.Namespace)
		release, err := helmManager.UpgradeChart(c.Request.Context(), name, req.ChartPath, req.Values)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Release upgraded successfully",
			"release": release.Name,
		})
	}
}

// UninstallHelmRelease uninstalls a Helm release
func UninstallHelmRelease(k8sClient *k8s.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Query("namespace")
		if namespace == "" {
			namespace = "default"
		}
		name := c.Param("name")

		helmManager := helm.NewManager(k8sClient.Config, namespace)
		err := helmManager.UninstallRelease(c.Request.Context(), name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Release uninstalled successfully"})
	}
}

// RollbackHelmRelease rolls back a Helm release
func RollbackHelmRelease(k8sClient *k8s.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")

		var req struct {
			Namespace string `json:"namespace"`
			Version   int    `json:"version" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Namespace == "" {
			req.Namespace = "default"
		}

		helmManager := helm.NewManager(k8sClient.Config, req.Namespace)
		err := helmManager.RollbackRelease(c.Request.Context(), name, req.Version)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Release rolled back successfully"})
	}
}

// GetHelmReleaseHistory returns the history of a Helm release
func GetHelmReleaseHistory(k8sClient *k8s.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Query("namespace")
		if namespace == "" {
			namespace = "default"
		}
		name := c.Param("name")

		helmManager := helm.NewManager(k8sClient.Config, namespace)
		history, err := helmManager.GetReleaseHistory(c.Request.Context(), name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, history)
	}
}
