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

resource "aws_iam_policy" "lambda_policy_for_kms_rotation" {
    name = "kms-rotation-policy"

    policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": "kms:ListKeys",
            "Resource": "*"
        },
        {
            "Sid": "VisualEditor1",
            "Effect": "Allow",
            "Action": [
                "kms:EnableKeyRotation",
                "kms:GetKeyRotationStatus",
                "kms:DescribeKey",
                "logs:PutLogEvents"
            ],
            "Resource": [
                "arn:aws:kms:*:428436893676:key/*",
                "arn:aws:logs:*:428436893676:log-group:*:log-stream:*"
            ]
        },
        {
            "Sid": "VisualEditor2",
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogStream",
                "logs:CreateLogGroup"
            ],
            "Resource": "arn:aws:logs:*:428436893676:log-group:*"
        }
    ]
}

EOF
}

resource "aws_iam_role_policy_attachment" "lambda_access" {
  role = aws_iam_role.iam_role_for_lambda.name
  policy_arn = aws_iam_policy.lambda_policy_for_kms_rotation.arn
}

output "lambda_role_arn" {
  value = aws_iam_role.iam_role_for_lambda.arn
}