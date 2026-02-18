output "vpc_id" {
  description = "VPC ID"
  value       = module.vpc.vpc_id
}

output "eks_cluster_endpoint" {
  description = "API endpoint for EKS cluster"
  value       = module.eks.cluster_endpoint
}

output "eks_cluster_name" {
  description = "EKS Cluster Name"
  value       = module.eks.cluster_name
}

output "rds_endpoint" {
  description = "RDS connection endpoint"
  value       = module.rds.rds_endpoint
}

output "kafka_brokers" {
  description = "MSK Bootstrap Brokers"
  value       = module.msk.bootstrap_brokers
}

output "redis_primary_endpoint" {
  description = "Redis primary endpoint"
  value       = module.elasticache.primary_endpoint_address
}
