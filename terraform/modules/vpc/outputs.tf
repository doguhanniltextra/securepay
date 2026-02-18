output "vpc_id" {
  description = "The ID of the VPC"
  value       = aws_vpc.main.id
}

output "public_subnet_ids" {
  description = "List of IDs of public subnets"
  value       = aws_subnet.public[*].id
}

output "private_subnet_ids" {
  description = "List of IDs of private subnets"
  value       = aws_subnet.private[*].id
}

output "subnet_cidr_blocks" {
  description = "List of CIDR blocks of the subnets"
  value       = concat(aws_subnet.public[*].cidr_block, aws_subnet.private[*].cidr_block)
}
