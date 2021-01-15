.PHONY: test
test:
	go test ./cmd
	go test ./tfdef
	go build
	./xo schema validate --schema=cmd/testdata/valid-schema.json --document=cmd/testdata/valid-document.json
	./xo provisioner compile terraform -s examples/compiling-schemas/variables.schema.json -o -	