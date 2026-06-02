#!/bin/bash
# THE EXECUTIVE PROTOCOL: REPOSITORY SYNCHRONIZATION & INTELLIGENT MERGE
# This script maintains repository health by reconciling all local feature branches with main.

echo "=== EXECUTIVE PROTOCOL: REPOSITORY SYNC & INTELLIGENT MERGE ==="

# 1. UPSTREAM TRACKING & SUBMODULE SANITIZATION
echo "Step 1: Upstream Tracking & Submodule Sanitization"

# Identify the upstream remote (prefer 'upstream' if it exists, otherwise 'origin')
REMOTE="origin"
if git remote | grep -q "^upstream$"; then
    REMOTE="upstream"
fi
echo "Using remote: $REMOTE"

# Detect main branch name (main or master)
MAIN_BRANCH="main"
if git branch -r | grep -q "$REMOTE/master"; then
    MAIN_BRANCH="master"
fi
echo "Main branch detected: $MAIN_BRANCH"

echo "Fetching all remotes and tags..."
git fetch --all --tags

# Ensure local main is updated
echo "Updating local $MAIN_BRANCH from $REMOTE..."
git checkout $MAIN_BRANCH
git merge $REMOTE/$MAIN_BRANCH --ff-only || git reset --hard $REMOTE/$MAIN_BRANCH


echo "Updating submodules recursively to their latest tracking commits..."
git submodule update --init --recursive

# Identify current branch to return to it later
INITIAL_BRANCH=$(git rev-parse --abbrev-ref HEAD)
echo "Initial branch: $INITIAL_BRANCH"

# 2. DUAL-DIRECTION INTELLIGENT MERGE ENGINE
echo "Step 2: Dual-Direction Intelligent Merge Engine"

# Get all local branches
ALL_LOCAL_BRANCHES=$(git branch --format='%(refname:short)')

for BRANCH in $ALL_LOCAL_BRANCHES; do
    if [ "$BRANCH" == "main" ] || [ "$BRANCH" == "master" ]; then
        continue
    fi

    echo "--- Reconciling Branch: $BRANCH ---"

    # Stash if necessary
    STASHED=false
    if [[ $(git status --porcelain) ]]; then
        echo "Stashing local changes..."
        git stash
        STASHED=true
    fi

    # Checkout the feature branch
    if ! git checkout "$BRANCH"; then
        echo "Failed to checkout $BRANCH, skipping."
        continue
    fi

    # Reverse Merge: Merging main into feature branch
    echo "Reverse Merge: Merging $REMOTE/$MAIN_BRANCH into $BRANCH..."
    if git merge "$REMOTE/$MAIN_BRANCH" --no-edit; then
        echo "Successfully caught up $BRANCH with $REMOTE/$MAIN_BRANCH."
    else
        echo "CONFLICT detected on $BRANCH. Aborting reverse merge."
        git merge --abort
    fi

    # Forward Merge: Merging feature branch back into main (Optional/Automated)
    # Directive: "Interrogate each active feature branch. If it contains unique development progress... merge it into main."
    # For safety in this script, we only forward merge branches with the 'feat/ready-' prefix.
    if [[ $BRANCH == feat/ready-* ]]; then
        echo "Forward Merge: Merging $BRANCH back into $MAIN_BRANCH..."
        git checkout $MAIN_BRANCH
        if git merge "$BRANCH" --no-edit; then
            echo "Successfully merged $BRANCH into $MAIN_BRANCH."
            git checkout "$BRANCH" # Go back to feature branch for final state
        else
            echo "CONFLICT detected during forward merge of $BRANCH. Aborting."
            git merge --abort
            git checkout "$BRANCH"
        fi
    fi

    # Return to initial branch
    git checkout "$INITIAL_BRANCH"

    if [ "$STASHED" = true ]; then
        echo "Restoring stashed changes..."
        git stash pop
        STASHED=false
    fi
done

# Ensure we are back on the initial branch and it is caught up
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "$MAIN_BRANCH" ]; then
    echo "Finalizing sync for $CURRENT_BRANCH..."
    git merge "$REMOTE/$MAIN_BRANCH" --no-edit || echo "Final sync merge failed. Manual intervention required."
fi

echo "Step 3: Workspace Cleanup, Documentation, & Push"

# Governance: Source version
VERSION=$(cat VERSION.md)
echo "Current Version: $VERSION"

# Batch Script Validation: Ensure execution scripts use the correct MAIN_BRANCH
for script in build.sh start.sh; do
    if [ -f "$script" ]; then
        echo "Validating $script..."
        # Targeted replacement using word boundaries where possible to avoid collateral damage
        # Portable sed: use -i '' for BSD/macOS compatibility if needed, but here we detect
        if sed --version >/dev/null 2>&1; then
            sed -i "s/\bmain\b/$MAIN_BRANCH/g" "$script"
        else
            sed -i "" "s/[[:<:]]main[[:>:]]/$MAIN_BRANCH/g" "$script"
        fi
    fi
done

# Synchronize version across changelog and status
if [ -f "CHANGELOG.md" ]; then
    # Match both [1.0.0-alpha.x] and v1.0.0-alpha.x formats
    if sed --version >/dev/null 2>&1; then
        sed -i "s/1.0.0-alpha.[0-9]*/$VERSION/g" CHANGELOG.md
    else
        sed -i "" "s/1.0.0-alpha.[0-9]*/$VERSION/g" CHANGELOG.md
    fi
fi

# Stage all changes
git add .

# Atomic Commit if changes exist
if ! git diff-index --quiet HEAD --; then
    echo "Executing atomic commit for version $VERSION..."
    git commit -m "build: $VERSION - Executive Protocol Sync and feature updates"
else
    echo "No changes to commit."
fi

# Push to server
if [ "$HUSTLE_PUSH_ENABLED" = "true" ]; then
    echo "Pushing changes to $REMOTE..."
    git push $REMOTE $CURRENT_BRANCH
    git push $REMOTE --tags
else
    echo "Push disabled. Set HUSTLE_PUSH_ENABLED=true to enable automatic remote updates."
fi

echo "=== EXECUTIVE PROTOCOL COMPLETE ==="
