# --- Backend Dockerfile (Go) ---
# Build stage
FROM golang:1.22.2-alpine3.19 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o carcare ./cmd/carcare

# Production stage
FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache ca-certificates curl
COPY --from=builder /app/carcare ./carcare
COPY --from=builder /app/deployments/docker-entrypoint.sh ./entrypoint.sh
COPY --from=builder /app/internal/infrastructure/migration ./migration
RUN chmod +x ./entrypoint.sh
USER 1000:1000
EXPOSE 8080
ENTRYPOINT ["./entrypoint.sh"]
CMD ["./carcare"]
