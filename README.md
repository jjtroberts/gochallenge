# Rackner's K8s AWS Deploy Challenge

The purpose of this challenge is to test your ability to work with Kubernetes, Go, and AWS.

It's open ended in how you complete the challenge to achieve the final goal: creating a deploy script run a `nginx` deployment in EKS.

## Deliverables

You will create Go and Kubernetes YAML files to complete this project. Please save and include any of these files in the final submission. Please also screen shot the AWS console showing the EKS cluster.

## Challenge Steps

1. Create an EKS cluster in AWS (using any method you prefer: the console, Terraform, aws-cli, etc)
`eksctl create cluster --name gochallenge --region us-east-1`
2. Write Kubernetes yaml (using Helm, or your preferred templating method) to create a deployment for an nginx image
`kubectl create deployment nginx --image=nginx:1.17 -o yaml --dry-run=client > deployment.yaml`
3. Write a `deploy` go cli script which creates or updates the nginx deployment in k8s with a specific image tag. This image tag should be passed into the go script via a cli flag.
4. Document how to run the script.

## Prereqs
1. golint `go install golang.org/x/lint@latest`
2. shadow `go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest`

## Usage
As a script: `go run deploy.go -yamlpath ./deployment.yaml -image nginx:1.17`

As a command:
1. `make build`
2. `./deploy -yamlpath ./deployment.yaml -image nginx:1.17`