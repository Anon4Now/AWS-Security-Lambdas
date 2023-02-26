provider "aws" {
    region = "us-east-1"
}

resource "aws_kms_key" "key1" {
    description = "Used to test the key rotation Lambda."
    count = 2
}

module "lambdaModule" {
  source = "./lambda"
  lambda_func_name = "lambda_for_kms"
  lambda_role_arn = module.lambdaModule.lambda_role_arn
}