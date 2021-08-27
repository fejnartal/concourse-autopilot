# Autopilot Resource

Bootstrap and manage your Concourse pipelines in a GitOps fashion:<br />

## Quickstart

### My repository is already configured with Autopolot

```bash
# Clone your repository. Skip this step if you already have a local copy
git clone ${your_repo_url:-https://github.com/efejjota/concourse-autopilot-example.git} my-repo
```
```bash
# `cd` into your repository folder
cd my-repo
```
```bash
# Use `fly execute` to bootstrap Autopilot in your Concourse
fly -t ${concourse_target:?}                        \
    execute                                         \
    -v config=${config_path:-.autopilot/config.yml} \
    -i repository=.                                 \
    -o regenerated=regenerated                      \
    -c <(curl -L https://raw.githubusercontent.com/efejjota/autopilot/main/tasks/gen-pipeline.yml)  \
 && sh regenerated/set-autopilot.sh
```

### It's my first time using Autopilot

`TO DO`: Describe how to setup Autopilot in your repository

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
