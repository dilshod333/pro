RUN:
	@go run cmd/main.go 

tidy:
	go mod tidy 

generate-rpc:
	protoc \
	--go_out=genprotos \
	--go_opt=paths=source_relative \
	--go-grpc_out=genprotos \
	--go-grpc_opt=paths=source_relative \
	protos/countries.proto




