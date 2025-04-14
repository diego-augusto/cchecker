#!/bin/bash

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
  echo "Docker is not running. Please start Docker and try again."
  exit 1
fi

echo "Building and starting Docker Container Health Checker..."

# Build and run the Docker Compose setup
docker-compose down
docker-compose up --build

# Script will exit when docker-compose is terminated 