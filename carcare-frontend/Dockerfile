FROM node:20-alpine AS builder
WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM nginx:1.27-alpine
COPY --from=builder /app/build/ /usr/share/nginx/html/
COPY docker/nginx.conf /etc/nginx/nginx.conf
USER 101
EXPOSE 80
ENTRYPOINT ["nginx", "-g", "daemon off;"]
