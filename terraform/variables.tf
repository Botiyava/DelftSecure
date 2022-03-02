variable "aws_region" {
  description = "AWS region to deploy infrastructure"
  type = string
  default = "eu-west-3"
}

variable "HTTP-HTTPS"{
  description = "HTTP-HTTPS"
  type = list
  default = ["80", "443"]
}

variable "cidr_block_VPC"{
  description = "CIDR block for VPC"
  type = string
  default = "172.31.0.0/16"
}

variable "cidr_block_App_subnet" {
  description = "public subnet for App"
  default = "172.31.10.0/24"
}
variable "cidr_block_DB_subnet" {
  description = "public subnet for App"
  default = "172.31.20.0/24"
}
variable "aws_availability_zone" {
  default = "eu-west-3a"
}
variable "instance_type" {
  default = "t2.micro"
}
variable "common_tags" {
  description = "Tags for all resources"
  type = map
  default = {
    Owner = "github.com/Botiyava"
    Project = "github.com/Botiyava/DelftSecure"
    Managed = "Terraform"
  }
}

