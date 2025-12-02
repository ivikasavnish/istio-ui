package helm

import (
	"context"
	"fmt"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/client-go/rest"
)

// Manager manages Helm operations
type Manager struct {
	settings *cli.EnvSettings
	config   *rest.Config
}

// NewManager creates a new Helm manager
func NewManager(config *rest.Config, namespace string) *Manager {
	settings := cli.New()
	if namespace != "" {
		settings.SetNamespace(namespace)
	}
	
	return &Manager{
		settings: settings,
		config:   config,
	}
}

// Release represents a Helm release
type Release struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Version    int    `json:"version"`
	Status     string `json:"status"`
	Chart      string `json:"chart"`
	AppVersion string `json:"app_version"`
	Updated    string `json:"updated"`
}

// ListReleases lists all Helm releases in the namespace
func (m *Manager) ListReleases(ctx context.Context) ([]*Release, error) {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(m.settings.RESTClientGetter(), m.settings.Namespace(), os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {}); err != nil {
		return nil, fmt.Errorf("failed to initialize action config: %w", err)
	}

	client := action.NewList(actionConfig)
	client.AllNamespaces = false
	
	results, err := client.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	releases := make([]*Release, 0, len(results))
	for _, r := range results {
		releases = append(releases, &Release{
			Name:       r.Name,
			Namespace:  r.Namespace,
			Version:    r.Version,
			Status:     r.Info.Status.String(),
			Chart:      r.Chart.Metadata.Name,
			AppVersion: r.Chart.Metadata.AppVersion,
			Updated:    r.Info.LastDeployed.Format("2006-01-02 15:04:05"),
		})
	}

	return releases, nil
}

// GetRelease gets details of a specific release
func (m *Manager) GetRelease(ctx context.Context, name string) (*release.Release, error) {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(m.settings.RESTClientGetter(), m.settings.Namespace(), os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {}); err != nil {
		return nil, fmt.Errorf("failed to initialize action config: %w", err)
	}

	client := action.NewGet(actionConfig)
	
	result, err := client.Run(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get release: %w", err)
	}

	return result, nil
}

// InstallChart installs a Helm chart
func (m *Manager) InstallChart(ctx context.Context, releaseName, chartPath string, values map[string]interface{}) (*release.Release, error) {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(m.settings.RESTClientGetter(), m.settings.Namespace(), os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {}); err != nil {
		return nil, fmt.Errorf("failed to initialize action config: %w", err)
	}

	client := action.NewInstall(actionConfig)
	client.ReleaseName = releaseName
	client.Namespace = m.settings.Namespace()
	client.CreateNamespace = true

	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load chart: %w", err)
	}

	result, err := client.Run(chart, values)
	if err != nil {
		return nil, fmt.Errorf("failed to install chart: %w", err)
	}

	return result, nil
}

// UpgradeChart upgrades a Helm release
func (m *Manager) UpgradeChart(ctx context.Context, releaseName, chartPath string, values map[string]interface{}) (*release.Release, error) {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(m.settings.RESTClientGetter(), m.settings.Namespace(), os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {}); err != nil {
		return nil, fmt.Errorf("failed to initialize action config: %w", err)
	}

	client := action.NewUpgrade(actionConfig)
	client.Namespace = m.settings.Namespace()

	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load chart: %w", err)
	}

	result, err := client.Run(releaseName, chart, values)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade chart: %w", err)
	}

	return result, nil
}

// UninstallRelease uninstalls a Helm release
func (m *Manager) UninstallRelease(ctx context.Context, releaseName string) error {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(m.settings.RESTClientGetter(), m.settings.Namespace(), os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {}); err != nil {
		return fmt.Errorf("failed to initialize action config: %w", err)
	}

	client := action.NewUninstall(actionConfig)
	
	_, err := client.Run(releaseName)
	if err != nil {
		return fmt.Errorf("failed to uninstall release: %w", err)
	}

	return nil
}

// RollbackRelease rolls back a release to a previous version
func (m *Manager) RollbackRelease(ctx context.Context, releaseName string, version int) error {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(m.settings.RESTClientGetter(), m.settings.Namespace(), os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {}); err != nil {
		return fmt.Errorf("failed to initialize action config: %w", err)
	}

	client := action.NewRollback(actionConfig)
	client.Version = version
	
	if err := client.Run(releaseName); err != nil {
		return fmt.Errorf("failed to rollback release: %w", err)
	}

	return nil
}

// GetReleaseHistory gets the history of a release
func (m *Manager) GetReleaseHistory(ctx context.Context, releaseName string) ([]*release.Release, error) {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(m.settings.RESTClientGetter(), m.settings.Namespace(), os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {}); err != nil {
		return nil, fmt.Errorf("failed to initialize action config: %w", err)
	}

	client := action.NewHistory(actionConfig)
	
	results, err := client.Run(releaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to get release history: %w", err)
	}

	return results, nil
}
