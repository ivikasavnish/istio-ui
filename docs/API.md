# MeshControl Center API Documentation

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Currently, the API does not require authentication for local development. In production, you should configure appropriate authentication mechanisms.

## Endpoints

### Virtual Services

#### List Virtual Services
```
GET /virtualservices?namespace={namespace}
```

Query Parameters:
- `namespace` (optional): Filter by namespace. If empty, returns resources from all namespaces.

Response:
```json
[
  {
    "apiVersion": "networking.istio.io/v1beta1",
    "kind": "VirtualService",
    "metadata": {
      "name": "reviews",
      "namespace": "default"
    },
    "spec": {
      "hosts": ["reviews"],
      "http": [...]
    }
  }
]
```

#### Get Virtual Service
```
GET /virtualservices/{namespace}/{name}
```

#### Create Virtual Service
```
POST /virtualservices
Content-Type: application/json

{
  "apiVersion": "networking.istio.io/v1beta1",
  "kind": "VirtualService",
  "metadata": {
    "name": "reviews",
    "namespace": "default"
  },
  "spec": {
    "hosts": ["reviews"],
    "http": [
      {
        "route": [
          {
            "destination": {
              "host": "reviews",
              "subset": "v1"
            }
          }
        ]
      }
    ]
  }
}
```

#### Update Virtual Service
```
PUT /virtualservices/{namespace}/{name}
Content-Type: application/json
```

#### Delete Virtual Service
```
DELETE /virtualservices/{namespace}/{name}
```

### Destination Rules

Similar endpoints as Virtual Services:
- `GET /destinationrules`
- `GET /destinationrules/{namespace}/{name}`
- `POST /destinationrules`
- `PUT /destinationrules/{namespace}/{name}`
- `DELETE /destinationrules/{namespace}/{name}`

### Gateways

Similar endpoints as Virtual Services:
- `GET /gateways`
- `GET /gateways/{namespace}/{name}`
- `POST /gateways`
- `PUT /gateways/{namespace}/{name}`
- `DELETE /gateways/{namespace}/{name}`

### Service Entries

Similar endpoints as Virtual Services:
- `GET /serviceentries`
- `GET /serviceentries/{namespace}/{name}`
- `POST /serviceentries`
- `PUT /serviceentries/{namespace}/{name}`
- `DELETE /serviceentries/{namespace}/{name}`

### Authorization Policies

Similar endpoints as Virtual Services:
- `GET /authorizationpolicies`
- `GET /authorizationpolicies/{namespace}/{name}`
- `POST /authorizationpolicies`
- `PUT /authorizationpolicies/{namespace}/{name}`
- `DELETE /authorizationpolicies/{namespace}/{name}`

### Peer Authentications

Similar endpoints as Virtual Services:
- `GET /peerauthentications`
- `GET /peerauthentications/{namespace}/{name}`
- `POST /peerauthentications`
- `PUT /peerauthentications/{namespace}/{name}`
- `DELETE /peerauthentications/{namespace}/{name}`

### Request Authentications

Similar endpoints as Virtual Services:
- `GET /requestauthentications`
- `GET /requestauthentications/{namespace}/{name}`
- `POST /requestauthentications`
- `PUT /requestauthentications/{namespace}/{name}`
- `DELETE /requestauthentications/{namespace}/{name}`

### Scheduled Actions

#### List Schedules
```
GET /schedules
```

Response:
```json
[
  {
    "id": "uuid",
    "name": "Canary Rollout",
    "description": "Gradually shift traffic to v2",
    "schedule": "0 */6 * * *",
    "action_type": "update_virtualservice",
    "namespace": "default",
    "resource_name": "reviews",
    "payload": {},
    "created_at": "2023-11-20T10:00:00Z",
    "next_run": "2023-11-20T12:00:00Z",
    "enabled": true
  }
]
```

#### Create Schedule
```
POST /schedules
Content-Type: application/json

{
  "name": "Canary Rollout",
  "description": "Gradually shift traffic to v2",
  "schedule": "0 */6 * * *",
  "action_type": "update_virtualservice",
  "namespace": "default",
  "resource_name": "reviews",
  "payload": {
    "spec": {
      "hosts": ["reviews"],
      "http": [{
        "route": [{
          "destination": {"host": "reviews", "subset": "v2"},
          "weight": 50
        }]
      }]
    }
  }
}
```

#### Delete Schedule
```
DELETE /schedules/{id}
```

### Kubernetes Resources

#### List Namespaces
```
GET /namespaces
```

Response:
```json
["default", "kube-system", "istio-system"]
```

#### List Services
```
GET /services?namespace={namespace}
```

#### Get Service
```
GET /services/{namespace}/{name}
```

### Topology

#### Get Topology
```
GET /topology?namespace={namespace}
```

Response:
```json
{
  "nodes": [
    {
      "name": "reviews",
      "namespace": "default",
      "type": "service",
      "labels": {"app": "reviews"}
    }
  ],
  "edges": [
    {
      "source": "default/reviews-vs",
      "target": "default/reviews",
      "type": "http"
    }
  ]
}
```

## Error Responses

All endpoints may return error responses in the following format:

```json
{
  "error": "Error message describing what went wrong"
}
```

Common HTTP status codes:
- `200 OK` - Successful request
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request body or parameters
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Cron Schedule Format

The scheduled actions use standard cron expression format:

```
┌───────────── minute (0 - 59)
│ ┌───────────── hour (0 - 23)
│ │ ┌───────────── day of the month (1 - 31)
│ │ │ ┌───────────── month (1 - 12)
│ │ │ │ ┌───────────── day of the week (0 - 6) (Sunday to Saturday)
│ │ │ │ │
* * * * *
```

Examples:
- `0 0 * * *` - Daily at midnight
- `0 */6 * * *` - Every 6 hours
- `*/15 * * * *` - Every 15 minutes
- `0 9 * * 1-5` - Weekdays at 9 AM
