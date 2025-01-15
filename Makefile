run-consul:
	docker run -d -p 8500:8500 -p 8600:8600/udp --name=dev-consul hashicorp/consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0

gen-pb:
	protoc -I=proto --go_out=. --go-grpc_out=. movie.proto

run-all:
	go run metadata/cmd/main.go -domain "localhost" -port ":8080" -consul_address "localhost:8500" -kafka_address "localhost:9092" -db_driver "postgres" -db_host "postgres" -db_port "5432" -db_user "admin" -db_password "secret" -db_name "metadata_db" -db_sslmode "disable"; \
	go run rating/cmd/main.go -domain "localhost" -port ":8081" -consul_address "localhost:8500" -kafka_address "localhost:9092" -db_driver "postgres" -db_host "postgres" -db_port "5433" -db_user "admin" -db_password "secret" -db_name "ratings_db" -db_sslmode "disable"; \
	go run metadata/cmd/main.go -domain "localhost" -port ":8082" -consul_address "localhost:8500" -kafka_address "localhost:9092"