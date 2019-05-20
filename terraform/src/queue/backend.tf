terraform {
    backend "s3" {
    encrypt = true
    bucket = "mm-state-bucket"
    region = "us-east-1"
    key = "states/statefile"
    }
}
