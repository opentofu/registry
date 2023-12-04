/**
 * Provider validator script used for further custom validation consumed by https://github.comissue-ops/validator
 *
 * @param {string | string[] | {label: string; required: boolean }} field The input field.
 *
 * @returns {Promise<string>} An error message or `'success'`
 */
module.exports = async (field) => {
    const core = require('@actions/core');

    if (typeof field !== 'string') return 'Field type is invalid';

    if (!/^[a-zA-Z0-9-]+\/terraform-provider-[a-zA-Z0-9-]+$/.test(field))
        return 'Repository must be a valid terraform provider repo matching the pattern "<org>/<repo>/terraform-provider-<name>"';

    // Check if the provider repository exists
    core.info(`Checking if repository '${field}' exists`);
    const [owner, repo] = field.split('/');

    // Check if the repository exists
    const repoResponse = await fetch(`https://api.github.com/repos/${owner}/${repo}`);
    if (!repoResponse.ok) {
        if (repoResponse.status === 404) {
            // If the repo does not exist, return an error message
            core.error(`Repository '${field}' does not exist`);
            return `Repository '${field}' does not exist`;
        } else {
            throw new Error(`Failed to check repository existence: ${repoResponse.statusText}`);
        }
    }

    // Check if the provider repository exists in the registry
    const registryResponse = await fetch(`https://raw.githubusercontent.com/opentofu/registry/main/providers/${owner.charAt(0)}/${owner}/${repo}.json`);
    if (registryResponse.ok) {
        // If the provider exists in the registry, return an error message
        core.error(`Repository '${field}' already exists in the registry`);
        return `Provider '${field}' already exists in the registry`;
    } else if (registryResponse.status !== 404) {
        throw new Error(`Failed to check registry: ${registryResponse.statusText}`);
    }

    core.info(`Repository '${field}' exists`);
    return 'success';
};
