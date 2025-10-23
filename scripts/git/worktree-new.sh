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

# Fetch latest branches
echo "Fetching branches..."
git fetch --all --prune

# Get current branch as default base
current_branch=$(git branch --show-current)

# Ask for new branch name
echo ""
new_branch=$(go tool gum input --placeholder "Enter new branch name (e.g., feature/my-feature)")

if [ -z "$new_branch" ]; then
    echo "No branch name provided. Exiting."
    exit 0
fi

# Validate branch name doesn't already exist
if git show-ref --verify --quiet "refs/heads/$new_branch"; then
    echo "Error: Branch '$new_branch' already exists locally."
    exit 1
fi

if git show-ref --verify --quiet "refs/remotes/origin/$new_branch"; then
    echo "Error: Branch '$new_branch' already exists on remote."
    exit 1
fi

echo ""
echo "New branch: $new_branch"

# Get list of branches to base the new branch off
branches=$(git branch -a | \
    sed 's/^\*//' | \
    sed 's/^[[:space:]]*//' | \
    sed 's/remotes\/origin\///' | \
    grep -v 'HEAD ->' | \
    sort -u)

# Use gum filter for base branch selection
echo ""
base_branch=$(echo "$branches" | go tool gum filter --placeholder "Select base branch (default: $current_branch)..." --height 20)

# If no selection, use current branch
if [ -z "$base_branch" ]; then
    base_branch="$current_branch"
    echo "Using current branch as base: $current_branch"
else
    echo "Using base branch: $base_branch"
fi

# Create worktree directory name (sanitize branch name for filesystem)
worktree_name=$(echo "$new_branch" | sed 's/\//-/g')
worktree_path=$(cd "$REPO_ROOT/.." && pwd)/"$worktree_name"

# Check if worktree path already exists
if [ -d "$worktree_path" ]; then
    echo "Error: Directory already exists at $worktree_path"
    exit 1
fi

# Create the worktree with new branch
echo ""
echo "Creating new branch and worktree..."
git worktree add -b "$new_branch" "$worktree_path" "$base_branch"
echo "✓ Worktree created successfully at $worktree_path"
echo "✓ New branch '$new_branch' created from '$base_branch'"

# Build IDE options dynamically
ide_options=()

# Check for Claude Code options
if command -v claude &> /dev/null; then
    if command -v wezterm &> /dev/null; then
        ide_options+=("Claude Code in Wezterm")
    fi
    if command -v code &> /dev/null; then
        ide_options+=("Claude Code in VSCode")
    fi
fi

# Check for Codex options
if command -v codex &> /dev/null; then
    if command -v wezterm &> /dev/null; then
        ide_options+=("Codex in Wezterm")
    fi
    if command -v code &> /dev/null; then
        ide_options+=("Codex in VSCode")
    fi
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
echo ""
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
    "Claude Code in VSCode")
        echo "Launching Claude Code in VSCode..."
        code "$worktree_path"
        sleep 1
        claude "$worktree_path" &
        ;;
    "Codex in Wezterm")
        echo "Launching Codex in Wezterm..."
        wezterm start --cwd "$worktree_path" -- codex "$worktree_path" &
        ;;
    "Codex in VSCode")
        echo "Launching Codex in VSCode..."
        code "$worktree_path"
        sleep 1
        codex "$worktree_path" &
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

