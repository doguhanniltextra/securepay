variable "cluster_name" {
  description = "Name of the MSK cluster"
  type        = string
  default     = "securepay-msk"
}

variable "subnet_ids" {
  description = "List of private subnet IDs for MSK (requires at least 2)"
  type        = list(string)
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "kafka_version" {
  description = "Apache Kafka version"
  type        = string
  default     = "3.6.0"
}

variable "broker_node_type" {
  description = "MSK broker instance type"
  type        = string
  default     = "kafka.t3.small"
}

variable "number_of_broker_nodes" {
  description = "Total number of brokers"
  type        = number
  default     = 2 # Minimum for HA
}

variable "environment" {
  description = "Environment"
  type        = string
  default     = "dev"
}
