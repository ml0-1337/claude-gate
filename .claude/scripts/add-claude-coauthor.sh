#!/bin/bash

# Script to add Claude as co-author to commits
# Usage: ./scripts/add-claude-coauthor.sh [number_of_commits]
# Default: 10 commits

set -e

# Get number of commits from argument, default to 10
COMMIT_COUNT="${1:-10}"

echo "=== Adding Claude as co-author to last $COMMIT_COUNT commits ==="
echo ""

# Create backup branch
BACKUP_BRANCH="backup-before-coauthor-$(date +%Y%m%d-%H%M%S)"
echo "Creating backup branch: $BACKUP_BRANCH"
git branch "$BACKUP_BRANCH"

echo "Adding Claude as co-author to commits..."

# Use filter-branch to update commit messages
FILTER_BRANCH_SQUELCH_WARNING=1 git filter-branch -f --msg-filter '
    cat
    if ! grep -q "Co-Authored-By: Claude" ; then
        echo ""
        echo "🤖 Generated with [Claude Code](https://claude.ai/code)"
        echo ""
        echo "Co-Authored-By: Claude <noreply@anthropic.com>"
    fi
' HEAD~"$COMMIT_COUNT"..HEAD

echo ""
echo "✅ Done! Commits have been updated with Claude as co-author."
echo ""
echo "To verify the changes:"
echo "  git log --format=\"%h %s\" --grep=\"Co-Authored-By\" -$COMMIT_COUNT"
echo ""
echo "To see full commit messages:"
echo "  git log --format=\"%h %s%n%b%n---\" -5"
echo ""
echo "If something went wrong, you can restore with:"
echo "  git reset --hard $BACKUP_BRANCH"
echo "  git branch -D $BACKUP_BRANCH"
echo ""
echo "To push the changes (⚠️  this will rewrite history):"
echo "  git push --force origin main"
echo ""
echo "To clean up the backup branch after verifying:"
echo "  git branch -D $BACKUP_BRANCH"