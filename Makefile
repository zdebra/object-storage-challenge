run:
	go run main.go

generate_server:
	go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest -package=genserver -generate types,chi-server,spec -o gen-server/server.gen.go openapi.yaml
	go mod tidy
