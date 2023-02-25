variable "db_username" {
  description = "database username"
  type        = string
}

variable "db_password" {
  description = "database password"
  type        = string
}

variable "public_key" {
  description = "ssh public key for accessing the bastion instance"
  type        = string
}

variable "private_key_path" {
  description = "ssh private key path for accessing the bastion instance"
  type        = string
}
