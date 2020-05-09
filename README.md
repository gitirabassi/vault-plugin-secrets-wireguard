# Vault-secret-plugin-wireguard

## Paths

- GET/POST/DELETE `/config`
  - `cidr` --> `10.20.0.0/24`
- GET/POST/DELETE `/servers/:server_name`
  - `port` --> defaults to `51820`
  - `public_endpoint` -->
  - `post_up_script` 
  - `post_down_script`
  - `private_webhook_address` --> defaults to `public_endpoint` in case it's not specified
  - `allowed_ips` --> list of subnets to route traffic trhu this server can be `["0.0.0.0/0"]` or `["10.0.0.0/24","192.68.0.0/24"]`
- GET/POST/DELETE `/roles/:role_name`
  - `servers` --> list of servers to connect to `["default", "aws-europe", "gcp-testing"]` --> must exist `:server_name`
  - `dns` 
  - `client_persistent_keepalive` --> defaults to 25 seconds
  - `client_subnet_mask` --> defaults to 32
- GET `/creds/:role_name`
  - `conf` --> complete wireguard configuration to be used with wg-quick for a client
- GET `/server-creds/:server_name`
  - `conf` --> complete wireguard configuration to be used with wg-quick for a client
  - `webhook_secret` --> webhook secret that vault will use to POST updates to wireguard servers

## Apply policy
```shell
vault policy write wireguard_client_develper contrib/client.hcl
```

