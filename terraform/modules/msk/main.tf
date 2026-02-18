# Security Group for MSK
resource "aws_security_group" "msk" {
  name        = "${var.cluster_name}-sg"
  description = "Security group for SecurePay MSK"
  vpc_id      = var.vpc_id

  # Plaintext
  ingress {
    from_port   = 9092
    to_port     = 9092
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
    Name        = "${var.cluster_name}-sg"
    Environment = var.environment
  }
}

# MSK Cluster (No auth for demo/cost simplicity)
resource "aws_msk_cluster" "main" {
  cluster_name           = var.cluster_name
  kafka_version          = var.kafka_version
  number_of_broker_nodes = var.number_of_broker_nodes

  broker_node_group_info {
    instance_type   = var.broker_node_type
    client_subnets  = var.subnet_ids
    security_groups = [aws_security_group.msk.id]
    
    storage_info {
      ebs_storage_info {
        volume_size = 10 # GB
      }
    }
  }

  client_authentication {
    unauthenticated = true # For demo purposes
  }

  encryption_info {
    encryption_in_transit {
      client_broker = "PLAINTEXT"
    }
  }

  tags = {
    Environment = var.environment
  }
}
