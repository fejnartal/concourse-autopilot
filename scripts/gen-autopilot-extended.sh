#!/bin/bash
set -euo pipefail

target_path="$1"
target_abs="$(realpath "${target_path}")"
target_folder="$(dirname "${target_abs}")"

repodir="$(git -C "${target_folder}" rev-parse --show-toplevel)"
target_rel="$(realpath --relative-to="${repodir}" "${target_abs}")"

{
  config="$(yq eval -o=json "${target_path}" 2>/dev/null | jq -rc '.config')"
  vcs_source="$(yq eval -o=json "${target_path}" 2>/dev/null | jq -rc '.vcs_resource.source')"
  vcs_type="$(yq eval -o=json "${target_path}" 2>/dev/null | jq -rc '.vcs_resource.type')"
  vcs_registry_image="$(yq eval -o=json "${target_path}" 2>/dev/null | jq -rc '.vcs_resource["registry-image"] // empty')"
} || {
  config="$(yq read --tojson "${target_path}" | jq -rc '.config')"
  vcs_source="$(yq read --tojson "${target_path}" | jq -rc '.vcs_resource.source')"
  vcs_type="$(yq read --tojson "${target_path}" 2>/dev/null | jq -rc '.vcs_resource.type')"
  vcs_registry_image="$(yq read --tojson "${target_path}" | jq -rc '.vcs_resource["registry-image"] // empty')"
}

generated_repositories=''
generated_jobs=''
generated_groups=''

for config_wildcard in $(echo "${config}" | jq -r '.[][]'); do
for config_file in $(ls ${repodir}/${config_wildcard} 2> /dev/null); do
    {
      autopilot_config="$(yq eval -o=json "${config_file}" 2>/dev/null | jq -rc '.autopilot_config')"
    } || {
      autopilot_config="$(yq read --tojson "${config_file}" | jq -rc '.autopilot_config')"
    }
    generated_groups+="
- name: ${autopilot_config:?}
  jobs:
"

    # This way of looping through a json array is documented by Ruben Koster:
    # https://www.starkandwayne.com/blog/bash-for-loop-over-json-array-using-jq/
    {
      repositories="$(yq eval -o=json "${config_file}" 2>/dev/null | jq -rc '.repositories[] | @base64')"
    } || {
      repositories="$(yq read --tojson "${config_file}" | jq -rc '.repositories[] | @base64')"
    }
    for repository in ${repositories}; do
      _jq() {
      echo "${repository}" | base64 --decode | jq -r "${1}"
      }

      generated_repositories+="$(_jq '"
- name: \(.name)
  type: git
  source: \(.)
"')"
    done

    {
      pipelines="$(yq eval -o=json "${config_file}" 2>/dev/null | jq -r '.pipelines[] | @base64')"
    } || {
      pipelines="$(yq read --tojson "${config_file}" | jq -r '.pipelines[] | @base64')"
    }
    team="$(echo "${config}" | jq -r '.[] | to_entries[] | select(.value=="'"$config_wildcard"'").key')"

    for pipeline in ${pipelines}; do
      _jq() {
      echo "${pipeline}" | base64 --decode | jq -r "${1}" --arg team "${team}"
      }

      generated_jobs+="$(_jq '"
- name: set-\(.name)
  plan:
  - get: repository
    passed: [sync-pipelines]
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

generated_resource_types=""
if [[ -n "${vcs_registry_image}" ]];
then
  generated_resource_types="
resource_types:
- name: ${vcs_type}
  type: registry-image
  source: ${vcs_registry_image}
"
fi

generated_manifest="
---
groups:
- name: autopilot
  jobs:
  - sync-pipelines

${generated_groups}

${generated_resource_types}

resources:
- name: repository
  type: ${vcs_type}
  source: ${vcs_source}

${generated_repositories}

jobs:
- name: sync-pipelines
  plan:
  - get: repository
    trigger: true
  - task: regenerate-pipeline
    config:
      platform: linux
      image_resource:
        type: registry-image
        source: { repository: efejjota/concourse-autopilot-resource }
      inputs:
      - name: repository
      outputs:
      - name: regenerated
      run:
        path: /bin/bash
        args:
        - \"-e\"
        - \"-c\"
        - |
          /opt/resource/scripts/gen-autopilot-extended.sh \"repository/${target_rel}\" > regenerated/pipeline.yml

  - set_pipeline: self
    file: regenerated/pipeline.yml

${generated_jobs}
"

{
  echo "${generated_manifest}" | yq eval -o=json 2>/dev/null | jq '.resources |= unique' | yq eval --prettyPrint 2>/dev/null
} || {
  echo "${generated_manifest}" | yq read --tojson - | jq '.resources |= unique' | yq read --prettyPrint -
}
