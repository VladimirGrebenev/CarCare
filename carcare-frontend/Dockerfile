# --- Frontend Dockerfile (SvelteKit) ---
# Build stage
FROM node:20.11.0-alpine3.19 AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Production stage
FROM gcr.io/distroless/nodejs20-debian12:nonroot
WORKDIR /app
COPY --from=builder /app/build ./build
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json
USER nonroot
EXPOSE 3000
HEALTHCHECK --interval=30s --timeout=3s --retries=3 --start-period=40s CMD ["node", "build"]
CMD ["node", "build"]
