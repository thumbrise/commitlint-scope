# commitlint-scope

[![CI](https://github.com/thumbrise/commitlint-scope/actions/workflows/ci.yml/badge.svg)](https://github.com/thumbrise/commitlint-scope/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/thumbrise/commitlint-scope.svg)](https://pkg.go.dev/github.com/thumbrise/commitlint-scope)
[![Latest Release](https://img.shields.io/github/v/release/thumbrise/commitlint-scope?label=latest&color=blue)](https://github.com/thumbrise/commitlint-scope/releases)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](/LICENSE)
[![Coverage Status](https://coveralls.io/repos/github/thumbrise/commitlint-scope/badge.svg?branch=main)](https://coveralls.io/github/thumbrise/commitlint-scope?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/thumbrise/commitlint-scope)](https://goreportcard.com/report/github.com/thumbrise/commitlint-scope)
[![Docker Image Size](https://img.shields.io/docker/image-size/thumbrise/commitlint-scope/latest-alpine?label=image%20size)](https://hub.docker.com/r/thumbrise/commitlint-scope)

Git commit conventional scope linter.

## Table of contents

<!-- TOC -->
* [commitlint-scope](#commitlint-scope)
  * [Table of contents](#table-of-contents)
  * [Purpose](#purpose)
  * [Quick start](#quick-start)
    * [Install](#install)
      * [Docker Hub](#docker-hub)
      * [Binaries](#binaries)
    * [Run](#run)
      * [Binary](#binary)
      * [Docker](#docker)
  * [Configuration file](#configuration-file)
    * [Init](#init)
    * [Overview](#overview)
  * [Zero Configuration](#zero-configuration)
  * [CI](#ci)
    * [GitHub Actions](#github-actions)
    * [GitLab CI](#gitlab-ci)
    * [Bitbucket Pipelines](#bitbucket-pipelines)
    * [Important Notes for All CI Systems](#important-notes-for-all-ci-systems)
  * [JSON schema](#json-schema)
  * [License](#license)
<!-- TOC -->

## Purpose

Lint changed files against configured scope. Similar to CODEOWNERS.
Useful if your dev-flow requires strict scoped file changes control over CI process.

## Quick start

### Install

#### Docker Hub

```shell
docker pull thumbrise/commitlint-scope:latest-alpine
```

#### Binaries

https://github.com/thumbrise/commitlint-scope/releases

<details>
<summary>
No CGO
</summary>
Binaries are statically compiled (CGO_ENABLED=0) and work without external dependencies on Linux, macOS, and Windows.
</details>

### Run

Go to your repository root.

```shell
cd /path/to/your/repo/
```

#### Binary

```shell
commitlint-scope run --from main --to feature
```

#### Docker

```bash
docker run --rm -v "$(pwd):/repo" -w /repo thumbrise/commitlint-scope:latest-alpine run --from main --to feature --verbose
```

## Configuration file

### Init

Generate .commitlint-scope.yml file.

```shell
commitlint-scope init
```

### Overview

```yaml
#$schema: https://github.com/thumbrise/commitlint-scope/blob/main/docs/schema/config.v3.json

# Scope parsing customization. Not required, if you follow common conventional header. In example: 'type!(scope): subject'
#scopeRegex: ^[a-z]+(?:\((?P<scope>[^)]+)\))?!?:\s

# Patterns list: each item specifies a list of scopes and the corresponding file glob patterns.
patterns:
  - scopes: ["auth"]
    files: ["services/auth/**"]

  - scopes: ["migrations", "sql"]
    files: ["database/migrations/*.sql"]

  - scopes: ["frontend", "assets"]
    files: ["**/assets/**", "**/frontend/**"]

  - scopes: ["docs", "md"]
    files: ["**/*.md"]

  - scopes: ["some.dot.scope", "any-anotherscope"]
    files: ["**/rail.v1.json"]
```

## Zero Configuration

Without a configuration file, each scope is used as a glob pattern by appending /**.
For example, a commit with scope auth will check if any changed file matches auth/**.
This behaviour keeps things simple for repositories where directory names mirror commit scopes.

For example:

```text
feat(auth): Some subject
```

Linter will compare changed files against glob pattern `auth/**`.


## CI

### GitHub Actions

```yaml
name: Lint Commit Scopes
on:
  pull_request:

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository (full history)
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Lint commit scopes
        run: |
          docker run --rm \
            -v "${{ github.workspace }}:/repo" \
            -w /repo \
            thumbrise/commitlint-scope:latest-alpine \
            run \
            --from ${{ github.event.pull_request.base.sha }} \
            --to ${{ github.event.pull_request.head.sha }} \
            --verbose
```

### GitLab CI

```yaml
commitlint:
  image: thumbrise/commitlint-scope:latest-alpine
  stage: test
  variables:
    GIT_DEPTH: 0
  script:
    - commitlint-scope run --from $CI_MERGE_REQUEST_DIFF_BASE_SHA --to $CI_COMMIT_SHA
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
```

### Bitbucket Pipelines

```yaml
pipelines:
  pull-requests:
    '**':
      - step:
          name: Lint commit scopes
          image: thumbrise/commitlint-scope:latest-alpine # For production - use specific version
          script:
            - git fetch --unshallow || true
            - commitlint-scope run --from $BITBUCKET_PR_DESTINATION_BRANCH --to $BITBUCKET_COMMIT
```

### Important Notes for All CI Systems

- Always do a full clone (`fetch-depth: 0` or equivalent) so `git rev-list` can find all commits.
- For shallow clones, use `git fetch --unshallow` or ensure the required commit range is available, otherwise `git rev-list` will fail.
- If your CI doesn't provide base/target commit variables, specify the SHAs manually.
- The image already includes safe directory configuration (`safe.directory '*'`) for smooth operation.
---

## JSON schema

https://github.com/thumbrise/commitlint-scope/blob/main/docs/schema/config.v3.json

## License

Apache 2.0
