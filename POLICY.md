# Registry Inclusion Policy

The OpenTofu Registry is an index of providers and modules that work with OpenTofu. The providers and modules themselves are hosted by GitHub, not the OpenTofu Registry.

The OpenTofu Registry service is operated by OpenTofu a Series of LF Projects, LLC under the laws of the United States of America. For terms of use, trademark policy, privacy policy and other project policies please see https://lfprojects.org/policies.

## Provider and Module Submission

Any user with a GitHub account is free to submit a provider or module for inclusion in the OpenTofu registry by using the GitHub issue system.

The following categories of modules and providers will not be included in the OpenTofu Registry and may be removed if found to be included. Note that the decision not to include a provider or module does not constitute legal advice or a finding of fact that a provider or module violates any of these policies, it is merely intended as a measure to protect the OpenTofu project and its maintainers from legal liability.

1. Modules and providers that likely violate [GitHub's Acceptable Use Policies](https://docs.github.com/en/site-policy/acceptable-use-policies/github-acceptable-use-policies) or [Cloudflare's Abuse Policy](https://www.cloudflare.com/trust-hub/abuse-approach/).
2. Modules and providers that promote, support or perform activities likely to be illegal under US law.
3. Modules and providers produced by or in support of entities that are likely to be under embargo, or entities headquartered in or have strong connections to countries that are under a technology embargo under US law.
4. Modules and providers that likely infringe on the intellectual property rights of others or are otherwise likely to be illegal under US law.
5. Modules and providers that contain, install, disseminate malware, disclose sensitive personal or otherwise sensitive information, or in other ways harm OpenTofu users.

Inclusion requests for providers or modules are reviewed by the maintainers and are typically processed without delay or further need for a review. In some cases, the maintainers may refer an inclusion request to the Technical Steering Committee for deliberation with a vote. Removals are always decided by the Technical Steering Committee with a vote. Technical Steering Committee decisions on inclusion or removal are carried out by the maintainers with at least two maintainers approving the pull request.

## Version Immutability

Once a provider or module version has been published and indexed, the OpenTofu Registry treats it as immutable. Consumers and automation rely on the assumption that a given version number always resolves to the same artifacts and checksums. The registry will never modify the artifacts or checksums of an already-indexed version in place; the only action ever available is removing a version from the registry entirely, and even that is reserved for exceptional circumstances, as removing a version on demand undermines that trust for everyone who depends on it, not just the requester.

For this reason, the OpenTofu Registry will **not** remove a version for reasons such as:

- A CI pipeline built and published the release more than once.
- A minor bug, regression, or mistake was discovered after release.
- A provider or module author simply changed their mind about a release.

In all of these cases, the correct fix is for the provider or module author to release a new version in their own repository. Consumers who need the fix can then upgrade to it, while anyone already relying on the original version is not disrupted by it changing or disappearing.

The registry may remove a version on a case-by-case basis for exceptional circumstances, including but not limited to:

- The release contains malicious code, malware, or evidence of a supply chain compromise.
- The release discloses secrets, credentials, or other sensitive information that must be revoked or scrubbed.
- The signing key used for the release has been compromised, making the originally indexed checksums impossible to verify going forward.
- A legal or security requirement (e.g. a takedown notice or court order) mandates removal.
- The release is causing extreme disruption to the registry or its users that cannot reasonably be resolved by releasing a corrected version instead.

These requests are reviewed by the maintainers, following the same escalation to the Technical Steering Committee described above for removals. Meeting one of the criteria above does not guarantee action will be taken. The default answer remains to release a new version in your own repository instead.

If you believe your situation qualifies, open a [Version Removal Request](https://github.com/opentofu/registry/issues/new?template=version-removal-request.yml) issue and make your case, including any relevant evidence such as security advisories or CVEs.

We also encourage provider and module authors to enable [GitHub's immutable releases](https://docs.github.com/en/code-security/concepts/supply-chain-security/immutable-releases) on their repositories. This provides an upstream, technical guarantee that complements this policy, though the registry's own practice of not removing versions applies regardless of whether a given repository has this feature enabled.

## Reporting violations

If you believe an included module or provider violates applicable laws, please primarily contact GitHub for removal of the repository hosting it. For details, please refer to the [GitHub Content Removal Policies](https://docs.github.com/en/site-policy/content-removal-policies).

In rare cases a provider or module may not be removed by GitHub, but its inclusion in the registry may still violate our policies. In this case, you may report policy violations by writing an email to [liaison@opentofu.org](mailto:liaison@opentofu.org). Please note, unless required by law, the OpenTofu team has sole discretion on removing content and may decide not to remove a provider or module if deemed to be in the best interests of the OpenTofu project and its users. As a general rule, actions taken (if any) will be documented on GitHub at the discretion of the TSC and your email will not receive a response unless required by law.

Please also note that should the need arise, we may publish your report, whether action is taken or not, as a measure of transparency, with sensitive information redacted.

## Alternatives to the OpenTofu Registry

Some organizations may have need to host their own registry for security, compliance, or legal purposes. More information on hosting your own registry can be found in the [OpenTofu documentation](https://opentofu.org/docs/cli/private_registry/).

## Changes to this policy

This policy may be changed at any time based on the decision of the OpenTofu Technical Steering Committee. Changes will be published in the OpenTofu Registry GitHub repository.
