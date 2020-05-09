variable backend_mount_path {
  type        = string
  description = "path at which to mount wireguard backend plugin"
  default     = "wireguard"
}

variable wireguard_cidr {
  type        = string
  description = "cidr to use in the wireguard mesh"

}

variable server {
  type = map(object({
    address  = string
    port     = string
    role_arn = string
    region   = string
    vpc_id   = string
    })
  )
}
