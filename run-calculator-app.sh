#!/bin/sh

# Force the build to be strictly sequential at the ENGINE level
export COMPOSE_PARALLEL_LIMIT=1

# Build the images one by one manually (On my macOS, memory allocation problem was happening)
echo "Starting sequential build..."
docker compose build addition-service
docker compose build subtraction-service
docker compose build multiplication-service
docker compose build division-service
docker compose build exponentiation-service
docker compose build squareroot-service
docker compose build percentage-service
docker compose build history-service
docker compose build calculator-backend
docker compose build calculator-frontend

# Start
echo "Starting the whole app..."
docker compose up
