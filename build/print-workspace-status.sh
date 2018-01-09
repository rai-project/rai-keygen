#!/usr/bin/env bash

# This command is used by bazel as the workspace_status_command
# to implement build stamping with git information.

set -o errexit
set -o nounset
set -o pipefail

GIT_COMMIT=$(git rev-parse --short HEAD)
GIT_TAG=$(git describe --abbrev=0 --tags 2>/dev/null || echo "0.0.0")
GIT_SUMMARY=$(git describe --tags --dirty --always)
GIT_BRANCH=$(git symbolic-ref -q --short HEAD)
GIT_STATE=$(git status --porcelain)
BUILD_DATE=$(date -u '+%Y-%m-%d_%I:%M:%S%p')
VERSION=$(cat VERSION)


cat <<EOF
BuildDate ${BUILD_DATE-}
GitCommit ${GIT_COMMIT-}
GitBranch ${GIT_BRANCH-}
GitState  ${GIT_STATE-}
GitTag    ${GIT_TAG-}
Version   ${VERSION-}
EOF
