# Contributing to OpenTofu Registry

This repository contains OpenTofu Registry, which includes the metadata and managing applications used to drive the OpenTofu Registry at registry.opentofu.org

This document provides guidance on OpenTofu contribution recommended practices. It covers how to submit issues, how to get involved in the discussion, how to work on the code, and how to contribute code changes.

The easiest way to contribute is by [opening an issue](https://github.com/opentofu/opentofu/issues/new/choose)! Bug reports, broken compatibility reports, feature requests, old issue reposts, and well-prepared RFCs are all very welcome.

All major changes to the OpenTofu Registry go through the public RFC process, including those proposed by the core team. Thus, if you'd like to propose such a change, please prepare an RFC, so that the community can discuss the change and everybody has a chance to voice their opinion. You're also welcome to voice your own opinion on existing RFCs! You can find them by [going to the issues view and filtering by the rfc label](https://github.com/opentofu/registry/issues?q=is%3Aopen+is%3Aissue+label%3Arfc).

Generally, we appreciate external contributions very much and would love to work with you on them. **However, please make sure to read the [Contributing a Code Change](#contributing-a-code-change) section prior to making a contribution.**

---

<!-- MarkdownTOC autolink="true" -->

- [Contributing a Code Change](#contributing-a-code-change)
- [Working on the Code](#working-on-the-code)
- [Adding or updating dependencies](#adding-or-updating-dependencies)
- [Acceptance Tests: Testing interactions with external services](#acceptance-tests-testing-interactions-with-external-services)
- [Generated Code](#generated-code)

<!-- /MarkdownTOC -->

## Contributing a Code Change

In order to contribute a code change, you should fork the repository, make your changes, and then submit a pull request. Crucially, all code changes should be preceded by an issue that you've been assigned to. If an issue for the change you'd like to introduce already exists, please communicate in the issue that you'd like to take ownership of it. If an issue doesn't yet exist, please create one expressing your interest in working on it and discuss it first, prior to working on the code. Code changes without a related issue will generally be rejected.

Only issues with the `accepted` label have been officially accepted for implementation, so please avoid working on issues without that label.

In order for a code change to be accepted, you'll also have to accept the Developer Certificate of Origin (DCO). It's very lightweight, and you can find it [here](https://developercertificate.org). Accepting is accomplished by signing off on your commits, you can do this by adding a `Signed-off-by` line to your commit message, like here:
```
This is my commit message

Signed-off-by: Random Developer <random@developer.example.org>
```
Git has a built-in flag to append this line automatically:
```
~> git commit -s -m 'This is my commit message'
```

You can find more details about the DCO checker in the [DCO app repo](https://github.com/dcoapp/app).

Additionally, please update [the changelog](CHANGELOG.md) if you're making any user-facing changes.