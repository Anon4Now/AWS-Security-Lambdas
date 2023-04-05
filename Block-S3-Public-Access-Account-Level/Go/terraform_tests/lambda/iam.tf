provider "aws" {
  region = "us-east-1"
}

resource "aws_iam_role" "iam_role_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "lambda.amazonaws.com"
            },
            "Action": "sts:AssumeRole"
        }
    ]
}
EOF
}


resource "aws_iam_policy" "lambda_policy_for_s3_account_block" {
    name = "s3_account_access_block"

    policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": "logs:PutLogEvents",
            "Resource": "arn:aws:logs:*:623714057814:log-group:*:log-stream:*"
        },
        {
            "Sid": "VisualEditor1",
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogStream",
                "logs:CreateLogGroup"
            ],
            "Resource": "arn:aws:logs:*:623714057814:log-group:*"
        },
        {
            "Sid": "VisualEditor2",
            "Effect": "Allow",
            "Action": [
                "s3:PutAccountPublicAccessBlock",
                "s3:GetAccountPublicAccessBlock"
            ],
            "Resource": "*"
        }
    ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_access" {
  role = aws_iam_role.iam_role_for_lambda.name
  policy_arn = aws_iam_policy.lambda_policy_for_s3_account_block.arn
}

output "lambda_role_arn" {
  value = aws_iam_role.iam_role_for_lambda.arn
}