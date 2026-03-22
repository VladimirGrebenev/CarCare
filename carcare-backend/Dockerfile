# --- Backend Dockerfile (Go) ---
# Build stage
FROM golang:1.22.2-alpine3.19 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o carcare ./cmd/carcare

# Production stage
FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=builder /app/carcare ./carcare
COPY --from=builder /app/deployments/docker-entrypoint.sh ./entrypoint.sh
COPY --from=builder /app/internal/infrastructure/migration ./migration
ADD https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz /tmp/migrate.tar.gz
RUN tar -xzf /tmp/migrate.tar.gz -C /usr/local/bin && rm /tmp/migrate.tar.gz
USER nonroot
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --retries=3 --start-period=40s CMD ["./carcare", "healthcheck"]
ENTRYPOINT ["./entrypoint.sh"]
