#!/bin/bash
set -x

keyfile="$1"
owner="$2"
provider_name="$3"
if [[ -z "${NUMBER}" ]]; then
  echo "Please run this script from a GitHub Action. The issue NUMBER is missing"
  exit 1
fi
if [[ -z "${keyfile}" ]]; then
  echo "no keyfile provided"
  exit 1
fi
if [[ -z "${owner}" ]]; then
  echo  "no owner provided"
  exit 2
fi
if [[ "${owner}" == "hashicorp" || "${owner}" == "opentofu" ]]; then
  echo "No checks required for 'hashicorp' or 'opentofu' namespaced providers"
  exit 0
fi

function check_repo_release() {
  local owner="${1}"
  local repo="${2}"
  local release="${3}"
  # download the GPG signature from the given provider release
  gh release download --repo "${owner}/${repo}" "${release}" -p "*SHA256*"
  # verify the signatures
  # shellcheck disable=SC2312
  sigfile=$(find . -name "*SHA256SUMS.sig" -print | head -1)
  # shellcheck disable=SC2312
  shafile=$(find . -name "*SHA256SUMS" -print | head -1)
  if gpg --verify "${sigfile}" "${shafile}" > /dev/null 2>&1
  then
    echo "Key is matching signatures from ${owner}/${repo}@${release}"
    rm "${shafile}" "${sigfile}"
    return 0
  fi
  echo "Key does not match signatures from ${owner}/${repo}@${release}"
  rm "${shafile}" "${sigfile}"
  return 1
}

function check_repo_versions() {
  local owner="${1}"
  local repo="${2}"
  local releases
  releases="$(gh release list --exclude-drafts --exclude-pre-releases --repo "${owner}/${repo}" -L 3 -O desc --json name -q '.[].name')"
  # check recent releases of the owner's repo (3 releases checked)
  while IFS= read -r release; do
    if check_repo_release "${owner}" "${repo}" "${release}"
    then
      # once one release is matching the signature, we are good to go, so return success
      return 0
    fi
  # list the latest 100 releases of the repository and get only the release names
  done <<< "${releases}"
  # if no release is matching the signature, return error
  return 1
}

function check_owner_repos() {
  local owner="${1}"
  # list first 100 repos of the owner and get all the terraform-provider-* repos to check their releases
  local repos
  repos="$(gh repo list "${owner}" --no-archived --source -L 100 --json name -q '.[].name | select(. | contains("terraform-provider-"))')"
  while IFS= read -r repo; do
    if check_repo_versions "${owner}" "${repo}" "${release}"
    then
      return 0
    fi
  done <<< "${repos}"
  return 1
}

# prepare gpg
apt update && apt install -y gpg
# import the submitted key
gpg --import "${keyfile}" 2>/dev/null
# trust the newly imported key
# shellcheck disable=SC2312
for fpr in $(gpg --list-keys --with-colons | grep "pub:" | awk -F: '{print $5}' | sort -u); do  echo -e "5\ny\n" | gpg -q --command-fd 0 --expert --edit-key "${fpr}" trust; done

if [[ -n "${provider_name}" ]]; then
  # if the submission contains also the provider name, we will check the signatures only of that particular provider
  repo="terraform-provider-${provider_name}"
  if ! check_repo_versions "${owner}" "${repo}"
  then
    gh issue comment "${NUMBER}" -b "Key is matching no recent release of ${owner}/${repo}"
    exit 0
  fi
else
  # if no provider name is given, will check the key against any terraform-provider-* repo of the owner
  if ! check_owner_repos "${owner}"
  then
    gh issue comment "${NUMBER}" -b "Key is matching no recent release from any 'terraform-provider-*' of ${owner}"
    exit 0
  fi
fi
gh issue comment "${NUMBER}" -b "Key provider signatures validation succeeded!"

# cleanup keys
# shellcheck disable=SC2312
for fpr in $(gpg --list-keys --with-colons -q | grep "pub:" | awk -F: '{print $5}' | sort -u); do  echo -e "y\n" | gpg --command-fd 0 --expert --delete-keys "${fpr}"; done

