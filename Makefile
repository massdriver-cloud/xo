MASSDRIVER_PATH?=../massdriver
MASSDRIVER_PROTOS=${MASSDRIVER_PATH}/protos

.PHONY: test
test:
	go test ./cmd
	go test ./src/bundles
	go test ./src/massdriver
	go test ./src/schemaloader
	go test ./src/tfdef
	go test ./src/generator
	go build
	./xo schema validate --schema=cmd/testdata/valid-schema.json --document=cmd/testdata/valid-document.json
	./xo provisioner terraform compile -s examples/compiling-schemas/variables.schema.json -o -
	./xo bundle build ./src/bundles/testdata/bundle.Build/bundle.yaml -o /tmp/test-bundle-build	

.PHONY: setup
setup: ## Install CLI/editor deps
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u github.com/twitchtv/twirp/protoc-gen-twirp
	go get google.golang.org/protobuf/reflect/protoreflect@v1.26.0
	go get google.golang.org/protobuf/runtime/protoimpl@v1.26.0

.PHONY: docker.build
docker.build:
	docker build -t massdriver/xo .

clean:
	rm -rf massdriver/*.{pb,twirp}.go

massdriver/workflow.pb.go:
	protoc --proto_path=$(GOPATH)/src:$(MASSDRIVER_PROTOS):. --twirp_out=. --go_out=. $(MASSDRIVER_PROTOS)/workflow.proto
