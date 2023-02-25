variable "name" {
  description = "Name of the VM"
  type        = string
}

variable "disk_id" {
  description = "Id of the disk to attach to the VM"
  type        = string
}

variable "disk_mode" {
  description = "READ_ONLY or READ_WRITE permissions on disk"
  type        = string
  default = "READ_WRITE"
}

variable "disk_size_gib" {
  description = "Size in Gb of the disk to attach to the VM"
  type        = string
  default = "10Gib"
}

variable "os_family" {
  description = "Family of the OS to use for the VM"
  type        = string
}

variable "os_project" {
  description = "Projetc of the OS to use for the VM"
  type        = string
}

variable "region" {
  description = "Region where the VM should be located"
  type        = string
}

variable "subnet_id" {
  description = "Id of the subnet to attach to the VM"
  type        = string
}

variable "tags" {
  description = "Tags to set onto the VM"
  type        = map(string)
}

variable "vm_size" {
  description = "Type of the VM size to use"
  type        = string
}

variable "hostname" {
  description = "Hostname of the VM"
  type        = string
  default = "default-vm"
}

variable "user" {
  description = "The user to use for SSH access to the instance."
  default = ""
  sensitive = true
  type = string
}

variable "ssh_public_key" {
  description = "The public key to use for SSH access to the instance."
  default = ""
  sensitive = true
  type = string
}
