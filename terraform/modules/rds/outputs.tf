output "rds_endpoint" {
  description = "RDS endpoint (host:port)"
  value       = aws_db_instance.main.endpoint
}

output "rds_address" {
  description = "RDS hostname"
  value       = aws_db_instance.main.address
}

output "rds_security_group_id" {
  description = "Security Group ID of the RDS Instance"
  value       = aws_security_group.rds.id
}
