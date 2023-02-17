# AWS Lambda Code for Security Functions

The code in the individual directories are geared towards specific needs from a security perspective in AWS. It has been written in both Python and Go, depending on the use cases and comfortability of implementation.

AWS has numerous controls that it needs to follow to be compliant with any number of compliance frameworks (NIST, CIS, ISO, etc..), these Lambda's will help enforce that compliance by automating remediation actions that tend to frequently crop up.

The Lambda's created thus far are as follows:

- Enable yearly rotation on Customer Managed Keys (CMKs)
- Enable Object Versioning for S3 Buckets
- Block Public Access from S3 Buckets (at the account level)

## Lambda Functionality:

- Will use Boto3 or Go SDK to programmtically interact with AWS
- Can be run from inside a Docker container or uploaded via zip file
- Will require appropriate permissions for Lambda execution role to perform tasks successfully
- ###IMPORTANT These Lambda's are geared towards a single account deployment strategy, if a centralized approach is needed the code will need some changes

## Quick Notes:

- This code can be altered to be used in a multi-account environment, or be used as part of a pipeline deployment
- This code contains tests associated with its base functionality, and these tests DO NOT interact with AWS because of mocks
- There will be Terraform configuration files that will create basic infrastucture to prove out each Lambda's functionality (don't forget to delete after tests)


