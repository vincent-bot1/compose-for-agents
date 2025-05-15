group default {
  targets = [
    "all",
  ]
}

group all {
  targets = [
    "agents",
    "agents-ui",
    "gateway",
  ]
}

# Required by docker/metadata-action and docker/bake-action gh actions.
target "docker-metadata-action" {}

target _base {
  inherits = ["docker-metadata-action"]
  output = ["type=docker"]
}

target agents {
  inherits = ["_base"]
  context = "agent"
  tags = ["demo/agents"]
}

target agents-ui {
  inherits = ["_base"]
  context = "agent-ui"
  tags = ["demo/ui"]
}

target gateway {
  inherits = ["_base"]
  target = "agents_gateway"
  tags = ["docker/agents_gateway"]
}
