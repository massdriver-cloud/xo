MASSDRIVER_PATH?=../massdriver
MASSDRIVER_PROTOS=${MASSDRIVER_PATH}/protos

.PHONY: test
test:
	go test ./cmd
	go test ./tfdef
	go build
	./xo schema validate --schema=cmd/testdata/valid-schema.json --document=cmd/testdata/valid-document.json
	./xo provisioner compile terraform -s examples/compiling-schemas/variables.schema.json -o -	

.PHONY: setup
setup: ## Install CLI/editor deps
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u github.com/twitchtv/twirp/protoc-gen-twirp
	go get google.golang.org/protobuf/reflect/protoreflect@v1.26.0
	go get google.golang.org/protobuf/runtime/protoimpl@v1.26.0

clean: 
	rm -rf massdriver/deployments.{pb,twirp}.go

massdriver/deployments.pb.go:
	protoc --proto_path=$(GOPATH)/src:$(MASSDRIVER_PROTOS):. --twirp_out=. --go_out=. $(MASSDRIVER_PROTOS)/deployments.proto
