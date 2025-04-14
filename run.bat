@echo off
echo Checking if Docker is running...

docker info > nul 2>&1
if errorlevel 1 (
    echo Docker is not running. Please start Docker Desktop and try again.
    exit /b 1
)

echo Building and starting Docker Container Health Checker...

docker-compose down
docker-compose up --build

echo Docker Container Health Checker has stopped. 