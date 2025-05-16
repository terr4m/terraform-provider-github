set windows-shell := ["pwsh", "-NoLogo", "-Command"]

default:
  just --list

setup: tidy
  {{ if os() == "windows" { 'New-Item -Name "dist" -ItemType "directory" -Force' } else { "mkdir -p ./dist" } }}

tidy:
  go mod tidy

fmt:
  gofmt -s -w -e .

lint:
  golangci-lint run --timeout 120s

test:
  go test -v -cover -timeout=120s -parallel=10 ./...

testacc $TF_ACC="1":
  go test -v -cover -timeout 120m ./... -sweep=tf-acc-

build:
  go build -o ./dist -v ./...

[working-directory: "tools"]
docs:
  go generate ./...
