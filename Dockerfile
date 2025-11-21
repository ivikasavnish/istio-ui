# Backend Dockerfile
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o meshcontrol-server cmd/server/main.go

# Frontend Dockerfile
FROM node:18-alpine AS frontend-builder

WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

# Final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy backend binary
COPY --from=backend-builder /app/meshcontrol-server .

# Copy frontend build
COPY --from=frontend-builder /app/build ./frontend/build

EXPOSE 8080

CMD ["./meshcontrol-server"]
