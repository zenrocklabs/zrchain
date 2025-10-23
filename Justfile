# Switch to a branch in a new git worktree with fuzzy search
worktree-switch:
    @./scripts/git/worktree-switch.sh
alias wts := worktree-switch

# Create a new branch in a new git worktree
worktree-new:
    @./scripts/git/worktree-new.sh
alias wtn := worktree-new

# Clean up git worktrees with multi-select
worktree-cleanup:
    @./scripts/git/worktree-cleanup.sh
alias wtc := worktree-cleanup