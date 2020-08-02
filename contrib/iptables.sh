#!/bin/sh

# Create new chain in nat table WIREGUARD_INBOUND
iptables -t nat -N WIREGUARD_INBOUND

# route all TCP traffic from wg0 to wireguard
iptables -t nat -A PREROUTING -p tcp -j WIREGUARD_INBOUND

iptables -t nat -A ISTIO_REDIRECT -p tcp -j REDIRECT --to-port 15001
