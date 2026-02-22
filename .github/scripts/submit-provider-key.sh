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

# Post initial comment with Actions run link and validation summary
ACTION_URL="${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/actions/runs/${GITHUB_RUN_ID}"
gh issue comment "${NUMBER}" -b "## Automated Validation Started

**GitHub Actions Run:** ${ACTION_URL}

### Validation Steps
- ✓ Checking provider namespace and name
- ✓ Validating GPG public key format
- ✓ Verifying key ownership
- ✓ Adding key to registry
- ✓ Opening pull request

Results will be posted here when validation completes."

namespace=$(echo "${BODY}" | grep "### Provider Namespace" -A2 | tail -n1 | tr "[:upper:]" "[:lower:]" | sed -e 's/[\r\n]//g')
providername=$(echo "${BODY}" | grep "### Provider Name" -A2 | tail -n1 | tr "[:upper:]" "[:lower:]" | sed -e 's/[\r\n]//g')
keydata=$(echo "${BODY}" | grep -A 1000 "BEGIN PGP PUBLIC KEY BLOCK"  | grep -B 1000 "END PGP PUBLIC KEY BLOCK")
echo "${keydata}" > tmp.key

if [[ ! "${namespace}" =~ ^[a-zA-Z0-9-]+$ ]]; then
  gh issue comment "${NUMBER}" -b "Failed validation: Invalid namespace: '${namespace}'"
  exit 1
fi

if [[ "${providername}" = "_no response_" ]]; then
  providername=""
fi
if [[ -n "${providername}" ]]; then
  if [[ ! "${providername}" =~ ^[a-zA-Z0-9-]+$ ]]; then
    gh issue comment "${NUMBER}" -b "Failed validation: Invalid provider name: '${providername}'"
    exit 1
  fi

  if [[ "${providername}" =~ ^terraform-provider-.*$ ]]; then
    gh issue comment "${NUMBER}" -b "Failed validation: It seems like you accidentally added the 'terraform-provider-' prefix: '${providername}'"
    exit 1
  fi
fi

set +e
go run ./cmd/verify-gpg-key -org "${namespace}" -provider-data "../providers" -provider-name "${providername}" -username "${GH_USER}" -key-file=tmp.key -output=./output.json
verification=$?
set -euo pipefail

gh issue comment "${NUMBER}" -b "$(jq -r '.' < ./output.json || true)"
if [[ "${verification}" != 0 ]]; then
  exit 1
fi

if [[ -z "${providername}" ]]; then
  keyfile="../keys/${namespace:0:1}/${namespace}/provider-$(date +%s).asc"
else
  keyfile="../keys/${namespace:0:1}/${namespace}/${providername}/provider-$(date +%s).asc"
fi
if [[ -d "$(dirname "${keyfile}")" ]]; then
  msg="Updated"
  #git rm $(dirname $keyfile)/*
else
  msg="Created"
fi
mkdir -p "$(dirname "${keyfile}")"
mv tmp.key "${keyfile}"

# Create Branch
branch="provider-key-submission_${namespace}-$(date +%s)"
set +e
if ! git checkout -b "${branch}"; then
  gh issue comment "${NUMBER}" -b "Failed validation: A branch already exists for this provider '${branch}'"
  exit 1
fi
set -euo pipefail

# Add result
git add "${keyfile}"

# Commit and push result
git config --global user.email "no-reply@opentofu.org"
git config --global user.name "OpenTofu Automation"
if [[ -n "${providername}" ]]; then
  git commit -s -m "Create provider key ${namespace}/${providername}"
else
  git commit -s -m "Create provider key ${namespace}"
fi
git push -u origin "${branch}"

# Create pull request and update issue
pr=$(gh pr create --title "${TITLE}" --body "${msg} ${keyfile/.././} for provider ${namespace}. Closes #${NUMBER}.") #--assignee opentofu/core-engineers)
gh issue comment "${NUMBER}" -b "Your submission has been validated and has moved on to the pull request phase (${pr}).  This issue has been locked."
gh issue lock "${NUMBER}" -r resolved
