# Autopilot Resource

Bootstrap and manage your Concourse pipelines in a GitOps fashion

## Quickstart

### Download the autopilot cli for your corresponding OS
https://github.com/efejjota/concourse-autopilot/releases/

_If you are running on Linux or MacOS, remember to make it executable_
```bash
chmod +x autopilot
```

### Clone the repository you want use with autopilot
```bash
git clone https://github.com/efejjota/concourse-autopilot-example.git my-repo/
```

### Bootstrap required files for autopilot to work
_You need to perform this step only once per repository_
```bash
autopilot init my-repo/
```

### Synchronize your Concourse with the repository
```bash
autopilot fly my-repo/
```
