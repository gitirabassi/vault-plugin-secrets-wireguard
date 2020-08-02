
build:
	GOOS=linux go build .

docker:
	cd ../.. && docker build -t quay.io/gtirabassi/wg-config-gen:latest -f ./cmd/wg-config-gen/Dockerfile .

test:
	docker run --rm -it -v $$PWD:/opt/ quay.io/gtirabassi/wg-config-gen:latest -address=12.12.12.12 -output-dir=/opt/tmp -with-dns-gateway=true -num-clients=16

clean:
	rm tmp/*

push:
	docker push quay.io/gtirabassi/wg-config-gen:latest


	#  ssh core@3.122.221.129
	# docker run --rm -it -v $PWD:/opt/ quay.io/gtirabassi/wg-config-gen:latest -interface-name=wg1 -address=3.122.221.129 -output-dir=/opt/ -with-dns-gateway=true -num-clients=16
	# docker run --network=host -v /etc/wireguard:/etc/coredns coredns/coredns:1.6.9 -conf /etc/coredns/Corefile
	# docker run -network=host -v /etc/envoy:/etc/envoy envoyproxy/envoy:v1.14.1
	docker run -d -network=host -v /etc/envoy:/etc/envoy envoyproxy/envoy:v1.14.1^C



PostUp = 

sysctl -w net.ipv4.ip_forward=1
iptables -A FORWARD -p udp --dport 53 -i wg0 -j ACCEPT; iptables -A FORWARD -p tcp --dport 22 -i wg0 -j ACCEPT; iptables -A FORWARD -p tcp --dport 80 -i wg0 -j ACCEPT; iptables -A FORWARD -p tcp --dport 443 -i wg0 -j ACCEPT
iptables -A FORWARD -i wg0 -j ACCEPT

iptables -t nat -A OUTPUT -p tcp -s 10.29.100.10 -j  REDIRECT --to-port 15001
