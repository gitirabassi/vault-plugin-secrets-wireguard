
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
