output "bootstrap_brokers" {
  description = "Plaintext connection string"
  value       = aws_msk_cluster.main.bootstrap_brokers
}

output "bootstrap_brokers_tls" {
  description = "TLS connection string"
  value       = aws_msk_cluster.main.bootstrap_brokers_tls
}

output "zookeeper_connect_string" {
  description = "Zookeeper connection string"
  value       = aws_msk_cluster.main.zookeeper_connect_string
}

output "msk_security_group_id" {
  description = "Security Group ID of the MSK Cluster"
  value       = aws_security_group.msk.id
}
