#!/bin/bash
set -e

echo "=== EXECUTIVE PROTOCOL: REPOSITORY SYNC ==="

# Failure handler
trap 'echo "ERROR: Sync failed. Aborting merge and restoring state..."; git merge --abort 2>/dev/null || true; exit 1' ERR

# 1. Fetch all
echo "Fetching all remotes and tags..."
git fetch --all --tags

# 2. Submodule Update
echo "Updating submodules recursively..."
git submodule update --init --recursive

# 3. Identify and Reconciliation
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
echo "Current branch: $CURRENT_BRANCH"

# Forward Merge: Features to Main
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo "Attempting to merge main into $CURRENT_BRANCH (Reverse Merge)..."
    git merge origin/main --no-edit
fi

echo "Reconciliation complete."
echo "Sync successful."
