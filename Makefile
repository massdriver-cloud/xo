.PHONY: test
test:
	go test ./cmd
	go test ./src/...
	go build
	./xo schema validate --schema=src/jsonschema/testdata/valid-schema.json --document=src/jsonschema/testdata/valid-document.json
	./xo bundle build ./src/bundles/testdata/bundle.Build/bundle.yaml -o /tmp/test-bundle-build

.PHONY: docker.build
docker.build:
	DOCKER_BUILDKIT=1 docker build -t 005022811284.dkr.ecr.us-west-2.amazonaws.com/massdriver-cloud/xo .

hack.build-to-massdriver:
	GOOS=linux GOARCH=amd64 go build && cp ./xo ../massdriver/xo-amd64

local.build-to-m1:
	GOOS=darwin GOARCH=arm64 go build
