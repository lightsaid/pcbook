gen:
	protoc --proto_path=proto proto/*.proto  --go_out=:pb --go-grpc_out=:pb

clean:
	rm pb/*


test:
	go test -cover -race ./...