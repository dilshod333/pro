tidy:
	@go mod tidy
run:
	@go run cmd/main.go

generate-rpc:
	protoc \
	--go_out=genprotos \
	--go_opt=paths=source_relative \
	--go-grpc_out=genprotos \
	--go-grpc_opt=paths=source_relative \
	protos/countries.proto
