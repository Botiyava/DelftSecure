provider "aws" {
  region = var.aws_region
}
locals {
  Project = "DelftSecure"
}
data "aws_ami" "latest_ubuntu" {
  owners      = ["099720109477"]
  most_recent = true
  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }
}

/*data "aws_ami" "latest_amazon2" {
  owners      = ["137112412989"]
  most_recent = true
  filter {
    name   = "name"
    values = ["amzn2-ami-kernel-5.10-hvm-2.0.20220218.0-arm64-gp2"]
  }
}*/

# VPC for all resources
resource "aws_vpc" "main_vpc" {
  cidr_block = var.cidr_block_VPC
  tags       = merge(var.common_tags, { Name = "${local.Project}_vpc" })
}

# Internet gateway: allow the VPC to connect to the internet
resource "aws_internet_gateway" "main_igw" {
  vpc_id = aws_vpc.main_vpc.id
  tags   = merge(var.common_tags, { Name = "${local.Project}_igw" })
}


# ----------BLOCK FOR PUBLIC RESOURCES----------------
# Route table for public subnet
resource "aws_route_table" "App_RT" {
  vpc_id = aws_vpc.main_vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main_igw.id
  }
  tags = merge(var.common_tags, { Name = "${local.Project}_App_RT" })
}

resource "aws_subnet" "App_subnet" {
  vpc_id                  = aws_vpc.main_vpc.id
  cidr_block              = var.cidr_block_App_subnet
  availability_zone       = var.aws_availability_zone
  map_public_ip_on_launch = true
  tags                    = merge(var.common_tags, { Name = "${local.Project}_app_subnet" })
}
resource "aws_route_table_association" "app_rta" {
  route_table_id = aws_route_table.App_RT.id
  subnet_id      = aws_subnet.App_subnet.id

}

# Public key to use to login to the EC2 instance
resource "aws_key_pair" "ssh_key" {
  key_name   = "tf_delft"
  public_key = file("delft.pub")
}
resource "aws_instance" "App_instance" {
  ami                    = data.aws_ami.latest_ubuntu.id
  subnet_id              = aws_subnet.App_subnet.id
  instance_type          = var.instance_type
  private_ip             = "172.31.10.100"
  vpc_security_group_ids = [aws_security_group.App_SG.id]

  /* user_data = file("external_files/app.sh")*/
  availability_zone = var.aws_availability_zone
  user_data = templatefile("external_files/app.sh", {
    config1 = file("external_files/config1.json")
  })

  key_name = aws_key_pair.ssh_key.key_name
  tags = merge(var.common_tags, {
    Name = "${local.Project}_app_instance"
  })
  connection {
    type     = "ssh"
    user     = "ubuntu"
    private_key = file("delft")
    host = self.public_ip
  }
  provisioner "remote-exec" {
    inline = [
      "sleep 10",
      "cd /home/ubuntu/DelftSecure",
      "nohup go run cmd/main.go &",
      "sleep 2",
    ]
  }
}

resource "aws_eip" "eip_for_nat" {

}
resource "aws_nat_gateway" "nat_for_db" {
  allocation_id = aws_eip.eip_for_nat.id
  subnet_id     = aws_subnet.App_subnet.id

  tags = merge(var.common_tags, {
    Name = "${local.Project}_nat_for_db"
  })

  # To ensure proper ordering, it is recommended to add an explicit dependency
  # on the Internet Gateway for the VPC.
  depends_on = [aws_internet_gateway.main_igw]
}
# ----------BLOCK FOR PRIVATE RESOURCES----------------
resource "aws_route_table" "db_rt" {
  vpc_id = aws_vpc.main_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_nat_gateway.nat_for_db.id
  }
  tags = merge(var.common_tags, {
    Name = "${local.Project}_DB_RT"
  })
}
resource "aws_subnet" "DB_subnet" {
  vpc_id            = aws_vpc.main_vpc.id
  cidr_block        = var.cidr_block_DB_subnet
  availability_zone = var.aws_availability_zone
  tags = merge(var.common_tags, {
    Name = "${local.Project}_DB_subnet"
  })

}
resource "aws_route_table_association" "db_rta" {
  route_table_id = aws_route_table.db_rt.id
  subnet_id      = aws_subnet.DB_subnet.id

}
resource "aws_instance" "DB_instance" {
  ami           = "ami-0b772a840c4c48ac2" # Amazon Linux 2
  private_ip    = "172.31.20.100"
  instance_type = "t4g.micro"

  key_name = aws_key_pair.ssh_key.key_name
  user_data = templatefile("external_files/install_postgres.sh", {
    pg_hba_file = templatefile("external_files/pg_hba.conf", { allowed_ip = "0.0.0.0/0" }),
  })
  subnet_id              = aws_subnet.DB_subnet.id
  vpc_security_group_ids = [aws_security_group.db_SG.id]
  availability_zone      = var.aws_availability_zone

  tags = merge(var.common_tags, {
    Name = "${local.Project}_DB_SG"
  })
}

/*resource "aws_route_table" "DB_RT" {
  vpc_id = aws_vpc.main_vpc.id
  route {
    cidr_block = "172.31.0.0/16"

  }
  tags = merge(var.common_tags, {Name = "${local.Project}_DB_RT"})
}*/
# --------------Security Groups----------------
# Security group for App
resource "aws_security_group" "App_SG" {

  vpc_id = aws_vpc.main_vpc.id
  dynamic "ingress" {
    for_each = var.HTTP-HTTPS
    content {
      to_port     = ingress.value
      from_port   = ingress.value
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    }
  }

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]

  }
  tags = merge(var.common_tags, { Name = "${local.Project}_App_SG" })
}

# SG for bastion

# SG for DB
resource "aws_security_group" "db_SG" {
  name        = "PostgreSQL"
  description = "Allow SSH and PostgreSQL inbound traffic"
  vpc_id      = aws_vpc.main_vpc.id

  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    description = "psql"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [var.cidr_block_VPC]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.common_tags, { Name = "${local.Project}_DB_SG" })
}

