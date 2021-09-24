.PHONY: test
test:
	go test ./cmd
	go test ./src/...
	go build
	./xo schema validate --schema=src/jsonschema/testdata/valid-schema.json --document=src/jsonschema/testdata/valid-document.json
	./xo bundle build ./src/bundles/testdata/bundle.Build/bundle.yaml -o /tmp/test-bundle-build	

.PHONY: docker.build
docker.build:
	DOCKER_BUILDKIT=1 docker build --ssh default -t massdriver/xo .
