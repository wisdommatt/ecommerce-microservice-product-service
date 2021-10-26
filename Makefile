protoc:
	protoc product.proto --go-grpc_out=. --go_out=.
	protoc user.proto --go-grpc_out=. --go_out=.
	
run:
	go run main.go

tests:
	go test ./... -race -cover

watch:
	go install github.com/cespare/reflex@latest
	reflex -s -- sh -c 'clear && PORT=2424 go run main.go'

gen-mocks:
	go get github.com/vektra/mockery/v2/.../
	mockery --all