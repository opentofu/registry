# Contributing to OpenTofu Registry

This repository contains the metadata and applications that drive the OpenTofu Registry at [registry.opentofu.org](https://registry.opentofu.org).

> [!IMPORTANT]
> Before writing any code, please read the two rules below. Pull requests that skip them will be closed.
>
> 1. **Every code change needs an issue with the `accepted` label, assigned to you.** This is to protect your time: discussing whether the change is still needed and agreeing on the approach _before_ you write code means you won't spend hours on a pull request we can't merge. Code changes without an accepted, assigned issue will generally be rejected.
> 2. **We do not accept AI-generated contributions from the community.** See our [AI Usage Policy](AI-USAGE-POLICY.md).

> [!NOTE]
> Want to add a provider, module, or GPG key to the registry? That doesn't require a code change. Simply use the issue forms linked in the [README](README.md).

## How to contribute a code change

1. Find an [existing issue](https://github.com/opentofu/registry/issues) or open a new one describing the change you'd like to make.
2. Comment that you'd like to work on it, and wait for a maintainer to add the `accepted` label and assign it to you.
3. Fork the repository and make your changes.
4. Sign off your commits (see [DCO](#developer-certificate-of-origin-dco) below) and open a pull request.

## AI policy

We do not accept AI-generated code from community contributors. This applies to any code generated, completed, or rewritten by an AI tool. If you found a bug or an improvement using AI tools, open an issue describing it in your own words (no code) instead of a pull request.

The full policy, including the rules maintainers follow, is in [AI-USAGE-POLICY.md](AI-USAGE-POLICY.md).

## Developer Certificate of Origin (DCO)

All commits must be signed off to accept the [Developer Certificate of Origin](https://developercertificate.org). Add a `Signed-off-by` line to your commit message:

```
This is my commit message

Signed-off-by: Random Developer <random@developer.example.org>
```

Git can do this for you:

```
~> git commit -s -m 'This is my commit message'
```
