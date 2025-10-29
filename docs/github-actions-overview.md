# GitHub Actions Overview

This document summarizes the CI/CD workflows used in the OpenTofu Registry project.

- **build.yml** → Runs unit tests, lints code, and builds the project.
- **deploy.yml** → Handles release and deployment to registry servers.
- **test.yml** → Executes integration tests with provider mock data.
- **docs.yml** → Validates documentation and triggers rebuilds for updated markdown.

## How to extend
To add a new workflow, create a YAML file under `.github/workflows` and ensure it
follows the naming convention `<task>.yml`.
