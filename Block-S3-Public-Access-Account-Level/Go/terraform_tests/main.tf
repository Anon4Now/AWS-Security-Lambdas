provider "aws" {
  region = "us-east-1"
}

module "iam_role" {
  source = "./lambda"
  lambda_func_name = "lambda_for_s3_blocking"
  lambda_role_arn = module.iam_role.lambda_role_arn
}