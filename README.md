<!-- markdownlint-disable MD041 -->
<div align=center>
  <a href="https://github.com/kemadev/go-framework/actions/workflows/go-cd.yaml"><img alt="Go - CD - main" src="https://github.com/kemadev/go-framework/actions/workflows/go-cd.yaml/badge.svg?branch=main&event=push"></a>
  <a href="https://github.com/kemadev/go-framework/actions/workflows/go-cd.yaml"><img alt="Go - CD - next" src="https://github.com/kemadev/go-framework/actions/workflows/go-cd.yaml/badge.svg?branch=next&event=push"></a>
  <a href="https://github.com/kemadev/go-framework/actions/workflows/go-release.yaml"><img alt="Go - Release" src="https://github.com/kemadev/go-framework/actions/workflows/go-release.yaml/badge.svg?branch=main&event=push"></a>
  <a href="https://github.com/kemadev/go-framework/actions/workflows/docker-ci.yaml"><img alt="Docker - CI" src="https://github.com/kemadev/go-framework/actions/workflows/docker-ci.yaml/badge.svg?branch=main&event=schedule"></a>
  <a href="https://github.com/kemadev/go-framework/actions/workflows/global-ci.yaml"><img alt="Global - CI" src="https://github.com/kemadev/go-framework/actions/workflows/global-ci.yaml/badge.svg?branch=main&event=schedule"></a>
  <a href="https://github.com/kemadev/go-framework/actions/workflows/go-ci.yaml"><img alt="Go - CI" src="https://github.com/kemadev/go-framework/actions/workflows/go-ci.yaml/badge.svg?branch=main&event=schedule"></a>
  <a href="https://github.com/kemadev/go-framework/actions/workflows/repo-template-stale.yaml"><img alt="Housekeeping - Repo template staleness check" src="https://github.com/kemadev/go-framework/actions/workflows/repo-template-stale.yaml/badge.svg?branch=main&event=schedule"></a>
  <a href="https://github.com/kemadev/go-framework/actions/workflows/branch-stale.yaml"><img alt="Housekeeping - Branches staleness check" src="https://github.com/kemadev/go-framework/actions/workflows/branch-stale.yaml/badge.svg?branch=main&event=schedule"></a>
  <a href="https://github.com/kemadev/go-framework/actions/workflows/issue-pr-stale.yaml"><img alt="Housekeeping - Issues & PRs staleness housekeeping" src="https://github.com/kemadev/go-framework/actions/workflows/issue-pr-stale.yaml/badge.svg?branch=main&event=schedule"></a>
  <a href="https://github.com/kemadev/go-framework/actions/workflows/markdown-ci.yaml"><img alt="Markdown - CI" src="https://github.com/kemadev/go-framework/actions/workflows/markdown-ci.yaml/badge.svg?branch=main&event=schedule"></a>
  <a href="https://github.com/kemadev/go-framework/actions/workflows/shell-ci.yaml"><img alt="Shell - CI" src="https://github.com/kemadev/go-framework/actions/workflows/shell-ci.yaml/badge.svg?branch=main&event=schedule"></a>
  <a href="https://github.com/kemadev/go-framework/actions/workflows/workflow-action-ci.yaml"><img alt="Workflows - CI" src="https://github.com/kemadev/go-framework/actions/workflows/workflow-action-ci.yaml/badge.svg?branch=main&event=schedule"></a>
</div>

# go-framework

<!-- Brief description of the project -->

## Welcome

<!-- Project presentation, motivation, and main features -->

## Demo

<!-- If applicable, project demo (video, screenshots, asciicinema, ...) -->

## Getting Started

<!-- Basic usage and main commands -->

## Contributing

### Global Guidelines

- You can find contributing guidelines in [CONTRIBUTING.md](CONTRIBUTING.md)
- Feeling like something could be improved? Let's do it together! From code to documentation, services to use, or linter rules, everything is discussable and improvable, make a PR! Let's move forward together.

### Concepts

