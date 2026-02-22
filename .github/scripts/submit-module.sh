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

# Post initial comment with Actions run link and validation summary
ACTION_URL="${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/actions/runs/${GITHUB_RUN_ID}"
gh issue comment "${NUMBER}" -b "## Automated Validation Started

**GitHub Actions Run:** ${ACTION_URL}

### Validation Steps
- ✓ Checking module repository format
- ✓ Validating module metadata
- ✓ Creating module JSON file
- ✓ Opening pull request

Results will be posted here when validation completes."

repository=$(echo "${BODY}" | grep "### Module Repository" -A2 | tail -n1 | tr "[:upper:]" "[:lower:]" | sed -e 's/[\r\n]//g')
repository=$(echo -n "${repository}" | sed -e 's|https://github.com/||' -e 's|github.com/||')

if [[ ! "${repository}" =~ ^[a-zA-Z0-9-]+/terraform-[a-zA-Z0-9-]+$ ]]; then
  gh issue comment "${NUMBER}" -b "Failed validation: Invalid repository name: '${repository}'. Please edit your issue to state the name of the repository in the format of ORGANIZATION/terraform-NAME-TARGETSYSTEM."
  exit 1
fi

# Detect if the user accidentally submitted a provider repo as a module.
# Provider repos follow the pattern terraform-provider-NAME, which also matches
# the module regex above. Catch this early with a helpful redirect message.
if [[ "${repository}" =~ ^[a-zA-Z0-9-]+/terraform-provider-[a-zA-Z0-9-]+$ ]]; then
  gh issue comment "${NUMBER}" -b "Failed validation: It looks like \`${repository}\` is a **provider** repository, not a module. Provider repos follow the \`terraform-provider-NAME\` naming convention. If you meant to submit a provider, please open a new issue using the [provider submission form](https://github.com/opentofu/registry/issues/new?template=provider.yml). Module repositories should follow the pattern \`ORGANIZATION/terraform-NAME-TARGETSYSTEM\`."
  exit 1
fi

set +e
if ! go run ./cmd/add-module -repository="${repository}" -output=./output.json ; then
  set -euo pipefail
  if [[ "$(jq -r '.exists' < ./output.json || true)" == "true" ]]; then
    gh issue close "${NUMBER}" -c "$(jq -r '.validation' < ./output.json || true)"
    exit 0
  else
    gh issue comment "${NUMBER}" -b "$(jq -r '.validation' < ./output.json || true)"
    exit 1
  fi
fi
set -euo pipefail
namespace=$(jq -r '.namespace' < ./output.json)
name=$(jq -r '.name' < ./output.json)
target=$(jq -r '.target' < ./output.json)
jsonfile=$(jq -r '.file' < ./output.json)


# Create Branch
branch=module-submission_${namespace}_${name}_${target}
set +e
if ! git checkout -b "${branch}"; then
  set -euo pipefail
  gh issue comment "${NUMBER}" -b "Failed validation: A branch already exists for this module '${branch}'"
  exit 1
fi
set -euo pipefail

# Add result
git add "${jsonfile}"

# Commit and push result
git config --global user.email "no-reply@opentofu.org"
git config --global user.name "OpenTofu Automation"
git commit -s -m "Create module ${namespace}/${name}/${target}"
git push -u origin "${branch}"

# Create pull request and update issue
pr=$(gh pr create --title "${TITLE}" --body "Created ${jsonfile/../src/} for module ${namespace}/${name}/${target}.  Closes #${NUMBER}.") #--assignee opentofu/core-engineers)
gh issue comment "${NUMBER}" -b "Your submission has been validated and has moved on to the pull request phase (${pr}).  This issue has been locked."
gh issue lock "${NUMBER}" -r resolved
