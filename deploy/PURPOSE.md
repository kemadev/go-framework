# `deploy`

## Directories in this directory

- Should be named with a prefix that defines the deployment order, such as `10-`, `20-`, ...
- Are deployed sequentially in alphabetical order, so the prefix is important
- Should be named after the component they deploy (e.g., `deploy/XX-network`, `deploy/XX-database`, `deploy/XX-app_1`, `deploy/XX-app_2`, ...)
- Starting with a `0` (eg. `00-dev`) are not deployed through CD pipelines, and are used to common manual steps (e.g. dev cluster creation)

## Files in this directory

- Are placed in subdirectories, see above
- Are related to deployments
- Should manage application deployment resources for applications, including different environments (e.g., `dev`, `next`, ...)
- Should name their projects according to the URL of their directory, replacing non-alphanumeric characters with `-` (e.g., `github-com-username-repo-deploy-XX-app_1` for `deploy/XX-app_1`).
- Should name their stacks according to the environment they deploy to (e.g., `next`, `main`, ...)
- Should implement GitOps best practices
- Should be moved to a separate repository if it can be reused across multiple projects, e.g. `30-queue` implements a global service used by may applications, and the current project is not meant to be a queue service
- Should be as simple as possible, making it possible to use stack references to share resources across projects
- Should manage deployment resources for applications
