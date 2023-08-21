gen:
	protoc --proto_path=proto proto/*.proto  --go_out=:pb --go-grpc_out=:pb

clean:
	rm pb/*

server:
	go run cmd/server/main.go -port 8080

client:
	go run cmd/client/main.go -address 127.0.0.1:8080

test:
	go test -cover -race ./...