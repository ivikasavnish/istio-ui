# MeshControl Center Architecture

## Overview

MeshControl Center is a full-stack web application designed to provide comprehensive management of Istio service mesh resources through an intuitive graphical interface.

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         User Interface                          │
│                    (React + Material-UI)                        │
└─────────────────────────────────────────────────────────────────┘
                                │
                                │ HTTP/REST
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API Gateway (Gin)                          │
│                         Port: 8080                              │
└─────────────────────────────────────────────────────────────────┘
                                │
                ┌───────────────┼───────────────┐
                │               │               │
                ▼               ▼               ▼
        ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
        │   Istio      │ │  Kubernetes  │ │  Scheduler   │
        │   Client     │ │   Client     │ │   Service    │
        └──────────────┘ └──────────────┘ └──────────────┘
                │               │               │
                └───────────────┼───────────────┘
                                │
                                ▼
                ┌────────────────────────────────┐
                │     Kubernetes API Server      │
                │         with Istio CRDs        │
                └────────────────────────────────┘
```

## Component Details

### Frontend (React Application)

**Technology Stack:**
- React 18
- Material-UI v5
- React Router v6
- Axios for API calls
- js-yaml for YAML processing

**Key Components:**
- **Layout**: Sidebar navigation with responsive design
- **Pages**: Dashboard, VirtualServices, DestinationRules, Gateways, Security Policies, Topology, Scheduled Actions
- **Services**: API client abstraction layer
- **Utils**: Helper functions for YAML conversion, date formatting

**Features:**
- Form-based resource creation and editing
- YAML preview and validation
- Real-time resource listing with DataGrid
- Responsive material design
- Error handling and user feedback

### Backend (Go Application)

**Technology Stack:**
- Go 1.21
- Gin web framework
- Kubernetes client-go
- Istio client-go
- Robfig cron scheduler

**Architecture Layers:**

1. **API Layer** (`internal/api/`)
   - REST endpoint handlers
   - Request validation
   - Response formatting
   - Error handling

2. **Client Layer** (`internal/istio/`, `internal/k8s/`)
   - Kubernetes API client wrapper
   - Istio CRD client wrapper
   - Resource CRUD operations
   - Connection management

3. **Scheduler Layer** (`internal/scheduler/`)
   - Cron-based task scheduling
   - Action execution engine
   - Schedule persistence
   - Task monitoring

4. **Application Layer** (`cmd/server/`)
   - Server initialization
   - Route configuration
   - Middleware setup
   - Lifecycle management

### Data Flow

#### Resource Read Operation
```
User → React UI → API Request → Gin Router → Handler
→ Istio Client → Kubernetes API → Istio CRD
→ Response ← Handler ← Client ← API ← CRD
→ Display ← UI ← JSON
```

#### Resource Create/Update Operation
```
User fills form → Validation → YAML generation
→ API POST/PUT → Handler → Validation
→ Istio Client → Create/Update CRD → K8s API
→ Success Response → UI Update → User feedback
```

#### Scheduled Action
```
Scheduler Cron Trigger → Action Executor
→ Fetch Current Resource → Apply Modifications
→ Update via Istio Client → K8s API
→ Log Result → Update Schedule Status
```

## Security Architecture

### RBAC (Kubernetes)
```yaml
ClusterRole:
  - networking.istio.io: [get, list, watch, create, update, delete]
  - security.istio.io: [get, list, watch, create, update, delete]
  - core: [get, list, watch] (namespaces, services)
```

### Authentication Flow (Future)
```
User → Login → JWT Token → API Request (with token)
→ Token Validation → RBAC Check → Resource Access
```

## Deployment Architecture

### Kubernetes Deployment
```
┌─────────────────────────────────────────────┐
│              Ingress Controller             │
│         (meshcontrol.example.com)           │
└─────────────────────────────────────────────┘
        │                           │
        │ /api                      │ /
        ▼                           ▼
┌──────────────────┐      ┌──────────────────┐
│  Backend Service │      │ Frontend Service │
│   (ClusterIP)    │      │   (ClusterIP)    │
└──────────────────┘      └──────────────────┘
        │                           │
        ▼                           ▼
┌──────────────────┐      ┌──────────────────┐
│  Backend Pods    │      │  Frontend Pods   │
│  (Go Server)     │      │  (Nginx + SPA)   │
└──────────────────┘      └──────────────────┘
        │
        ▼
┌──────────────────────────────────────────┐
│      Kubernetes API + Istio CRDs         │
└──────────────────────────────────────────┘
```

### Container Architecture
```
Backend Container:
  - Alpine Linux base
  - Go binary
  - No external dependencies
  - 50-100MB image size

Frontend Container:
  - Nginx alpine base
  - React build artifacts
  - Static file serving
  - 20-30MB image size
```

## Scalability Considerations

### Horizontal Scaling
- Backend: Stateless, can be scaled horizontally
- Frontend: Static content, highly cacheable
- Scheduler: Single instance with leader election (future)

### Performance Optimizations
- Connection pooling to K8s API
- Response caching for read operations
- Lazy loading of resources
- Pagination for large datasets

## Monitoring and Observability

### Metrics (Future Implementation)
- API request latency
- Resource operation counts
- Scheduler execution success rate
- Error rates by endpoint

### Logging
- Structured logging (JSON format)
- Request/response logging
- Error stack traces
- Scheduler action audit logs

### Health Checks
- `/health` endpoint for liveness
- Kubernetes readiness probes
- API connectivity checks

## API Design Patterns

### RESTful Conventions
```
GET    /api/v1/{resource}                    - List all
GET    /api/v1/{resource}/{namespace}/{name} - Get one
POST   /api/v1/{resource}                    - Create
PUT    /api/v1/{resource}/{namespace}/{name} - Update
DELETE /api/v1/{resource}/{namespace}/{name} - Delete
```

### Response Formats
```json
Success: { "data": [...] } or resource object
Error: { "error": "descriptive message" }
```

### Versioning
- API versioned via URL path (`/api/v1/`)
- Breaking changes require new version
- Deprecation notices for old versions

## Technology Decisions

### Why Go for Backend?
- Native Kubernetes and Istio client libraries
- Excellent concurrency support
- Fast compilation and execution
- Small binary size
- Strong typing and error handling

### Why React for Frontend?
- Component-based architecture
- Large ecosystem and community
- Material-UI for consistent design
- Excellent developer experience
- Wide industry adoption

### Why Gin Framework?
- High performance
- Simple and intuitive API
- Rich middleware ecosystem
- Good documentation
- Active maintenance

## Future Enhancements

### Phase 2 Features
- [ ] Multi-cluster support
- [ ] Advanced traffic analytics
- [ ] GitOps integration
- [ ] WebSocket for real-time updates
- [ ] Advanced RBAC with user management
- [ ] Integration with observability tools

### Phase 3 Features
- [ ] Service mesh health scoring
- [ ] Automated canary analysis
- [ ] Policy recommendations
- [ ] Disaster recovery features
- [ ] Compliance reporting

## Development Workflow

```
Developer → Code Changes → Local Testing
→ Build (go build / npm build)
→ Docker Build → Push to Registry
→ Update K8s Manifests → kubectl apply
→ Verify Deployment → Monitoring
```

## CI/CD Pipeline (Recommended)

```
git push → GitHub Actions
→ Run Tests → Build Backend → Build Frontend
→ Build Docker Images → Push to Registry
→ Deploy to Staging → Integration Tests
→ Manual Approval → Deploy to Production
```
