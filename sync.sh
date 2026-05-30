#!/bin/bash
echo "Fetching all remotes and tags..."
git fetch --all --tags
echo "Merging main into current branch..."
git merge origin/main
echo "Updating submodules recursively..."
git submodule update --init --recursive
echo "Sync complete."
