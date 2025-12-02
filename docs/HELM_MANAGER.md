# Helm Manager

## Overview

The Helm Manager provides a visual interface for managing Helm releases in your Kubernetes cluster. It allows you to install, upgrade, rollback, and uninstall Helm charts directly from the MeshControl Center UI.

## Features

- **List Releases**: View all Helm releases in a namespace
- **Install Charts**: Deploy new Helm charts with custom values
- **Upgrade Releases**: Update existing releases to new versions
- **Uninstall Releases**: Remove Helm releases
- **Rollback**: Revert to previous release versions
- **Release History**: View release history and revisions
- **Namespace Filtering**: Filter releases by namespace

## Usage

### In the UI

Navigate to "Helm Manager" from the sidebar menu.

#### Installing a Chart

1. Select the target namespace from the dropdown
2. Click "Install Chart"
3. Fill in the form:
   - **Release Name**: Unique name for the release
   - **Namespace**: Target namespace
   - **Chart Path or URL**: Path to the chart (e.g., `bitnami/nginx`)
   - **Values**: JSON object with custom values (optional)
4. Click "Install"

Example values:
```json
{
  "replicaCount": 2,
  "service": {
    "type": "LoadBalancer"
  }
}
```

#### Upgrading a Release

1. Click the upgrade icon next to a release
2. Specify the new chart path/version
3. Optionally update values
4. Click "Upgrade"

#### Uninstalling a Release

1. Click the delete icon next to a release
2. Confirm the deletion

### API Endpoints

#### List Helm Releases
```http
GET /api/v1/helm/releases?namespace=default
```

Response:
```json
[
  {
    "name": "nginx-ingress",
    "namespace": "default",
    "version": 1,
    "status": "deployed",
    "chart": "nginx-ingress",
    "app_version": "1.0.0",
    "updated": "2023-12-02 17:30:00"
  }
]
```

#### Get Release Details
```http
GET /api/v1/helm/releases/nginx-ingress?namespace=default
```

#### Install Chart
```http
POST /api/v1/helm/releases
Content-Type: application/json

{
  "name": "my-nginx",
  "namespace": "default",
  "chart_path": "bitnami/nginx",
  "values": {
    "replicaCount": 2
  }
}
```

#### Upgrade Release
```http
PUT /api/v1/helm/releases/my-nginx
Content-Type: application/json

{
  "namespace": "default",
  "chart_path": "bitnami/nginx",
  "values": {
    "replicaCount": 3
  }
}
```

#### Uninstall Release
```http
DELETE /api/v1/helm/releases/my-nginx?namespace=default
```

#### Rollback Release
```http
POST /api/v1/helm/releases/my-nginx/rollback
Content-Type: application/json

{
  "namespace": "default",
  "version": 1
}
```

#### Get Release History
```http
GET /api/v1/helm/releases/my-nginx/history?namespace=default
```

## Chart Sources

### Helm Repository Charts

Use chart names from added Helm repositories:
```
bitnami/nginx
stable/mysql
```

### Local Charts

Use absolute or relative paths to local charts:
```
/path/to/my-chart
./charts/my-app
```

### Chart URLs

Use direct URLs to chart archives:
```
https://charts.example.com/nginx-1.0.0.tgz
```

## Values Configuration

### JSON Format

Values must be provided as valid JSON:

```json
{
  "image": {
    "repository": "nginx",
    "tag": "1.21.0"
  },
  "service": {
    "type": "LoadBalancer",
    "port": 80
  },
  "ingress": {
    "enabled": true,
    "hosts": ["example.com"]
  }
}
```

### Common Values

#### Setting Replicas
```json
{
  "replicaCount": 3
}
```

#### Resource Limits
```json
{
  "resources": {
    "limits": {
      "cpu": "100m",
      "memory": "128Mi"
    }
  }
}
```

#### Environment Variables
```json
{
  "env": [
    {
      "name": "MY_VAR",
      "value": "my-value"
    }
  ]
}
```

## Implementation Details

### Backend

The Helm management is handled by the `Manager` in `internal/helm/manager.go`:

- Uses Helm v3 SDK
- Supports all standard Helm operations
- Integrates with Kubernetes contexts
- Provides structured release information

### Frontend

The `HelmManager` page (`frontend/src/pages/HelmManager.js`):

- DataGrid for release listing
- Install/upgrade dialogs
- JSON values editor
- Namespace filtering

## Best Practices

1. **Always Review Values**: Check the chart's default values before installing
2. **Use Namespaces**: Isolate releases by namespace
3. **Test in Dev First**: Test chart installations in development before production
4. **Keep Charts Updated**: Regularly upgrade to get security fixes
5. **Use Version Pinning**: Specify exact chart versions for production

## Troubleshooting

### Chart Not Found

If you get "chart not found" errors:
1. Verify the chart name is correct
2. Check that Helm repositories are added (`helm repo add`)
3. For local charts, ensure the path is accessible
4. Try using a direct chart URL

### Installation Fails

Common installation issues:
1. **Invalid Values**: Check JSON syntax in values field
2. **Resource Conflicts**: Release name or resources already exist
3. **Permission Issues**: Insufficient RBAC permissions
4. **Namespace Issues**: Target namespace doesn't exist

### Release in Failed State

To recover from a failed release:
1. Check release history
2. Rollback to a previous version
3. Or uninstall and reinstall

## Security Considerations

- Helm operations require appropriate Kubernetes RBAC permissions
- Values may contain sensitive data (passwords, API keys)
- Chart sources should be trusted
- Review chart contents before installation
- Use namespace isolation for sensitive workloads

## Future Enhancements

- [ ] Helm repository management (add/remove repos)
- [ ] Chart value schema validation
- [ ] Visual diff before upgrade
- [ ] Release notes display
- [ ] Chart search from repositories
- [ ] Automated chart updates
- [ ] Release health monitoring
- [ ] Integration with artifact registries
