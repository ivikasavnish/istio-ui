# MeshControl Center API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Currently uses basic authentication. Future versions will support JWT tokens.

## Endpoints

### Services

#### List All Services
```http
GET /services
```

**Response:**
```json
[
  {
    "id": 1,
    "namespace": "default",
    "name": "my-service",
    "version": "v1",
    "last_seen": "2024-01-15T10:30:00Z"
  }
]
```

#### List Services in Namespace
```http
GET /services/{namespace}
```

**Response:**
```json
{
  "namespace": "default",
  "services": [
    {
      "name": "my-service",
      "labels": {
        "app": "my-service",
        "version": "v1"
      },
      "ports": [8080, 8443]
    }
  ]
}
```

### VirtualServices

#### Create VirtualService
```http
POST /virtualservices
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "my-service",
  "namespace": "default",
  "hosts": ["my-service.default.svc.cluster.local"],
  "http": [
    {
      "route": [
        {
          "destination": {
            "host": "my-service",
            "subset": "v1"
          },
          "weight": 80
        },
        {
          "destination": {
            "host": "my-service",
            "subset": "v2"
          },
          "weight": 20
        }
      ]
    }
  ]
}
```

**Response:**
```json
{
  "message": "VirtualService created successfully"
}
```

#### List VirtualServices
```http
GET /virtualservices/{namespace}
```

#### Get VirtualService
```http
GET /virtualservices/{namespace}/{name}
```

#### Update VirtualService
```http
PUT /virtualservices/{namespace}/{name}
Content-Type: application/json
```

#### Delete VirtualService
```http
DELETE /virtualservices/{namespace}/{name}
```

### DestinationRules

#### Create DestinationRule
```http
POST /destinationrules
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "my-service",
  "namespace": "default",
  "host": "my-service.default.svc.cluster.local",
  "subsets": [
    {
      "name": "v1",
      "labels": {
        "version": "v1"
      }
    },
    {
      "name": "v2",
      "labels": {
        "version": "v2"
      }
    }
  ],
  "trafficPolicy": {
    "connectionPool": {
      "tcp": {
        "maxConnections": 100
      },
      "http": {
        "http1MaxPendingRequests": 10,
        "http2MaxRequests": 100
      }
    },
    "outlierDetection": {
      "consecutiveErrors": 5,
      "interval": "30s",
      "baseEjectionTime": "30s",
      "maxEjectionPercent": 50
    }
  }
}
```

### PeerAuthentication

#### Create PeerAuthentication
```http
POST /peerauthentications
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "default",
  "namespace": "default",
  "mtls_mode": "STRICT"
}
```

### Traffic Management

#### Update Traffic Weights
```http
POST /traffic/weights
Content-Type: application/json
```

**Request Body:**
```json
{
  "service": "my-service",
  "namespace": "default",
  "weights": {
    "v1": 70,
    "v2": 30
  }
}
```

#### Configure Canary Deployment
```http
POST /traffic/canary
Content-Type: application/json
```

**Request Body:**
```json
{
  "service": "my-service",
  "namespace": "default",
  "old_version": "v1",
  "new_version": "v2",
  "percentage": 10
}
```

### Snapshots

#### Create Snapshot
```http
POST /snapshots
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "pre-deployment",
  "description": "Configuration before v2 rollout",
  "config": "{...}",
  "created_by": "admin"
}
```

#### List Snapshots
```http
GET /snapshots
```

#### Get Snapshot
```http
GET /snapshots/{id}
```

#### Delete Snapshot
```http
DELETE /snapshots/{id}
```

#### Restore Snapshot
```http
POST /snapshots/{id}/restore
```

### Scheduled Actions

#### Create Scheduled Action
```http
POST /scheduled-actions
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Daily traffic shift",
  "action_type": "traffic_shift",
  "cron_expr": "0 0 * * *",
  "config": "{\"service\": \"my-service\", \"weights\": {\"v1\": 100}}",
  "enabled": true,
  "next_run": "2024-01-16T00:00:00Z"
}
```

#### List Scheduled Actions
```http
GET /scheduled-actions
```

#### Update Scheduled Action
```http
PUT /scheduled-actions/{id}
Content-Type: application/json
```

#### Delete Scheduled Action
```http
DELETE /scheduled-actions/{id}
```

### Audit Logs

#### List Audit Logs
```http
GET /audit-logs?limit=100
```

**Response:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "username": "admin",
    "action": "CREATE",
    "resource": "VirtualService",
    "details": "default/my-service",
    "success": true,
    "timestamp": "2024-01-15T10:30:00Z"
  }
]
```

## WebSocket Events

Connect to `/ws` for real-time events.

### Event Types

- `connected` - Initial connection
- `virtualservice_created` - New VirtualService created
- `virtualservice_updated` - VirtualService updated
- `virtualservice_deleted` - VirtualService deleted
- `destinationrule_created` - New DestinationRule created
- `gateway_created` - New Gateway created
- `traffic_weights_updated` - Traffic weights changed
- `snapshot_created` - New snapshot created
- `scheduled_action_created` - New scheduled action created

### Example Event

```json
{
  "type": "virtualservice_created",
  "payload": {
    "name": "my-service",
    "namespace": "default",
    "hosts": ["my-service"]
  }
}
```

## Error Responses

All error responses follow this format:

```json
{
  "error": "Error message description"
}
```

**Common HTTP Status Codes:**
- `200 OK` - Success
- `201 Created` - Resource created
- `400 Bad Request` - Invalid request
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error
- `503 Service Unavailable` - Service unavailable (e.g., K8s not connected)
