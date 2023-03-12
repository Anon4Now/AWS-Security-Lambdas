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

resource "aws_iam_policy" "lambda_policy_for_s3_versioning" {
  name = "s3_versioning_policy"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "s3:GetBucketVersioning",
                "logs:PutLogEvents",
                "s3:PutBucketVersioning"
            ],
            "Resource": [
                "arn:aws:s3:::*",
                "arn:aws:logs:*:108554963036:log-group:*:log-stream:*"
            ]
        },
        {
            "Sid": "VisualEditor1",
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogStream",
                "logs:CreateLogGroup"
            ],
            "Resource": "arn:aws:logs:*:108554963036:log-group:*"
        },
        {
            "Sid": "VisualEditor2",
            "Effect": "Allow",
            "Action": "s3:ListAllMyBuckets",
            "Resource": "*"
        }
    ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_access" {
  role = aws_iam_role.iam_role_for_lambda.name
  policy_arn = aws_iam_policy.lambda_policy_for_s3_versioning.arn
}

output "lambda_role_arn" {
  value = aws_iam_role.iam_role_for_lambda.arn
}