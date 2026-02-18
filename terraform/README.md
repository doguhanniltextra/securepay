# SecurePay Infrastructure (Terraform)

This directory contains the Infrastructure-as-Code (IaC) definitions for deploying the SecurePay platform on AWS.

> **CRITICAL NOTE:** This is a demonstration implementation. It is **NOT** intended to be applied as-is due to AWS Free Tier limitations. The purpose of this code is to showcase production-ready AWS architecture skills, modularity, and best practices for a CV portfolio.

## Structure

The project follows a standard Terraform layout:

```text
terraform/
├── modules/               # Reusable infrastructure components
│   ├── vpc/              # Networking (VPC, Subnets, NAT, IGW)
│   ├── eks/              # Compute (EKS Cluster, Node Groups, IAM)
│   ├── rds/              # Database (PostgreSQL RDS)
│   ├── msk/              # Event Bus (Managed Kafka)
│   └── elasticache/      # Caching (Redis ElastiCache)
├── environments/          # Environment-specific configurations
│   └── dev/              # Development environment (composing modules)
└── README.md             # This file
```

## Purpose

- Demonstrate ability to design scalable, secure cloud architecture.
- Showcase modular Terraform development.
- Implement production-grade patterns:
  - Private subnets for workloads and data.
  - IAM Roles for Service Accounts (IRSA).
  - Managed services (EKS, RDS, MSK, ElastiCache) for operational excellence.
  - Encryption at rest and in transit.
  - Security Groups for least-privilege access.

## Usage

**Prerequisites:**
- Terraform v1.5+
- AWS CLI configured

**To plan (Dry Run):**

```bash
cd environments/dev
terraform init
terraform plan
```

**WARNING:** Running `terraform apply` will incur significant costs (NAT Gateways, EKS Control Plane, MSK Cluster, etc.). Do not execute unless you are prepared for the bill.
