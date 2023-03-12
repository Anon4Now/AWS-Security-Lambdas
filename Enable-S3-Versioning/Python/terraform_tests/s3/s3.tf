provider "aws" {
    region = "us-east-1"
}

variable "bucket_names" {
  type = list(string)
}

resource "aws_s3_bucket" "myBucket" {
    count = length(var.bucket_names)
    bucket = var.bucket_names[count.index]
}

resource "aws_s3_bucket_versioning" "myBucketVersioning" {
    bucket = aws_s3_bucket.myBucket[0].id
    versioning_configuration {
      status = "Enabled"
    }
}

output "S3Buckets" {
  value = [aws_s3_bucket.myBucket.*.id]
}