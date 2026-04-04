group "default" {
  targets = [
    "curl",
    "server",
    "ui",
  ]
}

function "gentags" {
  params = [image]
  result = ["pneutrinoutil/${image}:local"]
}

target "curl" {
  context = "./docker/curl"
  dockerfile = "Dockerfile"
  tags = gentags("curl")
}

target "server" {
  context = "."
  dockerfile = "./server/Dockerfile"
  tags = gentags("server")
}

target "ui" {
  context = "./ui"
  dockerfile = "Dockerfile"
  tags = gentags("ui")
}
