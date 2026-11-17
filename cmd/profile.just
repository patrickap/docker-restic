config_file := "${DOCKER_RESTIC_CONFIG_DIR}/profiles.toml"

[private]
default:
  @just --list profile

list:
  @set -euo pipefail \
    && yq 'keys | .[]' "{{config_file}}"

expand profile *args:
  #!/bin/sh

  set -euo pipefail

  yq -o=json '.' "{{config_file}}" \
    | jq -r --arg profile "{{profile}}" --arg words "{{args}}" '
        # scalars skips nested tables, and false is falsy so it is skipped too
        def flags: objects | to_entries[] | select(.value | scalars) | "--\(.key) \(.value)" | rtrimstr(" true");

        (.[$profile] // error("unknown profile: " + $profile)) as $node
        | [$node | flags]
          + [foreach ($words | scan("\\S+")) as $word (
              [$node, null];
              ((.[0][$word] | objects) // null) as $child | [$child // .[0], $child];
              $word, (.[1] | flags))]
        | join(" ")
      ' \
    | envsubst

get profile path:
  @set -euo pipefail \
    && yq '.{{profile}}.{{path}} // ""' "{{config_file}}" \
    | envsubst
