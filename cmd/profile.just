set shell := ["/bin/sh", "-euo", "pipefail", "-uc"]

config_file := "${DOCKER_RESTIC_CONFIG_DIR}/profiles.toml"

[private]
default:
  @just --list profile

list:
  #!/usr/bin/env python3

  import os, tomllib, shlex

  with open(os.path.expandvars("{{config_file}}"), "rb") as f:
    config = tomllib.load(f)

  print("\n".join(shlex.quote(str(key)) for key in config))

expand profile *args:
  #!/usr/bin/env python3

  import os, tomllib, shlex
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
        arg
        for k, v in child.items()
        if not k.startswith("_") and isinstance(v, (str, int, float, bool))
        for arg in (
          [f"--{k}"] if v is True else [f"--{k}", os.path.expandvars(str(v))]
        )
      ]
    )

  _, args = reduce(expand, shlex.split("{{profile}} {{args}}"), (config, []))
  result = shlex.join(args[1:])
  
  print(result)

get profile path:
  #!/usr/bin/env python3

  import os, tomllib, shlex
  from functools import reduce

  with open(os.path.expandvars("{{config_file}}"), "rb") as f:
    config = tomllib.load(f)

  def get(node, key):
    return node.get(key, {})

  def render(value, key=""):
    if isinstance(value, dict):
      for k, v in value.items():
        yield from render(v, f"{key}.{k}" if key else k)
    elif isinstance(value, list):
      for v in value:
        yield from render(v, key)
    else:
      value = shlex.quote(os.path.expandvars(str(value)))
      yield f"{key}={value}" if key else value

  value = reduce(get, "{{profile}}.{{path}}".split("."), config)
  result = "\n".join(render(value))

  print(result)
