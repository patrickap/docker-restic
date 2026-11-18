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
        if k != "env" and isinstance(v, (str, int, float, bool))
      ]
    )

  _, args = reduce(expand, "{{profile}} {{args}}".split(), (config, []))
  result = os.path.expandvars(" ".join(args[1:]))
  
  print(result)

get profile path:
  #!/usr/bin/env python3

  import os, tomllib
  from functools import reduce

  with open(os.path.expandvars("{{config_file}}"), "rb") as f:
    config = tomllib.load(f)

  def get(node, key):
    return node.get(key, {})

  def render(value, key=""):
    if isinstance(value, dict):
      for k, v in value.items():
        render(v, f"{key}.{k}" if key else k)
    elif isinstance(value, list):
      for v in value:
        render(v, key)
    else:
      return f"{key}={value}" if key else value

  value = reduce(get, "{{profile}}.{{path}}".split("."), config)
  result = render(value)
  
  print(result)

env profile:
  #!/usr/bin/env python3

  import os, tomllib, shlex

  with open(os.path.expandvars("{{config_file}}"), "rb") as f:
    config = tomllib.load(f)

  def env(value):
    return " ".join(
      f"{k}={os.path.expandvars(str(v))}"
      for k, v in value.items()
    )

  value = config["{{profile}}"].get("env", {})
  result = env(value)

  print(result)
