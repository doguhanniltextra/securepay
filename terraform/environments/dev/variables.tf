variable "environment" {
  description = "Deployment environment (e.g., dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "project_name" {
  description = "Project name for resource naming"
  type        = string
  default     = "securepay"
}

variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

# DB Password should come from environment or secret manager (tfvars via -var-file=secrets.tfvars)
# For demo, variable is fine
variable "db_password" {
  description = "Password for RDS database"
  type        = string
  sensitive   = true
}