- A repository such as this one is representing a project
- A project is basically an application
- An application is a microservice that works with other microservices to to achieve project's goals
- Microservices are small, loosely coupled, and independently deployable and scalable
- Each microservice should be agnostic of it downstreams. However, it should expose a clear and well-defined API to its downstreams for them to consume (that is, the microservice itself uses its upstreams' API)

### Documentation

- Go code documentation is accessible thanks to [pkgsite](https://pkg.go.dev/golang.org/x/pkgsite/cmd/pkgsite)! Just run `pkgsite` in the root of the project to serve its documentation
- Global project documentation is available in the [doc](./doc) directory

### Development Guidelines and Conventions

- All major directories contain a `PURPOSE.md` file with a brief description of directory's content and instructions on how to use it
- You are encouraged to create your own `README.md` files in subdirectories to provide project-specific instructions, and to document more-widely scoped topics in [doc](./doc) directory
- Code sharing is encouraged, such code should be placed in [pkg](pkg) and [internal/pkg](internal/pkg) directories, as-per Go's conventions
- Importing other applications libraries and packages is encouraged, following code sharing encouragement
- First class code documentation (following [Go doc comment guidelines](https://go.dev/doc/comment)) as well as project documentation is encouraged
- Following [Learn Go with tests](https://github.com/quii/learn-go-with-tests) is encouraged
- Following [Effective Go](https://go.dev/doc/effective_go) and [Google's styleguide](https://google.github.io/styleguide/go/) is encouraged
- Following [locality of behaviour](https://htmx.org/essays/locality-of-behaviour/) and [principle of least astonishment](https://en.wikipedia.org/wiki/Principle_of_least_astonishment) is encouraged
- Checking [Go proverbs](https://go-proverbs.github.io/) is encouraged
- Variables, functions, methods, ... should be named in a short and descriptive way

### Project development

- Deploy `dev` stacks of [the layers](./deploy), in ascending order of layer number (`XX-` prefix)
- Make changes and enjoy hot reloading!
- Once you are happy with your changes, create a PR!

#### Prerequisites

- [Docker](https://github.com/docker/cli) & [Docker Compose](https://github.com/docker/compose) to run applications in containers. You should configure your credentials store and credential helpers for Docker to work with your container registry
- [Go](https://github.com/golang/go) to install applications dependencies as needed

#### Running the project

- Common tasks such as running, testing, creating new IaC components, updating Cloud resources, ... are done by using [kemutil](https://github.com/kemadev/ci-cd/tree/main/tool/kemutil). You are encouraged to install and use it!

#### Debugging

- Debugger support is available in VSCode, using [vscode-go](https://github.com/golang/vscode-go/wiki/debugging) extension. A task is available in [.vscode/tasks.json](.vscode/tasks.json) to run the debugger.
  Please note that you first need to run the application through docker compose with `debug` profile, e.g. via `kemutil`

#### CI / CD

##### Locally

- CI pipelines can be mimicked locally using `ci-cd` image, mounting project's directory as a volume in `/src`, and running the same commands as in the CI pipeline
- That is, you can run the following command to run the whole CI pipeline locally:

  ```bash
  kemutil ci [--fix] [--hot] ci
  ```

- When using `--hot`, your need to export `GIT_TOKEN` environment variable to propagate your git credentials to the container, so that it can fetch private dependencies. This is typically done by running:

  ```bash
  export GIT_TOKEN=$(gh auth token)
  ```

- Other commands are available, feel free to run `kemutil help` to see the list of available commands and their usage

##### False positives

- CI Pipelines can sometime report false positives. Here is what you can do to remediate (be as specific as possible on silences to avoid shadowing real issues):
  - `golangci-lint`: Add a `nolint:<linter>[,<linter>] // <explanation>` comment. See [this doc](https://golangci-lint.run/usage/false-positives/)
  - `semgrep`: Add a `nosemgrep: <rule-id>` comment. See [this doc](https://semgrep.dev/docs/ignoring-files-folders-code)
  - `gitleaks`: Add finding fingerprint to [.gitleaksignore](config/gitleaks/.gitleaksignore). See [this doc](https://github.com/gitleaks/gitleaks#additional-configuration). Please note that **any leaked secret should be revoked and replaced as soon as possible**
  - `markdownlint`: Add a `markdownlint-disable <rule>` comment. See [this doc](https://github.com/DavidAnson/markdownlint/blob/main/README.md#configuration)
  - `shellcheck`: Add a `shellcheck disable=<rule>` comment. See [this doc](https://github.com/koalaman/shellcheck/wiki/Ignore)
  - `hadolint`: Add a `hadolint ignore=<rule>` comment. See [this doc](https://github.com/hadolint/hadolint/blob/master/README.md#ignoring-rules)
  - `actionlint`: In case of a `shellcheck` error, refer to the `shellcheck` section. Otherwise, you can pass arguments to the linting action to ignore specific rules. See [this doc](https://github.com/rhysd/actionlint/blob/main/docs/usage.md#ignore-some-errors)
  - `grype`: Add an ignore in [.grype.yaml](https://github.com/kemadev/go-framework/tree/main/config/grype/.grype.yaml). See [this doc](https://github.com/anchore/grype#specifying-matches-to-ignore).
    Please note that **any vulnerability should be remediated as soon as possible**, only add real real positives to the ignore list. Prefer deploying with a non-exploitable vulnerability reported rather than ignoring it.
