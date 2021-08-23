# Autopilot Resource

Bootstrap and manage your Concourse pipelines in a GitOps fashion:<br />

## Source Configuration

- `uri`: _Required._ The location of the [**autopilot repo**](#autopilot-repo).
- `branch`: The branch to track. If unset, the default branch is used.
- `config`: _Required_ A key:value pair of **team** : [**config file**](#config-file)
  - key: The Concourse team to which you want to set the pipelines listed in the **config file**
  - value: The path to a **config file**. It must be relative to the root of the **autopilot repo**<br />
    [_You can use wildcards to grant your autopilot pipeline superpowers_](#autopilot-superpowers)

## Behaviour

### `check`

Triggers a new version whenever any of these happens:
- The **autopilot repo** branch you are tracking gets a new commit
- You update the **autopilot pipeline** using `fly set-pipeline`
- A new version of the `concourse-autopilot-resource` is published

This ensures that the pipelines are always fresh and immediately react
to any changes made to the configuration or the autopilot repo.

### `get`

Generates a `pipeline.yml` with a set of `set_pipeline` steps which will keep your Concourse instance in sync with the **autopilot repo** you are tracking.

The recommended way to set this `pipeline.yml` is to use `set_pipeline: self` in your **autopilot pipeline**.

Your **autopilot pipeline** should not contain any custom resource, job, etc. apart from the ones listed in the example below _since they will be overwritten by the `set_pipeline: self` step._

### `put`

NO-OP

## Examples

* [Example autopilot repository](https://github.com/efejjota/concourse-autopilot-example): _ready-to-run experiment to test autopilot resource_
* Example autopilot pipeline:
```yaml
---
resource_types:
- name: autopilot
  type: registry-image
  source:
    repository: efejjota/concourse-autopilot-resource

resources:
- name: autopilot
  type: autopilot
  source:
     uri: https://your-autopilot-repo-url
     config:
     - team1: autopilot/team1/*.yml
     - team2: autopilot/team2/*.yml
     - main: autopilot/some-pipelines.yml

jobs:
- name: sync-pipelines
  plan:
  - get: autopilot
    trigger: true

  - set_pipeline: self
    file: autopilot/pipeline.yml
```

## Additional Info

### Autopilot Superpowers

By using a wildcard you can easily track entire folders, automatically detecting the creation, removal or rearrangement of **config files**.

Also, by mapping different folders to different Concourse teams you can easily move pipelines from one team to another by simply moving config files around.

### Autopilot Repo

`TO DO`: Describe the the autopilot repo concept

Look at the [concourse-autopilot-example](https://github.com/efejjota/concourse-autopilot-example) repo for a reference

### Config File

`TO DO`: Describe autopilot format

Review files in the [concourse-autopilot-example/autopilot/](https://github.com/efejjota/concourse-autopilot-example/tree/main/autopilot) folder for a reference
