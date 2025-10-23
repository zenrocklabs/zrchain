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

# Remember the branch we started on so we can restore later
original_branch=$(git branch --show-current)

# Fetch latest branches
echo "Fetching branches..."
git fetch --all --prune

# Get current branch as default base
current_branch=$(git branch --show-current)
default_base="$current_branch"

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
base_branch=$(echo "$branches" | go tool gum filter --placeholder "Select base branch (default: $default_base)..." --height 20)

# If no selection, use current branch
if [ -z "$base_branch" ]; then
    base_branch="$default_base"
    echo "Using default base branch: $base_branch"
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

# Decide whether to use Graphite (if available)
use_graphite=0
if [ "$GT_AVAILABLE" -eq 1 ]; then
    echo ""
    if go tool gum confirm --default "Use Graphite CLI for branch creation?"; then
        use_graphite=1
    fi
fi

if [ "$use_graphite" -eq 1 ]; then
    echo ""
    echo "Checking Graphite tracking for base branch '$base_branch'..."
    set +e
    gt branch info "$base_branch" >/dev/null 2>&1
    tracked_status=$?
    set -e
    if [ $tracked_status -ne 0 ]; then
        echo "Branch '$base_branch' is not tracked by Graphite."
        if go tool gum confirm --default "Track '$base_branch' with Graphite now?"; then
            if gt track "$base_branch"; then
                echo "✓ Branch '$base_branch' is now tracked."
            else
                echo "Failed to track '$base_branch' with Graphite. Exiting."
                exit 1
            fi
        else
            echo "Graphite requires tracked parent branches. Exiting."
            exit 0
        fi
    fi

    echo "Using Graphite to create the new branch..."
    if ! gt checkout "$base_branch" >/dev/null 2>&1; then
        echo "Graphite checkout failed, falling back to git checkout..."
        git checkout "$base_branch"
    fi

    default_message="chore: start ${new_branch}"
    echo ""
    graphite_message=$(go tool gum input --placeholder "Initial commit message for gt create" --value "$default_message")
    if [ -z "$graphite_message" ]; then
        graphite_message="$default_message"
    fi

    gt create "$new_branch" -m "$graphite_message"

    created_branch=$(git branch --show-current)
    if [ -z "$created_branch" ]; then
        created_branch="$new_branch"
    fi
    if [ "$created_branch" != "$new_branch" ]; then
        echo "Graphite created branch '$created_branch' (requested '$new_branch')."

        # Adjust worktree name/path if Graphite chose a different branch name
        adjusted_worktree_name=$(echo "$created_branch" | sed 's/\//-/g')
        if [ "$adjusted_worktree_name" != "$worktree_name" ]; then
            worktree_name="$adjusted_worktree_name"
            worktree_path=$(cd "$REPO_ROOT/.." && pwd)/"$worktree_name"
            if [ -d "$worktree_path" ]; then
                echo "Error: Directory already exists at $worktree_path"
                echo "Cannot create worktree for Graphite branch '$created_branch'. Exiting."
                exit 1
            fi
        fi
    fi

    echo "Switching back to base branch ($base_branch)..."
    if ! gt checkout "$base_branch" >/dev/null 2>&1; then
        git checkout "$base_branch"
    fi

    echo "Creating worktree..."
    git worktree add "$worktree_path" "$created_branch"

    echo "✓ Worktree created successfully at $worktree_path"
    echo "✓ Graphite branch '$created_branch' created from '$base_branch'"

    if [ -n "$original_branch" ] && [ "$original_branch" != "$base_branch" ]; then
        echo "Restoring repository branch to '$original_branch'..."
        if ! gt checkout "$original_branch" >/dev/null 2>&1; then
            git checkout "$original_branch"
        fi
    fi
else
    # Create the worktree with new branch via plain git
    echo ""
    echo "Creating new branch and worktree..."
    git worktree add -b "$new_branch" "$worktree_path" "$base_branch"
    echo "✓ Worktree created successfully at $worktree_path"
    echo "✓ New branch '$new_branch' created from '$base_branch'"

    if [ "$GT_AVAILABLE" -eq 1 ]; then
        echo ""
        echo "Tracking branch with Graphite metadata..."
        set +e
        if ! gt track --parent "$base_branch" "$new_branch" >/dev/null 2>&1; then
            echo "Graphite tracking failed (branch may already be tracked or parent not tracked)."
            echo "You can run 'gt track $new_branch' later if needed."
        else
            echo "✓ Branch registered with Graphite (parent: $base_branch)"
        fi
        set -e
    fi
fi

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
