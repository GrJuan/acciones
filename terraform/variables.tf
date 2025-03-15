variable "aws_region" {
  description = "La región de AWS donde se desplegará la infraestructura"
  type        = string
  default     = "us-west-2"
}

variable "app_name" {
  description = "Nombre de la aplicación"
  type        = string
  default     = "swe_app"
}

variable "instance_type" {
  description = "Tipo de instancia EC2"
  type        = string
  default     = "t3.micro"
}

variable "public_key_path" {
  description = "Ruta de la clave pública SSH para acceder a la instancia EC2"
  type        = string
  default     = "~/.ssh/id_rsa.pub"
}

variable "ami_id" {
  description = "AMI ID de la imagen de la instancia EC2"
  type        = string
  default     = "ami-0c55b159cbfafe1f0" # Amazon Linux 2 AMI (us-west-2)
}
