# SecurePay - Dev Environment

This directory contains the Terraform configuration for the Development environment.
It orchestrates the entire infrastructure stack by composing reusable modules.

> **WARNING:** This infrastructure is designed for production-grade demonstration. Running `terraform apply` will provision multiple managed services (EKS, RDS, MSK, ElastiCache, NAT Gateway) which will incur significant AWS costs immediately. Do not execute unless you intend to pay for these resources.

## Architecture

- **VPC**: 10.0.0.0/16, 3 AZs
  - Public Subnets (NAT Gateway, ALB)
  - Private Subnets (EKS Nodes, RDS, MSK, Redis)
- **EKS Cluster**: "securepay-dev-eks", Kubernetes 1.29
  - Managed Node Group in Private Subnets
  - OIDC Provider for IRSA
- **RDS**: "securepay-dev-db", PostgreSQL 16
  - db.t3.micro (Free-tier eligible instance type chosen for demo)
- **MSK**: "securepay-dev-msk", Kafka 3.6.0
  - 2 Brokers, kafka.t3.small
- **ElastiCache**: "securepay-dev-redis", Redis 7.0
  - Single node, cache.t3.micro

## Usage

1. **Initialize Terraform:**
   ```bash
   terraform init
   ```

2. **Configure Variables:**
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   # Edit terraform.tfvars to set db_password and region
   ```

3. **Plan (Dry Run):**
   ```bash
   terraform plan
   ```

4. **Apply (Cost Warning!):**
   ```bash
   terraform apply
   ```

## Cost Estimation (Monthly, Approximate)

| Resource | Type | Cost (est.) |
|----------|------|-------------|
| **Control Plane** | EKS Cluster | ~$73 |
| **Compute** | 2x t3.medium (Nodes) | ~$60 |
| **Database** | db.t3.micro (RDS) | ~$12 (or Free Tier) |
| **Streaming** | 2x kafka.t3.small (MSK) | ~$50 |
| **Cache** | cache.t3.micro (Redis) | ~$12 (or Free Tier) |
| **Network** | NAT Gateway | ~$32 + Data Transfer |
| **Total** | | **~$240+ / month** |

*Note: Prices vary by region.*
