[private]
default:
  @just --list profile

list:
  #!/bin/sh

  set -euo pipefail

  export config="${DOCKER_RESTIC_CONFIG_DIR}/profiles.toml"

  yq -r 'keys | .[]' "${config}"

expand profile *args:
  #!/bin/sh

  set -euo pipefail

  export config="${DOCKER_RESTIC_CONFIG_DIR}/profiles.toml"
  export PROFILE="{{profile}}"
  export args="{{args}}"

  for command in "${PROFILE}" ${args}; do
    export query="${query:+${query}.}${command}"

    flags=$(yq -r '
      .. | select(strenv(query) == (path | join(".")))
      | select(kind == "map")
      | to_entries
      | map(select((.value | kind) == "scalar" and (.value // false))
        | "--" + .key + ((.value | select(tag != "!!bool") | " " + tostring) // ""))
      | join(" ")
    ' "${config}")

    [ "${query}" = "${command}" ] || printf '%s ' "${command}"
    printf '%s ' "${flags}"
  done | envsubst
  printf '\n'

get profile path:
  #!/bin/sh

  set -euo pipefail

  export config="${DOCKER_RESTIC_CONFIG_DIR}/profiles.toml"
  export PROFILE="{{profile}}"
  export path="{{path}}"

  yq -r '.[strenv(PROFILE)][strenv(path)] // ""' "${config}" | envsubst
