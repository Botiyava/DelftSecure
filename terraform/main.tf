provider "aws"{
  region = var.region
}

//TODO add tf env in .gitignore

# Security group for App
resource "aws_security_group" "App_SG"{
  name = var.HTTP-HTTPS.description

  dynamic "ingress" {
    for_each = var.HTTP-HTTPS
    content {
      to_port = ingress.value
      from_port = ingress.value
      protocol = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    }
  }

  ingress{
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = []
  }

  egress{
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]

  }
  tags = {
    Name = "App_SG"
    Owner = "Leonid Shalaev"
    Managed = "Terraform"
  }
}

# Security group for Database
resource "aws_security_group" "DB_SG" {


  tags = {
    Name = "DB_SG"
    Owner = "Leonid Shalaev"
    Managed = "Terraform"
  }
}


resource "aws_db_instance" "DB " {
  instance_class = ""
}