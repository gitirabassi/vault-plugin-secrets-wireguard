# Vault-secret-plugin-wireguard

## Paths

- GET/POST/DELETE `/config`
  - `port`
  - `public_endpoint`
  - `client_persistent_keepalive`
  - `save_config`
  - `post_up_script`
  - `post_down_script`
  - `server_cidr` --> `10.20.0.1/24`
  - `webhook_address`
  - `webhook_secret`
- GET/POST/DELETE `/roles/:role_name`
  - `allowed_ips`
  - `dns`
  - `client_subnet_mask` --> 32
- GET `/creds/:role_name`
- GET `/server-creds`

## Apply policy
```shell
vault policy write wireguard_client_develper contrib/client.hcl
```

