# xo - eXecution Orchestrator

`xo` is a small tool to handle munging data in our bundle and architecture workflows.


## Usage

**Validate a JSON Schema**:

Useful for:

* Validating _our_ bundles' `manifest.json` in _our_ CI before release
* Validating _user_ input to **payloads** at the beginning of a workflow

```shell
xo schema validate --schema=cmd/testdata/valid-schema.json --document=cmd/testdata/valid-document.json

# or
cd massdriver-bundles
xo schema validate -s ./definitions/bundle-metadata.json -i ./bundles/$BUNDLE_NAME/metadata.json
```

**Generating variable definitions**:

[Terraform Variable Types](https://www.terraform.io/docs/configuration/expressions/types.html#types)

```shell
xo provisioner definitions terraform -s examples/tf-json-internals/variables.schema.json
```

## Development

### Adding Commands

Add commands using the [Cobra Generator](https://github.com/spf13/cobra/blob/master/cobra/README.md).

Commands should be scoped (subcommand) under a parent "command" to facilitate organization.

Blogs on writing Cobra commands:

* https://towardsdatascience.com/how-to-create-a-cli-in-golang-with-cobra-d729641c7177