output "primary_endpoint_address" {
  description = "Redis Primary Endpoint Address"
  value       = aws_elasticache_replication_group.main.primary_endpoint_address
}

output "primary_endpoint_port" {
  description = "Redis Primary Endpoint Port"
  value       = aws_elasticache_replication_group.main.port
}

output "security_group_id" {
  description = "Security Group ID of the Redis Cluster"
  value       = aws_security_group.redis.id
}
