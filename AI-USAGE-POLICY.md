> [!NOTE]
>  All revisions to this document must be handwritten and cannot use AI in either assistance or ideation.

Each repository in the OpenTofu org has its own AI policy, which is determined on a case-by-case basis. This policy applies solely to this repository. 

### Reason for this policy

Automated tooling, whether using LLMs or other static analysis, is easy to scale. We have yet to determine a way to scale our maintainers' attention. The quality standards for code merged into this repository have not changed because AI tools exist. We have seen a rise in AI-generated pull requests across the open source community as a whole including patches that authors cannot explain, fixes for problems that do not exist and “improvements” that break working code.

Beyond quality, the OpenTofu maintainers cannot verify the origins of any given model’s training data. A model may reproduce code from a license-incompatible source, most notably the BUSL-licensed codebase of HashiCorp’s Terraform, from which OpenTofu at its core was forked. Because we are unable to audit this, we have decided that the risk is too high.

This policy is purposely stricter than the [Linux Foundation's Generative AI Policy](https://www.linuxfoundation.org/legal/generative-ai), which explicitly permits individual projects to set their own rules. OpenTofu's situation as a fork from a now-proprietary codebase warrants additional caution.

### Rules for community contributors

OpenTofu welcomes human-written contributions from the open source community. We are happy to review pull requests that you wrote yourself, line by line. 

**We do not accept AI-generated code from community contributors.** This applies to any code generated, completed or rewritten by an AI tool. If an AI wrote part of your contribution, we cannot accept it. 

**If you discover a bug or possible improvement using AI tools**, feel free to open an issue instead of a pull request. Describe the problem in your own words and a maintainer or other community member may discuss this, pick it up, and implement a fix under the same rules of this policy and contributing guidelines. We emphasize that no code should be included in the description of these issues.

Pull requests identified as AI-generated will be closed without further review.

### Rules for Maintainers

The rules below apply to OpenTofu maintainers who are using AI tooling when contributing to this repository. These rules sit alongside the existing review and quality standards and not in place of them. AI does not change what is expected from a maintainer’s pull request, the standards bar is still the same.

#### Approved Tooling

Github Copilot Enterprise is the only approved AI tool on this repository as of 25th June 2026. Other tools may be evaluated and added to this list in the future, but until then, all other AI coding assistants are not permitted on this codebase.

#### Disclosure

If you did use AI tooling in any part of this contribution, you must disclose this in the pull request description. “Any part” here includes generated code, AI suggested edits, AI assisted refactoring, and documentation. See Commit attribution below.

#### Author responsibility
You must be able to understand and explain every line you submit. Maintainers remain fully responsible for the quality, correctness, and provenance of the code they contribute. The standards for pull requests have not changed because AI tooling exists. 

If you as the author of a contribution cannot explain a line of code on your own pull request during review then you should not be submitting it.

#### Commit messages

Commit messages must be human-written. The commit message is a description of your intent as a human. As this is a critical part of conveying understanding of the contribution it is important to ensure that AI tooling is not used to create commit messages.

#### Commit attribution

Every commit that uses AI assistance must include an `Assisted-by` trailer at the end of the commit message.

Format :
```
Assisted-by: Github Copilot (Model: <modelname>)
```

For example:
```
Assisted-by: GitHub Copilot (Model: Claude Opus 4.7)
```

**Do not use `Co-authored-by:` for AI tools.**

> [!IMPORTANT] Squash merging
> When merging a pull request into its target branch, GitHub may not carry individual commit trailers into the squashed commit message. The maintainer performing the merge must ensure that a copy of the `Assisted-by:` trailer (or trailers) is made into the final squashed commit message.

#### Reviewing contributions

When reviewing a pull request, your review must be human-written. AI may be used to **assist** a review, for example as a second pass after you have formed your own opinion, but this cannot replace your judgement as a maintainer. Form your own initial view first and then use AI assistance to check it against AI input.

The human review is always accountable for the review they post as if it was a contribution. AI usage here should be purely additive and never a replacement.
