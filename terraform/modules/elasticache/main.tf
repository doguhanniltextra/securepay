# Security Group for Redis
resource "aws_security_group" "redis" {
  name        = "${var.identifier}-sg"
  description = "Security group for SecurePay Redis"
  vpc_id      = var.vpc_id

  ingress {
    from_port   = var.port
    to_port     = var.port
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "${var.identifier}-sg"
    Environment = var.environment
  }
}

# Redis Subnet Group
resource "aws_elasticache_subnet_group" "main" {
  name       = "${var.identifier}-subnet-group"
  subnet_ids = var.subnet_ids
}

# Redis Replication Group (Cluster Mode Disabled)
resource "aws_elasticache_replication_group" "main" {
  replication_group_id = var.identifier
  description          = "SecurePay Redis Replication Group"
  node_type            = var.node_type
  port                 = var.port
  num_cache_clusters   = var.num_cache_nodes

  engine               = "redis"
  engine_version       = var.engine_version
  parameter_group_name = "default.redis7"

  subnet_group_name  = aws_elasticache_subnet_group.main.name
  security_group_ids = [aws_security_group.redis.id]

  at_rest_encryption_enabled = true
  transit_encryption_enabled = false # Simpler for demo (no TLS)

  tags = {
    Name        = var.identifier
    Environment = var.environment
  }
}
