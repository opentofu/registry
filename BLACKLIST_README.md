# Version Blacklist Documentation

## Overview
The OpenTofu Registry supports blacklisting specific versions of providers and modules to prevent them from being added during automated version updates.

## Configuration
The blacklist is configured in the `versions_blacklist.json` file at the repository root.

### Structure
```json
{
  "providers": [
    {
      "namespace": "hashicorp",
      "name": "aws",
      "version": "6.1.0",
      "reason": "Critical bug - see https://github.com/hashicorp/terraform-provider-aws/issues/43213"
    }
  ],
  "modules": [
    {
      "namespace": "terraform-aws-modules",
      "name": "vpc",
      "target_system": "aws",
      "version": "5.0.0",
      "reason": "Breaking changes not compatible with our infrastructure"
    }
  ]
}
```

## How It Works
1. During the automated version bump process (runs every 15 minutes), the system checks each new version against the blacklist
2. If a version is found in the blacklist, it will be skipped and a warning will be logged
3. The blacklisted version will never be added to the registry, even if it exists in the upstream repository

## Adding a Blacklisted Version
1. Edit `versions_blacklist.json`
2. Add an entry to either the `providers` or `modules` array
3. Commit and push the changes
4. The blacklist takes effect immediately on the next version bump run

## Removing a Blacklisted Version
1. Remove the entry from `versions_blacklist.json`
2. Commit and push the changes
3. The version will be eligible for addition on the next version bump run

## Important Notes
- The blacklist only prevents NEW versions from being added
- If a blacklisted version already exists in the registry, it won't be automatically removed
- To remove an existing version, you should:
  1. Manually edit the provider/module JSON file to remove the version
  2. Add the version to `versions_blacklist.json` to prevent re-addition
  3. Submit both changes in the same PR
- The version string must match exactly (e.g., "6.1.0" not "v6.1.0")