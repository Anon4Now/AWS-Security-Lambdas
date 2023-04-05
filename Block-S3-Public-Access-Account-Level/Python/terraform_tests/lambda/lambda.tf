variable "lambda_func_name" {
  type = string
}

variable "lambda_role_arn" {
  type = string
}

resource "aws_lambda_function" "s3_blocking_lambda" {
  function_name = var.lambda_func_name
  filename = "./main.zip"
  role = var.lambda_role_arn
  handler = "s3_block_public_access_account_settings"
  runtime = "python3.9"

}
