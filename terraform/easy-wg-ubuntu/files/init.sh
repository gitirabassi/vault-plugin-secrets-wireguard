#!/bin/bash

set -x
set -e

apt update
apt install wireguard curl tar resolvconf iptables -y

mkdir -p /etc/wireguard
cat <<EOF >/etc/wireguard/wg0.conf
[Interface]
PrivateKey = ${server_private_key}
Address = ${wireguard_gateway_address}/${wireguard_gateway_netmask}
ListenPort = ${server_port}
PostUp = sysctl -w net.ipv4.ip_forward=1; iptables -A FORWARD -i %i -j ACCEPT; iptables -t nat -A POSTROUTING -o ${main_interface_name} -j MASQUERADE
PostDown = iptables -D FORWARD -i %i -j ACCEPT; iptables -t nat -D POSTROUTING -o ${main_interface_name} -j MASQUERADE

${peers}

EOF

systemctl enable wg-quick@wg0.service
systemctl start wg-quick@wg0.service

cat <<EOF >/opt/clients.conf
${client_configs_joined}
EOF

cat <<EOF >/opt/Corefile
. {
    bind ${wireguard_gateway_address}
    forward . 1.1.1.1 8.8.8.8
    log
    errors
}
EOF

curl -L -O https://github.com/coredns/coredns/releases/download/v1.6.9/coredns_1.6.9_linux_amd64.tgz
tar zxvf coredns_1.6.9_linux_amd64.tgz
mv coredns /usr/bin/coredns
chmod +x /usr/bin/coredns

cat <<EOF >/etc/systemd/system/coredns.service
[Unit]
Description=coredns
Wants=network-online.target
Requires=network-online.target
After=network-online.target

[Service]
ExecStart=/usr/bin/coredns -conf=/opt/Corefile
Restart=always
StartLimitInterval=5
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable coredns.service
systemctl start coredns.service
