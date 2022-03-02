output "App_server_IP_address" {
  value = aws_instance.App_instance.public_ip
}