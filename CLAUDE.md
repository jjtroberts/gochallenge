# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Kubernetes deployment tool written in Go that demonstrates deploying nginx to AWS EKS. The main component is a CLI tool (`deploy.go`) that reads Kubernetes YAML files and creates or updates deployments in a Kubernetes cluster.

## Architecture

- **deploy.go**: Main CLI application that accepts `-yamlpath` and `-image` flags to deploy/update Kubernetes deployments
- **deployment.yaml**: Kubernetes deployment manifest for nginx
- **Makefile**: Build pipeline with formatting, linting, vetting, and building steps
- Uses Kubernetes client-go library for cluster interaction
- Supports both creating new deployments and updating existing ones

## Common Commands

### Prerequisites
Install required tools:
```bash
go install golang.org/x/lint@latest
go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest
```

### Development Commands
- **Format code**: `make fmt` or `go fmt ./...`
- **Lint code**: `make lint` (runs golint after formatting)
- **Vet code**: `make vet` (runs go vet and shadow after formatting)
- **Build**: `make build` (runs full pipeline: fmt → lint → vet → build)

### Running the Application
- **As script**: `go run deploy.go -yamlpath ./deployment.yaml -image nginx:1.17`
- **As binary**: 
  1. `make build`
  2. `./deploy -yamlpath ./deployment.yaml -image nginx:1.17`

## Key Implementation Details

- The deploy tool parses YAML files and handles apps/v1 Deployment resources
- Uses kubeconfig from `~/.kube/config` for cluster authentication
- Automatically detects if deployment exists and either creates or updates accordingly
- Updates the first container's image in the deployment spec
- Defaults to "default" namespace if none specified in YAML

## Dependencies

- Go 1.15+
- Kubernetes client-go v0.21.4
- kubectl and EKS cluster setup required for actual deployment