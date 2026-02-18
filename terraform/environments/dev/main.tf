module "vpc" {
  source       = "../../modules/vpc"
  project_name = var.project_name
  environment  = var.environment
  vpc_cidr     = var.vpc_cidr
  # Default subnets are fine
}

module "eks" {
  source           = "../../modules/eks"
  cluster_name     = "${var.project_name}-${var.environment}-eks"
  environment      = var.environment
  vpc_id           = module.vpc.vpc_id
  subnet_ids       = module.vpc.private_subnet_ids # Control plane & nodes in private subnets
  
  # Node Group Config
  node_group_name  = "general-workers"
  desired_size     = 2
  min_size         = 1
  max_size         = 3
  instance_types   = ["t3.medium"]
}

module "rds" {
  source            = "../../modules/rds"
  identifier        = "${var.project_name}-${var.environment}-db"
  environment       = var.environment
  vpc_id            = module.vpc.vpc_id
  subnet_ids        = module.vpc.private_subnet_ids
  
  db_name           = "securepay"
  db_username       = "securepay_admin"
  db_password       = var.db_password # In real world, use Secrets Manager
  instance_class    = "db.t3.micro"
  allocated_storage = 20
}

module "msk" {
  source                 = "../../modules/msk"
  cluster_name           = "${var.project_name}-${var.environment}-msk"
  environment            = var.environment
  vpc_id                 = module.vpc.vpc_id
  subnet_ids             = module.vpc.private_subnet_ids
  number_of_broker_nodes = 2
  broker_node_type       = "kafka.t3.small"
}

module "elasticache" {
  source          = "../../modules/elasticache"
  identifier      = "${var.project_name}-${var.environment}-redis"
  environment     = var.environment
  vpc_id          = module.vpc.vpc_id
  subnet_ids      = module.vpc.private_subnet_ids
  node_type       = "cache.t3.micro"
  num_cache_nodes = 1
}
