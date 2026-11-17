config_file := "${DOCKER_RESTIC_CONFIG_DIR}/profiles.toml"

[private]
default:
  @just --list profile

list:
  @set -euo pipefail \
    && yq -r 'keys | .[]' "{{config_file}}"

expand profile *args:
  #!/bin/sh

  set -euo pipefail

  for command in "{{profile}}" {{args}}; do
    query="${query:+${query}.}${command}"

    flags=$(yq -r '
      .. | select(strenv(query) == (path | join(".")))
      | select(kind == "map")
      | to_entries
      | map(select((.value | kind) == "scalar" and (.value // false))
        | "--" + .key + ((.value | select(tag != "!!bool") | " " + tostring) // ""))
      | join(" ")
    ' "{{config_file}}")

    [ "${query}" = "${command}" ] || printf '%s ' "${command}"
    printf '%s ' "${flags}"
  done | envsubst
  printf '\n'

get profile path:
  @set -euo pipefail \
    && yq -r '.[strenv({{profile}})][strenv({{path}})] // ""' "{{config_file}}" \
    | envsubst
