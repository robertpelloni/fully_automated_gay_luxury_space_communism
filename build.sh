#!/bin/bash
echo "Building Orchestrator..."
cd orchestrator && go build ./... && cd ..
echo "Building Research Hustle..."
cd hustle/research && go build ./... && cd ..
echo "Build complete."
