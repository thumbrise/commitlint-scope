# commitlint-scope

A new CLI generated with oclif

[![oclif](https://img.shields.io/badge/cli-oclif-brightgreen.svg)](https://oclif.io)
[![Version](https://img.shields.io/npm/v/commitlint-scope.svg)](https://npmjs.org/package/commitlint-scope)
[![Downloads/week](https://img.shields.io/npm/dw/commitlint-scope.svg)](https://npmjs.org/package/commitlint-scope)

<!-- toc -->

- [Usage](#usage)
- [Commands](#commands)
<!-- tocstop -->

# Usage

<!-- usage -->

```sh-session
$ npm install -g commitlint-scope
$ commitlint-scope COMMAND
running command...
$ commitlint-scope (--version)
commitlint-scope/0.0.0 darwin-arm64 node-v24.15.0
$ commitlint-scope --help [COMMAND]
USAGE
  $ commitlint-scope COMMAND
...
```

<!-- usagestop -->

# Commands

<!-- commands -->

- [`commitlint-scope hello PERSON`](#commitlint-scope-hello-person)
- [`commitlint-scope hello world`](#commitlint-scope-hello-world)
- [`commitlint-scope help [COMMAND]`](#commitlint-scope-help-command)
- [`commitlint-scope plugins`](#commitlint-scope-plugins)
- [`commitlint-scope plugins add PLUGIN`](#commitlint-scope-plugins-add-plugin)
- [`commitlint-scope plugins:inspect PLUGIN...`](#commitlint-scope-pluginsinspect-plugin)
- [`commitlint-scope plugins install PLUGIN`](#commitlint-scope-plugins-install-plugin)
- [`commitlint-scope plugins link PATH`](#commitlint-scope-plugins-link-path)
- [`commitlint-scope plugins remove [PLUGIN]`](#commitlint-scope-plugins-remove-plugin)
- [`commitlint-scope plugins reset`](#commitlint-scope-plugins-reset)
- [`commitlint-scope plugins uninstall [PLUGIN]`](#commitlint-scope-plugins-uninstall-plugin)
- [`commitlint-scope plugins unlink [PLUGIN]`](#commitlint-scope-plugins-unlink-plugin)
- [`commitlint-scope plugins update`](#commitlint-scope-plugins-update)

## `commitlint-scope hello PERSON`

Say hello

```
USAGE
  $ commitlint-scope hello PERSON -f <value>

ARGUMENTS
  PERSON  Person to say hello to

FLAGS
  -f, --from=<value>  (required) Who is saying hello

DESCRIPTION
  Say hello

EXAMPLES
  $ commitlint-scope hello friend --from oclif
  hello friend from oclif! (./src/commands/hello/index.ts)
```

_See code: [src/commands/hello/index.ts](https://github.com/thumbrise/commitlint-scope/blob/v0.0.0/src/commands/hello/index.ts)_

## `commitlint-scope hello world`

Say hello world

```
USAGE
  $ commitlint-scope hello world

DESCRIPTION
  Say hello world

EXAMPLES
  $ commitlint-scope hello world
  hello world! (./src/commands/hello/world.ts)
```

_See code: [src/commands/hello/world.ts](https://github.com/thumbrise/commitlint-scope/blob/v0.0.0/src/commands/hello/world.ts)_

## `commitlint-scope help [COMMAND]`

Display help for commitlint-scope.

```
USAGE
  $ commitlint-scope help [COMMAND...] [-n]

ARGUMENTS
  [COMMAND...]  Command to show help for.

FLAGS
  -n, --nested-commands  Include all nested commands in the output.

DESCRIPTION
  Display help for commitlint-scope.
```

_See code: [@oclif/plugin-help](https://github.com/oclif/plugin-help/blob/6.2.49/src/commands/help.ts)_

## `commitlint-scope plugins`

List installed plugins.

```
USAGE
  $ commitlint-scope plugins [--json] [--core]

FLAGS
  --core  Show core plugins.

GLOBAL FLAGS
  --json  Format output as json.

DESCRIPTION
  List installed plugins.

EXAMPLES
  $ commitlint-scope plugins
```

_See code: [@oclif/plugin-plugins](https://github.com/oclif/plugin-plugins/blob/5.4.69/src/commands/plugins/index.ts)_

## `commitlint-scope plugins add PLUGIN`

Installs a plugin into commitlint-scope.

```
USAGE
  $ commitlint-scope plugins add PLUGIN... [--json] [-f] [-h] [-s | -v]

ARGUMENTS
  PLUGIN...  Plugin to install.

FLAGS
  -f, --force    Force npm to fetch remote resources even if a local copy exists on disk.
  -h, --help     Show CLI help.
  -s, --silent   Silences npm output.
  -v, --verbose  Show verbose npm output.

GLOBAL FLAGS
  --json  Format output as json.

DESCRIPTION
  Installs a plugin into commitlint-scope.

  Uses npm to install plugins.

  Installation of a user-installed plugin will override a core plugin.

  Use the COMMITLINT_SCOPE_NPM_LOG_LEVEL environment variable to set the npm loglevel.
  Use the COMMITLINT_SCOPE_NPM_REGISTRY environment variable to set the npm registry.

ALIASES
  $ commitlint-scope plugins add

EXAMPLES
  Install a plugin from npm registry.

    $ commitlint-scope plugins add myplugin

  Install a plugin from a github url.

    $ commitlint-scope plugins add https://github.com/someuser/someplugin

  Install a plugin from a github slug.

    $ commitlint-scope plugins add someuser/someplugin
```

## `commitlint-scope plugins:inspect PLUGIN...`

Displays installation properties of a plugin.

```
USAGE
  $ commitlint-scope plugins inspect PLUGIN...

ARGUMENTS
  PLUGIN...  [default: .] Plugin to inspect.

FLAGS
  -h, --help     Show CLI help.
  -v, --verbose

GLOBAL FLAGS
  --json  Format output as json.

DESCRIPTION
  Displays installation properties of a plugin.

EXAMPLES
  $ commitlint-scope plugins inspect myplugin
```

_See code: [@oclif/plugin-plugins](https://github.com/oclif/plugin-plugins/blob/5.4.69/src/commands/plugins/inspect.ts)_

## `commitlint-scope plugins install PLUGIN`

Installs a plugin into commitlint-scope.

```
USAGE
  $ commitlint-scope plugins install PLUGIN... [--json] [-f] [-h] [-s | -v]

ARGUMENTS
  PLUGIN...  Plugin to install.

FLAGS
  -f, --force    Force npm to fetch remote resources even if a local copy exists on disk.
  -h, --help     Show CLI help.
  -s, --silent   Silences npm output.
  -v, --verbose  Show verbose npm output.

GLOBAL FLAGS
  --json  Format output as json.

DESCRIPTION
  Installs a plugin into commitlint-scope.

  Uses npm to install plugins.

  Installation of a user-installed plugin will override a core plugin.

  Use the COMMITLINT_SCOPE_NPM_LOG_LEVEL environment variable to set the npm loglevel.
  Use the COMMITLINT_SCOPE_NPM_REGISTRY environment variable to set the npm registry.

ALIASES
  $ commitlint-scope plugins add

EXAMPLES
  Install a plugin from npm registry.

    $ commitlint-scope plugins install myplugin

  Install a plugin from a github url.

    $ commitlint-scope plugins install https://github.com/someuser/someplugin

  Install a plugin from a github slug.

    $ commitlint-scope plugins install someuser/someplugin
```

_See code: [@oclif/plugin-plugins](https://github.com/oclif/plugin-plugins/blob/5.4.69/src/commands/plugins/install.ts)_

## `commitlint-scope plugins link PATH`

Links a plugin into the CLI for development.

```
USAGE
  $ commitlint-scope plugins link PATH [-h] [--install] [-v]

ARGUMENTS
  PATH  [default: .] path to plugin

FLAGS
  -h, --help          Show CLI help.
  -v, --verbose
      --[no-]install  Install dependencies after linking the plugin.

DESCRIPTION
  Links a plugin into the CLI for development.

  Installation of a linked plugin will override a user-installed or core plugin.

  e.g. If you have a user-installed or core plugin that has a 'hello' command, installing a linked plugin with a 'hello'
  command will override the user-installed or core plugin implementation. This is useful for development work.


EXAMPLES
  $ commitlint-scope plugins link myplugin
```

_See code: [@oclif/plugin-plugins](https://github.com/oclif/plugin-plugins/blob/5.4.69/src/commands/plugins/link.ts)_

## `commitlint-scope plugins remove [PLUGIN]`

Removes a plugin from the CLI.

```
USAGE
  $ commitlint-scope plugins remove [PLUGIN...] [-h] [-v]

ARGUMENTS
  [PLUGIN...]  plugin to uninstall

FLAGS
  -h, --help     Show CLI help.
  -v, --verbose

DESCRIPTION
  Removes a plugin from the CLI.

ALIASES
  $ commitlint-scope plugins unlink
  $ commitlint-scope plugins remove

EXAMPLES
  $ commitlint-scope plugins remove myplugin
```

## `commitlint-scope plugins reset`

Remove all user-installed and linked plugins.

```
USAGE
  $ commitlint-scope plugins reset [--hard] [--reinstall]

FLAGS
  --hard       Delete node_modules and package manager related files in addition to uninstalling plugins.
  --reinstall  Reinstall all plugins after uninstalling.
```

_See code: [@oclif/plugin-plugins](https://github.com/oclif/plugin-plugins/blob/5.4.69/src/commands/plugins/reset.ts)_

## `commitlint-scope plugins uninstall [PLUGIN]`

Removes a plugin from the CLI.

```
USAGE
  $ commitlint-scope plugins uninstall [PLUGIN...] [-h] [-v]

ARGUMENTS
  [PLUGIN...]  plugin to uninstall

FLAGS
  -h, --help     Show CLI help.
  -v, --verbose

DESCRIPTION
  Removes a plugin from the CLI.

ALIASES
  $ commitlint-scope plugins unlink
  $ commitlint-scope plugins remove

EXAMPLES
  $ commitlint-scope plugins uninstall myplugin
```

_See code: [@oclif/plugin-plugins](https://github.com/oclif/plugin-plugins/blob/5.4.69/src/commands/plugins/uninstall.ts)_

## `commitlint-scope plugins unlink [PLUGIN]`

Removes a plugin from the CLI.

```
USAGE
  $ commitlint-scope plugins unlink [PLUGIN...] [-h] [-v]

ARGUMENTS
  [PLUGIN...]  plugin to uninstall

FLAGS
  -h, --help     Show CLI help.
  -v, --verbose

DESCRIPTION
  Removes a plugin from the CLI.

ALIASES
  $ commitlint-scope plugins unlink
  $ commitlint-scope plugins remove

EXAMPLES
  $ commitlint-scope plugins unlink myplugin
```

## `commitlint-scope plugins update`

Update installed plugins.

```
USAGE
  $ commitlint-scope plugins update [-h] [-v]

FLAGS
  -h, --help     Show CLI help.
  -v, --verbose

DESCRIPTION
  Update installed plugins.
```

_See code: [@oclif/plugin-plugins](https://github.com/oclif/plugin-plugins/blob/5.4.69/src/commands/plugins/update.ts)_

<!-- commandsstop -->
