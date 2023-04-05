# AWS Lambda Code for Rotating CMK's

The Python code in this directory is designed to be used with AWS Lambda. It will look to see is the S3 account level public access block is set and enable it if it is not. The practice of blocking S3 public access, follows best-practice for a number of compliance frameworks. Below is a basic listing of some frameworks that this control implements:

- CISAWSF
- APRA
- MAS
- NIST4

## Lambda Functionality:

- Will use Boto3 to programmtically interact with AWS
- Can be run from inside a Docker container or uploaded via zip file
- Will require appropriate permissions for Lambda execution role to perform S3Control tasks successfully
- **IMPORTANT** These Lambda's are geared towards a single account/per region deployment strategy, if a centralized approach is needed the code will need some changes

## Quick Notes:

- This code can be altered to be used in a multi-account environment, or be used as part of a pipeline deployment
- This code contains unit tests associated with its base functionality, and these tests DO NOT interact with AWS because of mocks
- There are Terraform configuration files that will create basic infrastucture to prove out each Lambda's functionality (don't forget to delete after tests)
	- NEED TO ADD A LAMBDA ENV VAR THAT TELLS LAMBDA WHAT LOG LEVEL TO USE , ELSE IT WILL ERROR(e.g., log_level = INFO)
- To use the Terraform configuration files as is, the Python code will need to be converted to a zip file and dropped in the terraform_tests directory

### Windows Example (PowerShell):
```
$compress = @{
Path = ".\file1.py", ".\file2.py"
DestinationPath = ".\app.zip"
}

Compress-Archive @compress
```

