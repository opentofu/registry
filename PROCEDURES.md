# Procedures for the OpenTofu Registry

This document outlines the procedures the core team follows when working with the registry. It is meant as a guide, not as a ruleset, and occasionally we may deviate from this guide.

## Provider submissions

Anyone can submit a provider to the OpenTofu Registry, it does not have to be the provider author. The automation will validate the submission and create a pull request from it. Core maintainers can merge the pull request immediately if it matches the [inclusion policy](POLICY.md), no approval is needed.

Typical errors:

- **Incorrect repository name format:** The automation expects the provider repo in the format of `NAMESPACE/terraform-provider-NAME`. Often users enter it differently. If it is clear what the submitter meant, please edit the submission issue to correct the mistake. Editing the issue will re-run the automation.

If you can't correct an error on behalf of a submitter, please comment explicitly and explain to the submitter what the problem is and how to fix it.

## Module submissions

Module submissions work similar to provider submissions. As long as they meet the inclusion policy, they can be immediately merged.

Typical errors:

- **Incorrect repository name format:** The automation expects the module repo in the format of `NAMESPACE/terraform-TARGETSYSTEM-NAME`. Often users enter it differently. If it is clear what the submitter meant, please edit the submission issue to correct the mistake. Editing the issue will re-run the automation.

If you can't correct an error on behalf of a submitter, please comment explicitly and explain to the submitter what the problem is and how to fix it.

## GPG key submissions

Only provider authors may submit a GPG key on behalf of an organization.

Typical errors:

- **Incorrect repository name:** Submitters can decide to submit a key for the entire namespace or for a single provider. However, often users enter the repo name or namespace incorrectly. If it is clear what the submitter meant, please edit the issue to correct the mistake.
- **User verification failed:** The OpenTofu Registry validates this by checking if the submitting user is a **public** member of the organization the provider belongs to. Third parties are not allowed to submit GPG keys. Please ask the submitter to change their visibility status and then edit the issue and add a space at the end. You may also edit the issue and do so yourself if you can see that the user has already made their org membership public. 
- **Expired key:** The Terraform registry does not verify the GPG key expiry. OpenTofu does a check at submission time. When you encounter this issue, you can help the submitter by explaining how to extend their key:
  > Hey NAME, it looks like your key is expired. You can extend the key signature with the following commands:
  > 1. `gpg --list-keys` to get the key ID.
  > 2. `gpg --edit-key KEYID`
  > 3. Select the correct key by typing `key NUMBER`
  > 4. Run `expire` and answer the questions.
  > 5. Type `save`.
  > 6. Export your newly extended key by running `gpg --export KEYID`

## User complaint: the provider is not signed with a valid signing key

Often provider authors submit keys that are valid, but some or all of the providers are signed with a different key. This typically happens a few hours-days after the submission.

You can debug this on behalf of the user by performing the following steps:

1. Import the submitted key by running `gpg --import FILENAME`
2. Run `gpg --edit-key KEYID` and then type `trust` to trust the key.
3. Download the `...SHA256SUMS` and `...SHA256SUMS.sig` file of the provider.
4. Run `gpg --verify ...SHA256SUMS.sig ...SHA256SUMS`

## User complaint: provider version is not available

In rare cases it can happen that a provider version is not available in the OpenTofu Registry. Check the following:

1. Has the provider version been released less than 30 minutes ago? If so, you may just need to wait.
2. Is the provider available in the [OpenTofu Registry Search](https://search.opentofu.org)? If not, the provider may just not be indexed in the search interface. Check the scheduled jobs in the [registry-ui repository](https://github.com/opentofu/registry-ui).
3. Is the provider in the [registry dataset](https://github.com/opentofu/registry/tree/main/providers)? If not, then the provider version may be faulty. Check the scheduled jobs in the [registry repository](https://github.com/opentofu/registry).
4. If the provider is in the registry dataset, but not available from `tofu init`, check if the generated file is available in the Cloudflare R2 bucket. If so, you may need to clear the Cloudflare cache. 
5. If the provider is one of the HashiCorp-mirrored ones, check if the provider has a release in the repository in the OpenTofu organization. The mirroring may be broken.