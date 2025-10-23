#!/usr/bin/env bash
set -e

# Check if go is installed
if ! command -v go &> /dev/null; then
    echo "Error: go is not installed."
    exit 1
fi

# Detect if Graphite CLI is available
GT_AVAILABLE=0
if command -v gt &> /dev/null; then
    GT_AVAILABLE=1
fi

# Get the repository root
REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

# Fetch latest branches
echo "Fetching branches..."
git fetch --all --prune

# Decide whether to use Graphite helpers
use_graphite=0
if [ "$GT_AVAILABLE" -eq 1 ]; then
    echo ""
    if go tool gum confirm --default "Use Graphite CLI while switching?"; then
        use_graphite=1
    fi
fi

# Get list of branches (remote and local, cleaned up and deduplicated)
branches=$(git branch -a | \
    sed 's/^\*//' | \
    sed 's/^[[:space:]]*//' | \
    sed 's/remotes\/origin\///' | \
    grep -v 'HEAD ->' | \
    sort -u)

# Use gum filter for fuzzy branch selection
selected_branch=$(echo "$branches" | go tool gum filter --placeholder "Select a branch..." --height 20)

if [ -z "$selected_branch" ]; then
    echo "No branch selected. Exiting."
    exit 0
fi

echo "Selected: $selected_branch"

# Create worktree directory name (sanitize branch name for filesystem)
worktree_name=$(echo "$selected_branch" | sed 's/\//-/g')
worktree_path=$(cd "$REPO_ROOT/.." && pwd)/"$worktree_name"

starting_branch=$(git branch --show-current)

if [ "$use_graphite" -eq 1 ]; then
    echo ""
    echo "Syncing branch via Graphite..."
    set +e
    if ! gt get "$selected_branch" >/dev/null 2>&1; then
        echo "Graphite sync failed or branch not tracked yet; continuing with git data."
    fi
    set -e

    echo "Ensuring Graphite metadata for '$selected_branch'..."
    set +e
    if gt checkout "$selected_branch" >/dev/null 2>&1; then
        if [ -n "$starting_branch" ] && [ "$starting_branch" != "$selected_branch" ]; then
            gt checkout "$starting_branch" >/dev/null 2>&1 || git checkout "$starting_branch"
        fi
    else
        echo "Graphite checkout failed; branch may be untracked. You can run 'gt track $selected_branch' later."
    fi
    set -e
fi

# Check if worktree already exists
if [ -d "$worktree_path" ]; then
    echo "Worktree already exists at $worktree_path"
    
    # Ask if user wants to open it anyway
    if go tool gum confirm "Open existing worktree?"; then
        # Skip to IDE selection
        :
    else
        echo "Exiting."
        exit 0
    fi
else
    # Create the worktree
    echo "Creating worktree at $worktree_path for branch $selected_branch..."
    git worktree add "$worktree_path" "$selected_branch"
    echo "✓ Worktree created successfully!"

    if [ "$use_graphite" -eq 0 ] && [ "$GT_AVAILABLE" -eq 1 ]; then
        echo ""
        echo "Tracking branch with Graphite metadata..."
        set +e
        if ! gt track "$selected_branch" >/dev/null 2>&1; then
            echo "Graphite tracking failed (branch may already be tracked)."
            echo "Run 'gt track $selected_branch' manually if you want Graphite visibility."
        else
            echo "✓ Branch registered with Graphite."
        fi
        set -e
    fi
fi

# Build IDE options dynamically
ide_options=()

# Check for Claude Code in Wezterm
if command -v claude &> /dev/null && command -v wezterm &> /dev/null; then
    ide_options+=("Claude Code in Wezterm")
fi

# Check for Codex in Wezterm
if command -v codex &> /dev/null && command -v wezterm &> /dev/null; then
    ide_options+=("Codex in Wezterm")
fi

# Standard Cursor/VSCode options
if command -v cursor &> /dev/null; then
    ide_options+=("Cursor")
fi
if command -v code &> /dev/null; then
    ide_options+=("VSCode")
fi

# Always provide option to skip
ide_options+=("Skip (don't open)")

# Prompt for IDE choice
if [ ${#ide_options[@]} -eq 1 ]; then
    # Only "Skip" option available
    echo "No IDE detected. Worktree is ready at: $worktree_path"
    exit 0
fi

selected_ide=$(printf '%s\n' "${ide_options[@]}" | go tool gum filter --placeholder "Open with..." --height 10)

if [ -z "$selected_ide" ]; then
    echo "No IDE selected. Worktree is ready at: $worktree_path"
    exit 0
fi

# Launch the selected IDE
case "$selected_ide" in
    "Claude Code in Wezterm")
        echo "Launching Claude Code in Wezterm..."
        wezterm start --cwd "$worktree_path" -- claude "$worktree_path" &
        ;;
    "Codex in Wezterm")
        echo "Launching Codex in Wezterm..."
        wezterm start --cwd "$worktree_path" -- codex "$worktree_path" &
        ;;
    "Cursor")
        echo "Opening in Cursor..."
        cursor "$worktree_path" &
        ;;
    "VSCode")
        echo "Opening in VSCode..."
        code "$worktree_path" &
        ;;
    "Skip (don't open)")
        echo "Worktree is ready at: $worktree_path"
        ;;
esac

echo "✓ Done!"
