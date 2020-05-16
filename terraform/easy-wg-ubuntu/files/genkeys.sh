#!/bin/sh

PRIVATE=$(wg genkey)
PUBLIC=$(echo $PRIVATE |wg pubkey)

cat <<EOF
{
 "private": "${PRIVATE}",
 "public":"${PUBLIC}"
}
EOF
