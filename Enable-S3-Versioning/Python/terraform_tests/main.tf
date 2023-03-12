provider "aws" {
  region = "us-east-1"
}

module "s3Module" {
  source = "./s3"
  bucket_names = ["versionedbucketidu3iub3rb", "unversionedbucket32eb3b3ubf", "unversionedbucket3ibfin0dhf0ein"]
}

output "s3Output" {
  value = module.s3Module.S3Buckets
}

module "lambdaModule" {
  source = "./lambda"
  lambda_func_name = "lambda_for_s3"
  lambda_role_arn = module.lambdaModule.lambda_role_arn
}
