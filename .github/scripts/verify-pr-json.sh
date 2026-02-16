#!/bin/bash

set -euo pipefail

if [[ -z "${PR_NUMBER:-}" ]]; then
	echo "PR_NUMBER is required."
	exit 1
fi

failed=0
error_file="/tmp/verify-errors.md"

write_error_header() {
	if [[ ! -f "${error_file}" ]]; then
		cat >"${error_file}" <<'HEADER'
## JSON Verification Failed

The JSON files in this PR do not match what is generated from the source repositories. This may indicate tampered data or a race condition with new upstream releases.

HEADER
	fi
}

# Compare two JSON files semantically (normalized key order and formatting)
# Stores diff output in $json_diff_output. Sets return code 0 if equal, 1 if different.
json_diff() {
	json_diff_output=$(diff -u <(jq --sort-keys '.' "$1") <(jq --sort-keys '.' "$2") || true)
	[[ -z "${json_diff_output}" ]]
}

# Get changed provider and module JSON files via GitHub API
changed_providers=$(gh pr diff "${PR_NUMBER}" --name-only | grep '^providers/.*\.json$' || true)
changed_modules=$(gh pr diff "${PR_NUMBER}" --name-only | grep '^modules/.*\.json$' || true)

if [[ -z "${changed_providers}" && -z "${changed_modules}" ]]; then
	echo "No provider or module JSON changes detected."
	exit 0
fi

# Verify each changed provider JSON
for file in ${changed_providers}; do
	# Path format: providers/{first_char}/{namespace}/{name}.json
	namespace=$(echo "${file}" | cut -d'/' -f3)
	name=$(basename "${file}" .json)
	repo="${namespace}/terraform-provider-${name}"

	tmpdir=$(mktemp -d)

	echo "::group::Verifying provider ${repo}"
	echo "Regenerating JSON for ${repo}..."

	if ! go run ./cmd/add-provider -repository="${repo}" -provider-data="${tmpdir}" -output="${tmpdir}/output.json" 2>&1; then
		echo "::error::Failed to regenerate provider ${repo}"
		write_error_header
		cat >>"${error_file}" <<EOF
### Provider \`${repo}\`
Failed to regenerate JSON from source repository.

EOF
		failed=1
		rm -rf "${tmpdir}"
		echo "::endgroup::"
		continue
	fi

	first_char=$(echo "${namespace}" | cut -c1 | tr '[:upper:]' '[:lower:]')
	generated_file="${tmpdir}/${first_char}/${namespace}/${name}.json"

	if [[ ! -f "${generated_file}" ]]; then
		echo "::error::Regeneration produced no output file for provider ${repo}"
		write_error_header
		cat >>"${error_file}" <<EOF
### Provider \`${repo}\`
Regeneration produced no output file.

EOF
		failed=1
		rm -rf "${tmpdir}"
		echo "::endgroup::"
		continue
	fi

	# shellcheck disable=SC2310
	if ! json_diff "../${file}" "${generated_file}"; then
		echo "${json_diff_output}"
		echo "::error::Provider ${repo} JSON does not match regenerated output. PR may contain tampered data."
		write_error_header
		cat >>"${error_file}" <<EOF
### Provider \`${repo}\`
JSON does not match regenerated output:
\`\`\`diff
${json_diff_output}
\`\`\`
EOF
		failed=1
	else
		echo "PASS: Provider ${repo} JSON matches regenerated output."
	fi

	rm -rf "${tmpdir}"
	echo "::endgroup::"
done

# Verify each changed module JSON
for file in ${changed_modules}; do
	# Path format: modules/{first_char}/{namespace}/{name}/{target}.json
	namespace=$(echo "${file}" | cut -d'/' -f3)
	name=$(echo "${file}" | cut -d'/' -f4)
	target=$(basename "${file}" .json)
	repo="${namespace}/terraform-${target}-${name}"

	tmpdir=$(mktemp -d)

	echo "::group::Verifying module ${repo}"
	echo "Regenerating JSON for ${repo}..."

	if ! go run ./cmd/add-module -repository="${repo}" -module-data="${tmpdir}" -output="${tmpdir}/output.json" 2>&1; then
		echo "::error::Failed to regenerate module ${repo}"
		write_error_header
		cat >>"${error_file}" <<EOF
### Module \`${repo}\`
Failed to regenerate JSON from source repository.

EOF
		failed=1
		rm -rf "${tmpdir}"
		echo "::endgroup::"
		continue
	fi

	first_char=$(echo "${namespace}" | cut -c1 | tr '[:upper:]' '[:lower:]')
	generated_file="${tmpdir}/${first_char}/${namespace}/${name}/${target}.json"

	if [[ ! -f "${generated_file}" ]]; then
		echo "::error::Regeneration produced no output file for module ${repo}"
		write_error_header
		cat >>"${error_file}" <<EOF
### Module \`${repo}\`
Regeneration produced no output file.

EOF
		failed=1
		rm -rf "${tmpdir}"
		echo "::endgroup::"
		continue
	fi

	# shellcheck disable=SC2310
	if ! json_diff "../${file}" "${generated_file}"; then
		echo "${json_diff_output}"
		echo "::error::Module ${repo} JSON does not match regenerated output. PR may contain tampered data."
		write_error_header
		cat >>"${error_file}" <<EOF
### Module \`${repo}\`
JSON does not match regenerated output:
\`\`\`diff
${json_diff_output}
\`\`\`
EOF
		failed=1
	else
		echo "PASS: Module ${repo} JSON matches regenerated output."
	fi

	rm -rf "${tmpdir}"
	echo "::endgroup::"
done

if [[ "${failed}" -ne 0 ]]; then
	echo ""
	echo "::error::Verification failed. One or more JSON files do not match regenerated output."
	exit 1
fi

echo ""
echo "All provider and module JSON files verified successfully."
