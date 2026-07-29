#!/bin/bash

set -x
set -euo pipefail

if [[ -z "${BODY}" ]]; then
  echo "Please run this script from a GitHub Action."
  exit 1
fi
if [[ -z "${TITLE}" ]]; then
  echo "Please run this script from a GitHub Action."
  exit 1
fi
if [[ -z "${NUMBER}" ]]; then
  echo "Please run this script from a GitHub Action."
  exit 1
fi
if [[ -z "${GH_USER}" ]]; then
  echo "Please set GH_USER"
  exit 1
fi

namespace=$(echo "${BODY}" | grep "### Provider Namespace" -A2 | tail -n1 | tr "[:upper:]" "[:lower:]" | sed -e 's/[\r\n]//g')
providername=$(echo "${BODY}" | grep "### Provider Name" -A2 | tail -n1 | tr "[:upper:]" "[:lower:]" | sed -e 's/[\r\n]//g')
version=$(echo "${BODY}" | grep "### Version" -A2 | tail -n1 | sed -e 's/^v//' -e 's/[\r\n]//g')

if [[ ! "${namespace}" =~ ^[a-zA-Z0-9-]+$ ]]; then
  gh issue comment "${NUMBER}" -b "Failed validation: Invalid namespace: '${namespace}'"
  exit 1
fi

if [[ ! "${providername}" =~ ^[a-zA-Z0-9-]+$ ]]; then
  gh issue comment "${NUMBER}" -b "Failed validation: Invalid provider name: '${providername}'"
  exit 1
fi

if [[ ! "${version}" =~ ^[0-9]+\.[0-9]+\.[0-9]+.*$ ]]; then
  gh issue comment "${NUMBER}" -b "Failed validation: Invalid version: '${version}'"
  exit 1
fi

set +e
go run ./cmd/verify-reindex-request -namespace "${namespace}" -name "${providername}" -version "${version}" -username "${GH_USER}" -provider-data "../providers" -gpg-data "../keys" -output=./output.json
verification=$?
set -euo pipefail

gh issue comment "${NUMBER}" -b "$(jq -r '.' < ./output.json || true)"
if [[ "${verification}" != 0 ]]; then
  exit 1
fi

go run ./cmd/remove-provider-version -namespace "${namespace}" -name "${providername}" -version "${version}" -provider-data "../providers"

# Create Branch
branch="reindex-request_${namespace}-${providername}-${version}-$(date +%s)"
set +e
if ! git checkout -b "${branch}"; then
  gh issue comment "${NUMBER}" -b "Failed validation: A branch already exists for this re-index request '${branch}'"
  exit 1
fi
set -euo pipefail

# Add result
providerfile="../providers/${namespace:0:1}/${namespace}/${providername}.json"
git add "${providerfile}"

# Commit and push result
git config --global user.email "no-reply@opentofu.org"
git config --global user.name "OpenTofu Automation"
git commit -s -m "Fixes #${NUMBER}: removing ${namespace}/${providername} version ${version} for re-index"
git push -u origin "${branch}"

# Create pull request and update issue
# GITHUB_SERVER_URL, GITHUB_REPOSITORY, GITHUB_RUN_ID are default GitHub Actions env vars
# shellcheck disable=SC2154
run_url="${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/actions/runs/${GITHUB_RUN_ID}"
pr_body="Removes ${namespace}/${providername} version ${version} so it will be re-indexed on the next run. **This should be merged and then a resync triggered immediately after.**

- **Source Issue:** #${NUMBER}
- **Requested by:** @${GH_USER}
- **Created by:** [GitHub Actions Run](${run_url})

Closes #${NUMBER}."
pr=$(gh pr create --title "${TITLE}" --body "${pr_body}")
gh issue comment "${NUMBER}" -b "Your submission has been validated and has moved on to the pull request phase (${pr}). This issue has been locked."
gh issue lock "${NUMBER}" -r resolved
