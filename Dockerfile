# Multi-stage build

# Stage 1: Build frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build backend
FROM golang:1.22-alpine AS backend-builder
WORKDIR /app/backend
RUN apk add --no-cache git gcc musl-dev
COPY backend/go.* ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o meshcontrol-server cmd/server/main.go

# Stage 3: Final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite
WORKDIR /root/
COPY --from=backend-builder /app/backend/meshcontrol-server .
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

ENV STATIC_DIR=/root/frontend/dist
ENV PORT=8080
ENV DB_PATH=/data/meshcontrol.db

EXPOSE 8080

# Create data directory
RUN mkdir -p /data

CMD ["./meshcontrol-server"]
