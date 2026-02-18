output "cluster_id" {
  description = "EKS Cluster ID"
  value       = aws_eks_cluster.main.id
}

output "cluster_endpoint" {
  description = "EKS Cluster Endpoint"
  value       = aws_eks_cluster.main.endpoint
}

output "cluster_certificate_authority_data" {
  description = "EKS Cluster Certificate Authority Data"
  value       = aws_eks_cluster.main.certificate_authority[0].data
}

output "cluster_security_group_id" {
  description = "Security Group ID attached to the EKS cluster"
  value       = aws_eks_cluster.main.vpc_config[0].cluster_security_group_id
}

output "cluster_iam_role_name" {
  description = "IAM Role Name of the EKS cluster"
  value       = aws_iam_role.cluster_role.name
}

output "node_group_role_arn" {
  description = "IAM Role ARN of the Worker Nodes"
  value       = aws_iam_role.node_role.arn
}

output "oidc_provider_url" {
  description = "OIDC Provider URL for IRSA"
  value       = aws_iam_openid_connect_provider.eks.url
}

output "oidc_provider_arn" {
  description = "OIDC Provider ARN for IRSA"
  value       = aws_iam_openid_connect_provider.eks.arn
}
