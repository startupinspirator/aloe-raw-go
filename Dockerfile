# Build stage: Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
RUN npm run build

# Build stage: Backend
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
RUN apk add --no-cache gcc musl-dev
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN go build -o server .

# Final stage
FROM alpine:latest
WORKDIR /app
# Copy backend binary
COPY --from=backend-builder /app/server .
# Copy frontend dist folder (Go serves this)
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

EXPOSE 8080
CMD ["./server"]
