#!/bin/bash
set -euo pipefail

repodir="$1"
source="$2"
uri="$(jq -r '.uri' <(echo "${source}"))"
config="$(jq -r '.config' <(echo "${source}"))"

generated_repositories=''
generated_jobs=''
generated_groups=''

for config_wildcard in $(echo "${config}" | jq -r '.[][]'); do
for config_file in $(ls ${repodir}/${config_wildcard} 2> /dev/null); do
    autopilot_config="$(cat "${config_file}" | yq eval -o=j | jq -r '.autopilot_config')"
    generated_groups+="
- name: "${autopilot_config:?}"
  jobs:
"

    # This way of looping through a json array is documented by Ruben Koster:
    # https://www.starkandwayne.com/blog/bash-for-loop-over-json-array-using-jq/
    repositories="$(cat "${config_file}" | yq eval -o=j | jq -r '.repositories[] | @base64')"
    for repository in ${repositories}; do
      _jq() {
      echo ${repository} | base64 --decode | jq -r "${1}"
      }

      generated_repositories+="$(_jq '"
- name: \(.name)
  type: git
  source:
    uri: \(.uri)
    branch: \(.branch)
"')"
    done

    pipelines="$(cat "${config_file}" | yq eval -o=j | jq -r '.pipelines[] | @base64')"
    team="$(echo "${config}" | jq -r '.[] | to_entries[] | select(.value=="'$config_wildcard'").key')"

    for pipeline in ${pipelines}; do
      _jq() {
      echo ${pipeline} | base64 --decode | jq -r "${1}" --arg team "${team}"
      }

      generated_jobs+="$(_jq '"
- name: set-\(.name)
  plan:
  - get: autopilot
    passed: [sync-pipelines]
    trigger: true
  - get: \(.repository)
    trigger: true
  - set_pipeline: \(.name)
    team: \($team)
    file: \(.repository)/\(.manifest)
    vars: \(.vars)
"')"

      generated_groups+="$(_jq '"
  - set-\(.name)
"')"
    done
done
done

generated_manifest="
---
groups:
- name: autopilot
  jobs:
  - sync-pipelines

${generated_groups}

resource_types:
- name: autopilot
  type: registry-image
  source:
    repository: efejjota/concourse-autopilot-resource

resources:
- name: autopilot
  type: autopilot
  source:
    uri: "${uri}"
    config: "${config}"

${generated_repositories}

jobs:
- name: sync-pipelines
  plan:
  - get: autopilot
    trigger: true

  - set_pipeline: self
    file: autopilot/pipeline.yml

${generated_jobs}
"

echo "${generated_manifest}" | yq eval -o=json | jq '.resources |= unique' | yq eval --prettyPrint
