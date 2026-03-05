# OpenTofu Registry (registry.opentofu.org)

![](https://raw.githubusercontent.com/opentofu/brand-artifacts/main/full/transparent/SVG/on-dark.svg#gh-dark-mode-only)
![](https://raw.githubusercontent.com/opentofu/brand-artifacts/main/full/transparent/SVG/on-light.svg#gh-light-mode-only)

This repository is home to the metadata that drives the provider and module registry for [OpenTofu](https://github.com/opentofu/opentofu)

It also contains the applications used to manage version bumping, validation and API generation of the registry that is hosted at [registry.opentofu.org](https://registry.opentofu.org).

**Thanks to Cloudflare for sponsoring a Business plan to host the registry on!**

## Adding Providers, Modules or GPG Keys to the OpenTofu Registry
To add your provider, module or GPG key to the OpenTofu Registry you can submit an issue using one of the issue templates we provide in this repository.

- [Submit new Module](https://github.com/opentofu/registry/issues/new?assignees=&labels=module%2Csubmission&projects=&template=module.yml&title=Module%3A+)
- [Submit new Provider](https://github.com/opentofu/registry/issues/new?assignees=&labels=provider%2Csubmission&projects=&template=provider.yml&title=Provider%3A+)
- [Submit new Provider Signing Key](https://github.com/opentofu/registry/issues/new?assignees=&labels=provider-key%2Csubmission&projects=&template=provider_key.yml&title=Provider+Key%3A+)

Fill in the required fields and submit the issue. Once the issue has been submitted, the OpenTofu team will review this and either approve or deny the submission.

> [!IMPORTANT]
> **Submissions must be made through the GitHub issue form UI using the links above.** Do not open pull requests to add registry data directly, and do not create issues using the `gh` CLI, the GitHub API, or other tooling. The automated validation and processing pipeline depends on the structured data that only the issue form UI provides. Submissions made outside of the issue form UI will not be processed and will be closed.

## Contributing To The Codebase

Contributions are always welcome!

**Please see [`CONTRIBUTING.md`](CONTRIBUTING.md) for before making any contributions.**

## Registry Inclusion Policy

You can find information on the OpenTofu Registry Inclusion Policy in [POLICY.md](POLICY.md)

## Reporting security vulnerabilities
If you've found a vulnerability or a potential vulnerability in OpenTofu please follow [Security Policy](https://github.com/opentofu/opentofu/security/policy). We'll send a confirmation email to acknowledge your report, and we'll send an additional email when we've identified the issue positively or negatively.
