#!/bin/bash
set -e

echo "Pulling latest changes..."
cd /root/carcare
git pull origin main

echo "Rebuilding and restarting containers..."
docker compose up -d --build

echo "Deploy complete!"
docker compose ps
