/**
 * Provider validator script
 *
 * @param {string | string[] | {label: string; required: boolean }} field The input field.
 *
 * This can be one of several types:
 *  - `string` -> The value of the field (e.g. `'my-team'`)
 *  - `string[]` -> The value(s) of the field (e.g. `['team-1', 'team-2']`)
 *  - A checkboxes object with a `label` and `required` property (e.g.
 *    `{ label: 'my-team', required: true }`)
 *
 *  You do not need to handle them all! It is up to the individual validation
 *  script to define which type(s) to expect and how to handle them.
 *
 * @returns {Promise<string>} An error message or `'success'`
 */
module.exports = async (field) => {
    const { Octokit } = require('@octokit/rest')
    const core = require('@actions/core')

    const github = new Octokit({
        auth: core.getInput('github-token', { required: true })
    })

    if (typeof field !== 'string') return 'Field type is invalid';

    if (!/^[a-zA-Z0-9-]+\/terraform-provider-[a-zA-Z0-9-]+$/.test(field))
        return 'Repository must be a valid terraform provider repo matching the pattern "<org>/<repo>/terraform-provider-<name>"';

    try {
        // Check if the provider repository exists
        core.info(`Checking if repository '${field}' exists`)
        const [owner, repo] = field.split('/');

        await github.rest.repos.get({
            owner,
            repo,
        });

        try {
            await github.rest.repos.getContent({ owner: 'opentofu', repo: 'registry', path: `providers/${owner.charAt(0)}/${owner}/${repo}.json` });

            // if we did not get 404 from github - the provider already exists
            core.error(`Repository '${field}' already exists in the registry`)
            return `Provider '${field}' already exists in the registry`
        } catch (error) {
            if (error.status === 404) {
                // the provider does not exist in the registry, so good submission
            } else {
                // Otherwise, something else went wrong...
                throw error
            }
        }

        core.info(`Repository '${field}' exists`)
        return 'success'
    } catch (error) {
        if (error.status === 404) {
            // If the repo does not exist, return an error message
            core.error(`Repository '${field}' does not exist`)
            return `Repository '${field}' does not exist`
        } else {
            // Otherwise, something else went wrong...
            throw error
        }
    }
}