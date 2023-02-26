variable "lambda_func_name" {
  type = string
}

variable "lambda_role_arn" {
  type = string
}

resource "aws_lambda_function" "versioning_lambda" {
  function_name = var.lambda_func_name
  filename = "./main.zip"
  role = var.lambda_role_arn
  handler = "main"
  runtime = "go1.x"

}
