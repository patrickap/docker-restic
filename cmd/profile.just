set shell := ["/bin/sh", "-euo", "pipefail", "-uc"]

config_file := "${DOCKER_RESTIC_CONFIG_DIR}/profiles.toml"

[private]
default:
  @just --list profile

list:
  #!/usr/bin/env python3

  import os, tomllib

  with open(os.path.expandvars("{{config_file}}"), "rb") as f:
    config = tomllib.load(f)

  print("\n".join(config))

expand profile *args:
  #!/usr/bin/env python3

  import os, tomllib
  from functools import reduce

  with open(os.path.expandvars("{{config_file}}"), "rb") as f:
    config = tomllib.load(f)

  def expand(state, word):
    node, result = state
    child = node.get(word)
    child = child if isinstance(child, dict) else {}

    return (
      node.get(word, node) if child else node,
      result + [word] + [
        f"--{k}" if v is True else f"--{k} {v}"
        for k, v in child.items()
        if isinstance(v, (str, int, float, bool))
      ]
    )

  _, args = reduce(expand, "{{profile}} {{args}}".split(), (config, []))

  print(os.path.expandvars(" ".join(args[1:])))

get profile path:
  #!/usr/bin/env python3

  import os, tomllib
  from functools import reduce

  with open(os.path.expandvars("{{config_file}}"), "rb") as f:
    config = tomllib.load(f)

  def get(node, key):
    return node.get(key, {})

  value = reduce(get, "{{profile}}.{{path}}".split("."), config)

  print(os.path.expandvars(str(value) if isinstance(value, (str, int, float)) else str(value).lower() if isinstance(value, bool) else ""))
