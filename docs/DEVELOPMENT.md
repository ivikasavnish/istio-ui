# MeshControl Center - Development Guide

## Project Overview

MeshControl Center is a full-stack application for managing Istio service mesh configurations. It consists of:

- **Backend**: Go 1.22+ REST API server
- **Frontend**: React 18 + TypeScript + Vite application
- **Database**: SQLite (production can use PostgreSQL)
- **Deployment**: Kubernetes + Helm

## Development Setup

### Backend Setup

1. Install Go 1.22 or higher
2. Navigate to backend directory:
```bash
cd backend
```

3. Install dependencies:
```bash
go mod download
```

4. Run the server:
```bash
go run cmd/server/main.go
```

The server will start on port 8080.

### Frontend Setup

1. Install Node.js 18 or higher
2. Navigate to frontend directory:
```bash
cd frontend
```

3. Install dependencies:
```bash
npm install
```

4. Run development server:
```bash
npm run dev
```

The frontend will start on port 3000 with hot reload.

## Code Structure

### Backend Architecture

```
backend/
├── cmd/server/         # Application entry point
├── internal/
│   ├── api/           # HTTP handlers & WebSocket
│   ├── istio/         # Istio CRD operations
│   ├── kube/          # Kubernetes client
│   ├── models/        # Data models
│   ├── scheduler/     # Cron scheduler
│   ├── storage/       # Database layer
│   └── middleware/    # Auth & logging
```

### Frontend Architecture

```
frontend/src/
├── pages/             # Page components (routes)
├── components/        # Reusable UI components
├── services/          # API client
├── hooks/            # Custom React hooks
└── utils/            # Utility functions
```

## Adding New Features

### Adding a New Istio Resource

1. **Add model** in `backend/internal/models/models.go`
2. **Add manager methods** in `backend/internal/istio/manager.go`
3. **Add API handlers** in `backend/internal/api/handlers.go`
4. **Add routes** in `backend/internal/api/server.go`
5. **Add frontend API client** in `frontend/src/services/api.ts`
6. **Create UI page** in `frontend/src/pages/`

### Adding a New Page

1. Create component in `frontend/src/pages/NewPage.tsx`
2. Add route in `frontend/src/App.tsx`
3. Add navigation link in `frontend/src/components/Layout.tsx`

## Testing

### Backend Tests

```bash
cd backend
go test ./...
```

### Frontend Tests

```bash
cd frontend
npm test
```

## Building

### Build Backend Binary

```bash
cd backend
go build -o meshcontrol-server cmd/server/main.go
```

### Build Frontend for Production

```bash
cd frontend
npm run build
```

Output in `frontend/dist/`

### Build Docker Image

From project root:
```bash
docker build -t meshcontrol:latest .
```

## Deployment

### Local Kubernetes

```bash
kubectl apply -f deploy/k8s/deployment.yaml
kubectl port-forward -n meshcontrol svc/meshcontrol 8080:80
```

### Using Helm

```bash
helm install meshcontrol deploy/helm/meshcontrol \
  --namespace meshcontrol \
  --create-namespace
```

## Environment Variables

### Backend

- `PORT` - Server port (default: 8080)
- `DB_PATH` - Database path (default: ./meshcontrol.db)
- `KUBECONFIG` - Kubernetes config (optional)
- `STATIC_DIR` - Frontend files directory

### Frontend (Dev)

Set in `frontend/.env`:
```
VITE_API_URL=http://localhost:8080
```

## Database Migrations

Database schema is initialized automatically on first run. To reset:

```bash
rm backend/meshcontrol.db
```

## Troubleshooting

### Backend Issues

1. **Cannot connect to Kubernetes**
   - Check KUBECONFIG: `echo $KUBECONFIG`
   - Test connection: `kubectl cluster-info`

2. **Database locked**
   - Ensure only one instance is running
   - Check file permissions

### Frontend Issues

1. **Cannot connect to backend**
   - Verify backend is running on port 8080
   - Check browser console for CORS errors
   - Verify proxy settings in `vite.config.ts`

2. **Build fails**
   - Clear node_modules: `rm -rf node_modules && npm install`
   - Clear cache: `rm -rf .vite`

## Best Practices

### Backend

- Use context for cancellation
- Always handle errors
- Log important operations
- Validate input data
- Use transactions for multi-step operations

### Frontend

- Use TypeScript types
- Handle loading states
- Show user feedback (alerts/toasts)
- Validate form inputs
- Clean up effects and subscriptions

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## Resources

- [Istio Documentation](https://istio.io/docs/)
- [Kubernetes Client Go](https://github.com/kubernetes/client-go)
- [React Documentation](https://react.dev/)
- [Go Best Practices](https://golang.org/doc/effective_go)
