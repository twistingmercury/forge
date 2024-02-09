# Project Forge

Project Forge is an application that allows users to create new Go projects with basic boilerplate in place. It is built
using [Go](https://go.dev/) and is currently in development. I may be expanded in the future to support projects using
other languages.

## Installation

TBD

## Usage

```bash
$ ~ forge --projectDir myProject -modulePath github.com/username/myProject
```

### Flags

| Short Option | Long Option  | Description                                                               |
|--------------|--------------|---------------------------------------------------------------------------|
| -m           | --modulePath | :exclamation: Required: The module, or root package name for the project. |
| -p           | --projectDir | :exclamation: Required: The root directory name for the project.          |
| -v           | --version    | Show the current version of twistingmercury/forge.                        |
| -h           | --help       | Show help for twistingmercury/forge.                                      |

