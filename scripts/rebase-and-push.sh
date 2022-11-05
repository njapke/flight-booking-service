#!/usr/bin/env bash

set -euo pipefail

for branch in perf-issue-request-id perf-issue-clean-path perf-issue-basic-auth; do
    git rebase main "$branch"
    git push origin "$branch" -f
done
