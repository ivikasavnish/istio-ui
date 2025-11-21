# Getting Started with MeshControl Center

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.21+**: Download from [golang.org](https://golang.org/dl/)
- **Node.js 18+**: Download from [nodejs.org](https://nodejs.org/)
- **kubectl**: Kubernetes command-line tool
- **Access to a Kubernetes cluster** with Istio installed

## Installation

### Option 1: Running Locally

#### 1. Clone the Repository

```bash
git clone https://github.com/ivikasavnish/istio-ui.git
cd istio-ui
```

#### 2. Start the Backend

```bash
cd backend
go mod download
go run cmd/server/main.go
```

The backend API server will start on `http://localhost:8080`

#### 3. Start the Frontend (in a new terminal)

```bash
cd frontend
npm install
npm start
```

The React application will open in your browser at `http://localhost:3000`

### Option 2: Using Docker Compose

```bash
docker-compose up --build
```

This will start both the backend and frontend services. Access the application at `http://localhost:3000`

### Option 3: Deploying to Kubernetes

#### 1. Build Docker Images

```bash
# Build backend
cd backend
docker build -t meshcontrol-backend:latest .

# Build frontend
cd ../frontend
docker build -t meshcontrol-frontend:latest .
```

#### 2. Deploy to Kubernetes

```bash
kubectl apply -f deployment/kubernetes.yaml
```

#### 3. Access the Application

Port-forward to access locally:
```bash
kubectl port-forward -n meshcontrol-center svc/meshcontrol-frontend 3000:80
```

Or configure the Ingress with your domain in `deployment/kubernetes.yaml`

## Configuration

### Backend Environment Variables

- `KUBECONFIG`: Path to kubeconfig file (default: `~/.kube/config`)
- `PORT`: Server port (default: `8080`)
- `ENABLE_CORS`: Enable CORS (default: `true`)
- `LOG_LEVEL`: Logging level (default: `info`)

### Frontend Environment Variables

- `REACT_APP_API_URL`: Backend API URL (default: `http://localhost:8080`)

## First Steps

### 1. Check the Dashboard

Navigate to the dashboard to see an overview of your Istio resources:
- Virtual Services count
- Destination Rules count
- Gateways count
- Authorization Policies count

### 2. Create a Virtual Service

1. Navigate to "Virtual Services" from the sidebar
2. Click "Create Virtual Service"
3. Fill in the form:
   - Name: `my-service-vs`
   - Namespace: `default`
   - Host: `my-service.example.com`
   - Destination Host: `my-service.default.svc.cluster.local`
4. Click "Save"

### 3. View YAML

Click the code icon next to any resource to view its YAML configuration.

### 4. Create a Scheduled Action

1. Navigate to "Scheduled Actions"
2. Click "Create Schedule"
3. Configure:
   - Name: `Traffic Shift`
   - Schedule: `0 */6 * * *` (every 6 hours)
   - Action Type: Update VirtualService
   - Resource Name: Your VirtualService name

### 5. Explore the Topology

Navigate to "Topology" to see a visual representation of your service mesh.

## Troubleshooting

### Backend won't start

**Issue**: Cannot connect to Kubernetes cluster

**Solution**: Ensure your kubeconfig is correctly configured:
```bash
kubectl cluster-info
```

### Frontend shows API errors

**Issue**: Cannot connect to backend

**Solution**: Verify the backend is running and CORS is enabled:
```bash
curl http://localhost:8080/health
```

### Resources not appearing

**Issue**: No resources shown in the UI

**Solution**: Check namespace access and RBAC permissions:
```bash
kubectl get virtualservices --all-namespaces
```

## Next Steps

- Read the [API Documentation](API.md) for programmatic access
- Check out [examples](../examples/) for common Istio configurations
- Learn about [Scheduled Actions](SCHEDULED_ACTIONS.md) for automation

## Support

For issues and questions:
- GitHub Issues: https://github.com/ivikasavnish/istio-ui/issues
- Documentation: https://github.com/ivikasavnish/istio-ui/docs
