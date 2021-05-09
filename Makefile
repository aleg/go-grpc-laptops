install:
	go get google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc

gen:
	protoc --proto_path=proto proto/*.proto \
		--go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=require_unimplemented_servers=false:pb --go-grpc_opt=paths=source_relative

server:
	go run cmd/server/*.go --port 8080 -tls

client:
	go run cmd/client/*.go --address 0.0.0.0:8080 -tls

test:
	# go test github.com/aleg/go-grpc-laptops/serializer
	go test --cover --race ./...

cert:
	cd cert/; ./gen.sh; cd ..

clean-pb:
	rm pb/*.go

clean-tmp:
	rm -rf tmp/*.*

.PHONY: install gen server client test clean-pb clean-tmp cert
