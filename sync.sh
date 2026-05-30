#!/bin/bash
set -e

echo "=== EXECUTIVE PROTOCOL: REPOSITORY SYNC ==="

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
    git merge origin/main --no-edit || (echo "Merge conflict in reverse merge, please resolve manually" && exit 1)
fi

# Reverse Merge: Main into features (This script assumes we are on a feature branch usually)
# In a real "Executive Protocol", we might iterate through all local branches.
# For now, we ensure the current branch is caught up with main.

echo "Reconciliation complete."
echo "Sync successful."
