# MeshControl Center

**A comprehensive visual management platform for Istio service mesh**

MeshControl Center is a full-stack application that provides an intuitive UI for managing all aspects of Istio traffic management, security, and routing features. It combines the power of Kubernetes and Istio with a modern React frontend and a robust Go backend.

## Features

### 🎯 Core Capabilities
- **Traffic Management**: Visual configuration of VirtualServices, DestinationRules, and Gateways
- **Security Policies**: Manage AuthorizationPolicies, PeerAuthentication, and RequestAuthentication
- **Routing Configuration**: Interactive UI for complex routing rules with weight-based traffic splitting
- **Visual Graphs**: Real-time visualization of service mesh topology and traffic flows
- **YAML Preview**: Live YAML generation and editing for all Istio resources
- **Scheduled Actions**: Time-based automation for traffic shifting, canary deployments, and more
- **Form-based Configuration**: User-friendly forms with validation for all Istio resources

### 🏗️ Architecture
- **Backend**: Go-based REST API with native Kubernetes and Istio client libraries
- **Frontend**: Modern React application with Material-UI components
- **Database**: In-memory storage with persistence options
- **Scheduler**: Built-in cron scheduler for timed actions

## Quick Start

### Prerequisites
- Go 1.21+ 
- Node.js 18+
- kubectl configured with cluster access
- Istio 1.18+ installed on your Kubernetes cluster

### Running Locally

#### Backend
```bash
cd backend
go mod download
go run cmd/server/main.go
```

The API server will start on `http://localhost:8080`

#### Frontend
```bash
cd frontend
npm install
npm start
```

The UI will be available at `http://localhost:3000`

### Using Docker

Build and run with Docker Compose:
```bash
docker-compose up --build
```

Access the application at `http://localhost:3000`

## Project Structure

```
meshcontrol-center/
├── backend/                 # Go backend service
│   ├── cmd/
│   │   └── server/         # Main application entry point
│   ├── internal/
│   │   ├── api/            # HTTP handlers and routes
│   │   ├── istio/          # Istio client and resource management
│   │   ├── scheduler/      # Scheduled actions service
│   │   └── k8s/            # Kubernetes client utilities
│   ├── pkg/                # Public libraries
│   └── go.mod
├── frontend/               # React frontend
│   ├── public/
│   ├── src/
│   │   ├── components/     # React components
│   │   ├── pages/          # Page components
│   │   ├── services/       # API client services
│   │   └── utils/          # Utility functions
│   └── package.json
├── deployment/             # Kubernetes deployment manifests
├── docs/                   # Documentation
└── docker-compose.yml
```

## API Documentation

### Endpoints

#### VirtualServices
- `GET /api/v1/virtualservices` - List all VirtualServices
- `GET /api/v1/virtualservices/:namespace/:name` - Get specific VirtualService
- `POST /api/v1/virtualservices` - Create VirtualService
- `PUT /api/v1/virtualservices/:namespace/:name` - Update VirtualService
- `DELETE /api/v1/virtualservices/:namespace/:name` - Delete VirtualService

#### DestinationRules
- `GET /api/v1/destinationrules` - List all DestinationRules
- `GET /api/v1/destinationrules/:namespace/:name` - Get specific DestinationRule
- `POST /api/v1/destinationrules` - Create DestinationRule
- `PUT /api/v1/destinationrules/:namespace/:name` - Update DestinationRule
- `DELETE /api/v1/destinationrules/:namespace/:name` - Delete DestinationRule

#### Gateways
- `GET /api/v1/gateways` - List all Gateways
- `GET /api/v1/gateways/:namespace/:name` - Get specific Gateway
- `POST /api/v1/gateways` - Create Gateway
- `PUT /api/v1/gateways/:namespace/:name` - Update Gateway
- `DELETE /api/v1/gateways/:namespace/:name` - Delete Gateway

#### Security Policies
- `GET /api/v1/authorizationpolicies` - List all AuthorizationPolicies
- `POST /api/v1/authorizationpolicies` - Create AuthorizationPolicy
- Similar endpoints for PeerAuthentication and RequestAuthentication

#### Scheduled Actions
- `GET /api/v1/schedules` - List all scheduled actions
- `POST /api/v1/schedules` - Create scheduled action
- `DELETE /api/v1/schedules/:id` - Delete scheduled action

## Configuration

### Environment Variables

#### Backend
- `KUBECONFIG` - Path to kubeconfig file (default: `~/.kube/config`)
- `PORT` - Server port (default: `8080`)
- `ENABLE_CORS` - Enable CORS (default: `true`)
- `LOG_LEVEL` - Logging level (default: `info`)

#### Frontend
- `REACT_APP_API_URL` - Backend API URL (default: `http://localhost:8080`)

## Development

### Running Tests

Backend:
```bash
cd backend
go test ./...
```

Frontend:
```bash
cd frontend
npm test
```

### Building for Production

Backend:
```bash
cd backend
go build -o meshcontrol-server cmd/server/main.go
```

Frontend:
```bash
cd frontend
npm run build
```

## Deployment

Deploy to Kubernetes:
```bash
kubectl apply -f deployment/
```

This will create:
- MeshControl Center backend deployment
- Frontend deployment served by nginx
- Service and Ingress for external access
- RBAC permissions for Istio resource management

## Security Considerations

- The application requires read/write access to Istio CRDs in your cluster
- Use RBAC to limit access to specific namespaces
- Enable authentication in production environments
- Rotate service account tokens regularly

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Support

For issues and questions:
- GitHub Issues: https://github.com/ivikasavnish/istio-ui/issues
- Documentation: https://github.com/ivikasavnish/istio-ui/docs

## Roadmap

- [ ] Multi-cluster support
- [ ] Advanced traffic analytics
- [ ] GitOps integration
- [ ] Audit logging
- [ ] Role-based access control (RBAC)
- [ ] Integration with observability tools (Prometheus, Grafana, Jaeger)
- [ ] Service mesh health scoring
- [ ] Automated canary analysis 
