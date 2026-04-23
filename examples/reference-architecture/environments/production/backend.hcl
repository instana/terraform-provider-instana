# Production Environment - S3 Backend Configuration
# Usage: terraform init -backend-config=backend.hcl

bucket         = "your-terraform-state-bucket"
key            = "instana/production/terraform.tfstate"
region         = "us-east-1"
encrypt        = true
dynamodb_table = "terraform-state-lock"

# Optional: Additional backend configuration
# acl            = "private"
# kms_key_id     = "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"