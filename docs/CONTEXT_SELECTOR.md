# Kubernetes Context Selector

## Overview

The Kubernetes Context Selector allows users to switch between different Kubernetes contexts directly from the MeshControl Center UI. This is useful when managing multiple clusters or environments.

## Features

- **List Contexts**: View all available Kubernetes contexts from your kubeconfig
- **Current Context Indicator**: Clearly shows which context is currently active
- **Switch Contexts**: Change contexts with a single click
- **Automatic Reload**: After switching contexts, the page reloads to fetch data from the new cluster

## Usage

### In the UI

The context selector is located in the top-right corner of the application header.

1. Click on the dropdown to see all available contexts
2. The current context is marked with a "Current" chip
3. Select a different context to switch
4. The application will reload with data from the new context

### API Endpoints

#### List Contexts
```http
GET /api/v1/contexts
```

Response:
```json
[
  {
    "name": "minikube",
    "cluster": "minikube",
    "user": "minikube",
    "namespace": "default",
    "is_current": true
  },
  {
    "name": "production",
    "cluster": "prod-cluster",
    "user": "prod-admin",
    "namespace": "default",
    "is_current": false
  }
]
```

#### Get Current Context
```http
GET /api/v1/contexts/current
```

Response:
```json
{
  "context": "minikube"
}
```

#### Switch Context
```http
POST /api/v1/contexts/switch
Content-Type: application/json

{
  "context": "production"
}
```

Response:
```json
{
  "message": "Context switched successfully",
  "context": "production"
}
```

## Implementation Details

### Backend

The context management is handled by the `ContextManager` in `internal/k8s/context.go`:

- Reads kubeconfig file
- Lists all available contexts
- Switches between contexts
- Creates new Kubernetes clients for different contexts

### Frontend

The `ContextSelector` component (`frontend/src/components/ContextSelector.js`):

- Fetches available contexts from the API
- Displays them in a Material-UI Select dropdown
- Handles context switching
- Reloads the page after switching

## Security Considerations

- The application must have read access to the kubeconfig file
- Context switching affects all subsequent API calls
- Users should have appropriate permissions in each context
- No authentication is currently implemented (planned for future versions)

## Troubleshooting

### Context Not Switching

If the context doesn't switch:
1. Check that the kubeconfig file is accessible
2. Verify the context exists in your kubeconfig
3. Ensure you have valid credentials for the target context
4. Check the browser console for errors

### No Contexts Available

If no contexts are listed:
1. Verify that `KUBECONFIG` environment variable is set or `~/.kube/config` exists
2. Check that the kubeconfig file is valid
3. Ensure the backend has read permissions on the file

## Future Enhancements

- [ ] Context-specific namespace filtering
- [ ] Multi-cluster view (simultaneous access to multiple contexts)
- [ ] Context bookmarks/favorites
- [ ] Context health indicators
- [ ] Per-context authentication tokens
