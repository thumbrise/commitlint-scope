# commitlint-scope

[![CI](https://github.com/thumbrise/commitlint-scope/actions/workflows/ci.yml/badge.svg)](https://github.com/thumbrise/commitlint-scope/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/thumbrise/commitlint-scope.svg)](https://pkg.go.dev/github.com/thumbrise/commitlint-scope)
[![Latest Release](https://img.shields.io/github/v/release/thumbrise/commitlint-scope?label=latest&color=blue)](https://github.com/thumbrise/commitlint-scope/releases)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](/LICENSE)
[![Coverage Status](https://coveralls.io/repos/github/thumbrise/commitlint-scope/badge.svg?branch=main)](https://coveralls.io/github/thumbrise/commitlint-scope?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/thumbrise/commitlint-scope)](https://goreportcard.com/report/github.com/thumbrise/commitlint-scope)

Git commit conventional scope linter.

## Table of contents

<!-- TOC -->
* [commitlint-scope](#commitlint-scope)
  * [Table of contents](#table-of-contents)
  * [Purpose](#purpose)
  * [Quick start](#quick-start)
    * [Install](#install)
    * [Run](#run)
  * [Configuration file](#configuration-file)
  * [Zero Configuration](#zero-configuration)
  * [CI](#ci)
  * [JSON schema](#json-schema)
  * [License](#license)
<!-- TOC -->

## Purpose

Lint changed files against configured scope. Similar to CODEOWNERS.
Useful if your dev-flow requires strict scoped file changes control over CI process.

## Quick start

### Install

https://github.com/thumbrise/commitlint-scope/releases

### Run

```shell
commitlint-scope --from main --to feature
```

## Configuration file

**Init**

Generate .commitlint-scope.yml file.

```shell
commitlint-scope init
```

**Overview**

```yaml
#$schema: https://github.com/thumbrise/commitlint-scope/blob/main/docs/schema/config.json

# Scope parsing customization. Not required, if you follow common conventional header. In example: 'type!(scope): subject'
#scopeRegex: ^[a-z]+(?:\((?P<scope>[^)]+)\))?!?:\s

# Patterns map: each key is a scope name, value is a list of glob patterns that match files belonging to that scope.
patterns:
  "auth": [ "services/auth/**" ]
  "migrations": [ "database/migrations/*.sql" ]
  "frontend": [ "**/assets/**", "**/frontend/**" ]
  "docs": [ "**/*.md" ]
```

## Zero Configuration

Without configuration file - scopes are treated as file-glob matches, relative to repository root.

```text
feat(auth): Some subject
```

Linter will compare changed files against glob pattern `auth/**`


## CI

TODO: Add high-DX examples

---

## JSON schema

https://github.com/thumbrise/commitlint-scope/blob/main/docs/schema/config.json

## License

Apache 2.0
