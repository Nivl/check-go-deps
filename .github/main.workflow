workflow "Lint" {
  resolves = [
    "golangci-lint",
  ]
  on = "push"
}

action "golangci-lint" {
  uses = "cedrickring/golang-action@1.2.0"
  args = "./tools/lint.sh"
  env = {
    GO111MODULE = "on"
    GOFLAGS = "-mod=readonly"
  }
}
