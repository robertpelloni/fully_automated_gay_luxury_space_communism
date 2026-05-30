#!/bin/bash
echo "Building AI Hustle Machine..."

echo "Building Orchestrator..."
cd orchestrator && go build -o ../bin/orchestrator ./cmd/orchestrator && cd ..

echo "Building Research Hustle..."
cd hustle/research && go build -o ../../bin/research ./cmd/research && cd ..

echo "Build complete. Executables in bin/"
