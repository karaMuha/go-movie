run-consul:
	docker run -d -p 8500:8500 -p 8600:8600/udp --name=dev-consul hashicorp/consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0

gen-pb:
	protoc -I=proto --go_out=. --go-grpc_out=. movie.proto

run:
	docker-compose up -d

prepare:
	mkdir data/metadata; \
	mkdir data/ratings; \
	protoc -I=proto --go_out=. --go-grpc_out=. movie.proto