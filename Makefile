.PHONY: test
test:
	go test ./cmd
	go test ./tfdef
	go build
	./xo provisioner definitions terraform -s examples/tf-json-internals/variables.schema.json
	./xo schema validate --schema=cmd/testdata/valid-schema.json --document=cmd/testdata/valid-document.json