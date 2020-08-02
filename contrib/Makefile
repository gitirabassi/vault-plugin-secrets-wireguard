


validate:
	docker run --rm -it -p 15000:15000 -v $$PWD:/etc/envoy envoyproxy/envoy:v1.14.1

exporter:
	docker build -t quay.io/gtirabassi/prometheus_wireguard_exporter .
	docker push quay.io/gtirabassi/prometheus_wireguard_exporter
