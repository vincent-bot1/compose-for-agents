group default {
  targets = [
    "all",
  ]
}

group images {
  targets = [
    "agents",
    "agents-ui",
    "gateway",
  ]
}

group darwin {
  targets = [
    "images",
    "docker-compose-darwin",
    "docker-mcpgateway-darwin",
  ]
}

group windows {
  targets = [
    "images",
    "docker-compose-windows",
    "docker-mcpgateway-windows",
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

target docker-mcpgateway {
  inherits = ["_base"]
  target = "docker-mcpgateway"
  output = ["./bin"]
  platforms = [
    "darwin/arm64",
  ]
}

target docker-mcpgateway-darwin {
  inherits = ["docker-mcpgateway"]
  platforms = ["darwin/arm64"]
}

target docker-mcpgateway-windows {
  inherits = ["docker-mcpgateway"]
  platforms = ["windows/amd64"]
}

target docker-compose {
  context = "https://github.com/fiam/compose.git"
  output = ["./bin"]
}

target docker-compose-darwin {
  inherits = ["docker-compose"]
  platforms = ["darwin/arm64"]
  target = "binary-darwin"
}

target docker-compose-windows {
  inherits = ["docker-compose"]
  platforms = ["windows/amd64"]
  target = "binary-windows"
}