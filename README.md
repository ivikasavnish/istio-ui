# MeshControl Center

![MeshControl Center](https://img.shields.io/badge/version-1.0.0-blue.svg)
![Go](https://img.shields.io/badge/Go-1.22+-00ADD8.svg)
![React](https://img.shields.io/badge/React-18.2-61DAFB.svg)
![Istio](https://img.shields.io/badge/Istio-1.20+-466BB0.svg)

A comprehensive, production-ready visual management tool for Istio service mesh. MeshControl Center provides a modern UI for managing all Istio traffic management, security, and routing features with real-time monitoring, scheduled actions, and YAML preview capabilities.

## 🚀 Features

### Traffic Management
- **Weighted Routing**: Visual sliders for traffic split between service versions (v1/v2/v3)
- **Canary Deployments**: Progressive rollout with percentage-based traffic shifting
- **Blue-Green Deployments**: Instant traffic switching between versions
- **Traffic Mirroring**: Shadow traffic for testing
- **Circuit Breaking**: Configure connection pools and outlier detection

### Security
- **mTLS Configuration**: Toggle between STRICT, PERMISSIVE, and DISABLE modes
- **PeerAuthentication**: Manage service-to-service authentication
- **AuthorizationPolicy**: RBAC for service access control
- **JWT Authentication**: Token-based authentication support

### Fault Injection
- **Delay Injection**: Add latency to test resilience
- **Abort Injection**: Simulate HTTP errors (400, 500, 503, etc.)
- **gRPC Errors**: Test gRPC failure scenarios
- **Percentage-based**: Control what percentage of traffic is affected

### Observability
- **Service Discovery**: Auto-discover services across namespaces
- **Topology Graph**: Interactive visualization of service mesh (Cytoscape.js)
- **Real-time Events**: WebSocket-based live updates
- **Audit Logs**: Complete history of all changes

### Automation
- **Scheduler**: Cron-based scheduled actions
- **Traffic Snapshots**: Capture and restore configurations
- **Batch Operations**: Apply changes to multiple services
- **Rollback System**: Restore previous configurations

### Developer Experience
- **YAML Preview**: Live YAML generation from forms
- **Dark Mode**: Full dark theme support
- **REST API**: Complete programmatic access
- **WebSocket API**: Real-time event streaming

## 📦 Architecture

```
┌─────────────────────────────────────────────────────┐
│                  Frontend (React)                    │
│  ┌────────────┬──────────────┬──────────────────┐  │
│  │ Dashboard  │  Traffic Mgmt │   Security       │  │
│  ├────────────┼──────────────┼──────────────────┤  │
│  │  Topology  │  YAML Editor  │   Scheduler      │  │
│  └────────────┴──────────────┴──────────────────┘  │
└───────────────────────┬─────────────────────────────┘
                        │ REST + WebSocket
┌───────────────────────┴─────────────────────────────┐
│               Backend (Go)                           │
│  ┌────────────┬──────────────┬──────────────────┐  │
│  │  API Layer │  Istio Mgr   │  K8s Client      │  │
│  ├────────────┼──────────────┼──────────────────┤  │
│  │  Scheduler │   Storage    │  WebSocket       │  │
│  └────────────┴──────────────┴──────────────────┘  │
└───────────────────────┬─────────────────────────────┘
                        │
┌───────────────────────┴─────────────────────────────┐
│           Kubernetes + Istio                         │
│  ┌─────────────────────────────────────────────┐   │
│  │  VirtualServices  │  DestinationRules       │   │
│  │  Gateways         │  PeerAuthentication     │   │
│  │  AuthzPolicies    │  Services               │   │
│  └─────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────┘
```

## 🏗️ Project Structure

```
meshcontrol/
├── backend/                    # Go backend
│   ├── cmd/
│   │   └── server/
│   │       └── main.go        # Entry point
│   ├── internal/
│   │   ├── api/               # REST API handlers
│   │   ├── istio/             # Istio CRD management
│   │   ├── kube/              # Kubernetes client
│   │   ├── scheduler/         # Cron scheduler
│   │   ├── storage/           # Database layer
│   │   ├── models/            # Data models
│   │   └── middleware/        # Auth & RBAC
│   └── go.mod
├── frontend/                  # React frontend
│   ├── src/
│   │   ├── pages/            # Page components
│   │   ├── components/       # Reusable components
│   │   ├── services/         # API clients
│   │   └── hooks/            # Custom hooks
│   └── package.json
├── deploy/                   # Deployment configs
│   ├── k8s/
│   │   └── deployment.yaml
│   └── helm/
│       └── meshcontrol/
└── README.md
```

## 🚦 Getting Started

### Prerequisites

- **Kubernetes cluster** (v1.24+)
- **Istio** installed (v1.18+)
- **Go** 1.22+ (for backend development)
- **Node.js** 18+ (for frontend development)
- **Docker** (for building containers)
- **kubectl** configured

### Quick Start with Kubernetes

1. **Clone the repository:**
```bash
git clone https://github.com/ivikasavnish/istio-ui.git
cd istio-ui
```

2. **Deploy using kubectl:**
```bash
kubectl apply -f deploy/k8s/deployment.yaml
```

3. **Access the UI:**
```bash
kubectl port-forward -n meshcontrol svc/meshcontrol 8080:80
```

Open http://localhost:8080 in your browser.

### Quick Start with Helm

1. **Install using Helm:**
```bash
helm install meshcontrol deploy/helm/meshcontrol \
  --create-namespace \
  --namespace meshcontrol
```

2. **Access the UI:**
```bash
kubectl port-forward -n meshcontrol svc/meshcontrol 8080:80
```

Open http://localhost:8080 in your browser.

## 💻 Local Development

### Backend Development

1. **Navigate to backend directory:**
```bash
cd backend
```

2. **Install dependencies:**
```bash
go mod download
```

3. **Run the backend:**
```bash
go run cmd/server/main.go
```

The API server will start on http://localhost:8080

### Frontend Development

1. **Navigate to frontend directory:**
```bash
cd frontend
```

2. **Install dependencies:**
```bash
npm install
```

3. **Run the development server:**
```bash
npm run dev
```

The frontend will start on http://localhost:3000 with hot-reload enabled.

### Building for Production

#### Build Backend
```bash
cd backend
go build -o meshcontrol-server cmd/server/main.go
```

#### Build Frontend
```bash
cd frontend
npm run build
```

The production build will be in `frontend/dist/`

## 📖 API Documentation

### REST API Endpoints

#### Services
- `GET /api/v1/services` - List all services
- `GET /api/v1/services/{namespace}` - List services in namespace
- `GET /api/v1/namespaces` - List all namespaces

#### VirtualServices
- `GET /api/v1/virtualservices/{namespace}` - List VirtualServices
- `GET /api/v1/virtualservices/{namespace}/{name}` - Get VirtualService
- `POST /api/v1/virtualservices` - Create VirtualService
- `PUT /api/v1/virtualservices/{namespace}/{name}` - Update VirtualService
- `DELETE /api/v1/virtualservices/{namespace}/{name}` - Delete VirtualService

#### DestinationRules
- `GET /api/v1/destinationrules/{namespace}` - List DestinationRules
- `GET /api/v1/destinationrules/{namespace}/{name}` - Get DestinationRule
- `POST /api/v1/destinationrules` - Create DestinationRule
- `PUT /api/v1/destinationrules/{namespace}/{name}` - Update DestinationRule
- `DELETE /api/v1/destinationrules/{namespace}/{name}` - Delete DestinationRule

#### Traffic Management
- `POST /api/v1/traffic/weights` - Update traffic weights
- `POST /api/v1/traffic/canary` - Configure canary deployment

#### Snapshots & Scheduled Actions
- See full API documentation in `/docs/API.md`

### WebSocket API

Connect to `ws://localhost:8080/ws` for real-time events.

## 🔧 Configuration

### Environment Variables

- `PORT` - API server port (default: 8080)
- `DB_PATH` - SQLite database path (default: ./meshcontrol.db)
- `KUBECONFIG` - Kubernetes config path (auto-detected)
- `STATIC_DIR` - Frontend static files directory

## 🎯 Use Cases

### Canary Deployment

1. Navigate to **Traffic Management**
2. Select your service and namespace
3. Adjust sliders: v1=90%, v2=10%
4. Click "Create VirtualService"
5. Monitor traffic in **Dashboard**
6. Gradually increase v2 percentage

### mTLS Enforcement

1. Navigate to **Security**
2. Select namespace
3. Choose "STRICT" mTLS mode
4. Click "Apply mTLS Policy"
5. Verify in **Audit Logs**

### Fault Testing

1. Navigate to **Fault Injection**
2. Select service
3. Choose "Delay" with 5000ms
4. Set percentage to 50%
5. Click "Apply Fault Injection"
6. Test service resilience

## 🛡️ Security

- **RBAC**: Kubernetes RBAC for API access
- **Authentication**: JWT-based authentication (configurable)
- **Audit Logging**: All actions logged
- **TLS**: Support for TLS/HTTPS

## 🐛 Troubleshooting

### Backend won't start
- Check Kubernetes connectivity: `kubectl cluster-info`
- Verify Istio is installed: `kubectl get ns istio-system`
- Check logs: `kubectl logs -n meshcontrol deployment/meshcontrol`

### Frontend can't connect to backend
- Verify backend is running: `curl http://localhost:8080/health`
- Check proxy settings in `vite.config.ts`

## 📝 License

MIT License

## 🤝 Contributing

Contributions welcome! Please open an issue or submit a PR.

## 📧 Support

- GitHub Issues: https://github.com/ivikasavnish/istio-ui/issues

---

Made with ❤️ by the MeshControl Team
