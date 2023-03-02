terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.55.0"
    }
  }
  required_version = ">= 1.3.9"

  // prepare s3 bucket and dynamo-db in advacne.
  backend "s3" {
    bucket         = "bucket-yutaroyamanaka"
    key            = "terraform"
    region         = "ap-northeast-1"
    dynamodb_table = "terraform-state-lock"
  }
}

provider "aws" {
  region = "ap-northeast-1"
}

/* DO NOT APPLY FOLLOWING RESOURCES EVERY TIME
resource "aws_s3_bucket" "terraform-backend" {
  bucket = "bucket-yutaroyamanaka"
  acl    = "private"
}

resource "aws_dynamodb_table" "eks-lock-table" {
  name           = "terraform-state-lock"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "LockID"
  attribute {
    name = "LockID"
    type = "S"
  }
}
*/
