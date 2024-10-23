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

Inclusion requests for providers or modules are reviewed by the core developers and are typically processed without delay or further need for a review. In some cases, the core developers may refer an inclusion request to the Technical Steering Committee for deliberation with a vote. Removals are always decided by the Technical Steering Committee with a vote. Technical Steering Committee decisions on inclusion or removal are carried out by the core developers with at least two core developers approving the pull request. 

## Reporting violations

If you believe an included module or provider violates applicable laws, please primarily contact GitHub for removal of the repository hosting it. For details, please refer to the [GitHub Content Removal Policies](https://docs.github.com/en/site-policy/content-removal-policies).

In rare cases a provider or module may not be removed by GitHub, but its inclusion in the registry may still violate our policies. In this case, you may report policy violations by writing an email to [liaison@opentofu.org](mailto:liaison@opentofu.org). Please note, unless required by law, the OpenTofu team has sole discretion on removing content and may decide not to remove a provider or module if deemed to be in the best interests of the OpenTofu project and its users. As a general rule, actions taken (if any) will be documented on GitHub at the discretion of the TSC and your email will not receive a response unless required by law.

Please also note that should the need arise, we may publish your report, whether action is taken or not, as a measure of transparency, with sensitive information redacted.

## Alternatives to the OpenTofu Registry

Some organizations may have need to host their own registry for security, compliance, or legal purposes. More information on hosting your own registry can be found in the [OpenTofu documentation](https://opentofu.org/docs/cli/private_registry/).

## Changes to this policy

This policy may be changed at any time based on the decision of the OpenTofu Technical Steering Committee. Changes will be published in the OpenTofu Registry GitHub repository.
