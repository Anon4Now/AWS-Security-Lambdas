# AWS Lambda Code for Blocking S3 Public Access

The Go code in this directory is designed to be used with AWS Lambda. It will look at the current account's S3 public block access settings and enable if not set. The practice of blocking S3 public access, follows best-practice for a number of compliance frameworks. Below is a basic listing of some frameworks that this control implements:

- CISAWSF
- APRA
- MAS
- NIST4

## Lambda Functionality:

- Will use Go SDK to programmtically interact with AWS
- Can be run from inside a Docker container or uploaded via zip file
- Will require appropriate permissions for Lambda execution role to perform S3Control tasks successfully
- **IMPORTANT** These Lambda's are geared towards a single account/per region deployment strategy, if a centralized approach is needed the code will need some changes

## Quick Notes:

- This code can be altered to be used in a multi-account environment, or be used as part of a pipeline deployment
- This code contains unit tests associated with its base functionality, and these tests DO NOT interact with AWS because of mocks
- There are Terraform configuration files that will create basic infrastucture to prove out each Lambda's functionality (don't forget to delete after tests)
- To use the Terraform configuration files as is, the Go code will need to be converted to a zip file and dropped in the terraform_tests directory

### Linux Example (compile binary and convert to zip):
```
env GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./dist/main ./cmd/main.go
cd ./dist/ && zip main.zip main && cd ..
```

