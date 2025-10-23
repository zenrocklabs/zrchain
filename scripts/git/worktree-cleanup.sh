#!/usr/bin/env bash
set -e

# Check if go is installed
if ! command -v go &> /dev/null; then
    echo "Error: go is not installed."
    exit 1
fi

# Get the repository root
REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

# Get list of worktrees, excluding the main one
worktree_list=$(git worktree list --porcelain | awk '
    BEGIN { path=""; branch="" }
    /^worktree / { path=$2 }
    /^branch / { branch=$2; gsub(/^refs\/heads\//, "", branch) }
    /^$/ { 
        if (path != "" && branch != "") {
            print path "|" branch
        }
        path=""; branch=""
    }
    END {
        if (path != "" && branch != "") {
            print path "|" branch
        }
    }
' | grep -v "^$REPO_ROOT|")

# Check if there are any worktrees to clean up
if [ -z "$worktree_list" ]; then
    echo "No additional worktrees found to clean up."
    exit 0
fi

# Format worktrees for display (show path and branch)
display_list=$(echo "$worktree_list" | while IFS='|' read -r path branch; do
    basename_path=$(basename "$path")
    echo "$basename_path ($branch)"
done)

# Let user select multiple worktrees to delete
echo "Select worktrees to delete (use tab/space to select, enter to confirm):"
selected=$(echo "$display_list" | go tool gum choose --no-limit --height 15)

if [ -z "$selected" ]; then
    echo "No worktrees selected. Exiting."
    exit 0
fi

echo ""
echo "Selected worktrees for deletion:"
echo "$selected"
echo ""

# Confirm before proceeding
if ! go tool gum confirm "Proceed with deletion?"; then
    echo "Cancelled."
    exit 0
fi

# Process each selected worktree
echo "$selected" | while read -r selection; do
    # Extract the basename from the selection (format: "basename (branch)")
    basename_selected=$(echo "$selection" | sed 's/ (.*//')
    
    # Find the full path and branch for this worktree
    worktree_info=$(echo "$worktree_list" | grep -F "/$basename_selected|" || echo "")
    
    if [ -z "$worktree_info" ]; then
        echo "⚠️  Could not find worktree info for: $selection"
        continue
    fi
    
    IFS='|' read -r worktree_path branch_name <<< "$worktree_info"
    
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Processing: $basename_selected"
    echo "  Path: $worktree_path"
    echo "  Branch: $branch_name"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    # Remove the worktree
    echo "Removing worktree..."
    if git worktree remove "$worktree_path" 2>/dev/null; then
        echo "✓ Worktree removed: $worktree_path"
    elif git worktree remove --force "$worktree_path" 2>/dev/null; then
        echo "✓ Worktree force-removed: $worktree_path"
    else
        echo "⚠️  Failed to remove worktree: $worktree_path"
        continue
    fi
    
    # Ask if user wants to delete the branch too
    if go tool gum confirm "Also delete branch '$branch_name'?"; then
        # Check if branch exists locally
        if git show-ref --verify --quiet "refs/heads/$branch_name"; then
            if git branch -d "$branch_name" 2>/dev/null; then
                echo "✓ Branch deleted: $branch_name"
            else
                echo "⚠️  Branch has unmerged changes. Force delete?"
                if go tool gum confirm "Force delete branch '$branch_name'?"; then
                    git branch -D "$branch_name"
                    echo "✓ Branch force-deleted: $branch_name"
                else
                    echo "→ Branch kept: $branch_name"
                fi
            fi
        else
            echo "→ Branch doesn't exist locally (may be remote-only): $branch_name"
        fi
    else
        echo "→ Branch kept: $branch_name"
    fi
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✓ Cleanup complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Prune any stale worktree references
git worktree prune -v

