variable "region" {
  description = "AWS region to deploy infrastructure"
  default = "eu-west-3"
}

variable "HTTP-HTTPS"{
  description = "HTTP-HTTPS"
  type = list
  default = ["80", "443"]
}

variable "cidr_block_VPC"{
  description = "CIDR block for VPC"
  type = list
  default = ["172.31.0.0/16"]
}