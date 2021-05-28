# xo - eXecution Orchestrator

`xo` is a small tool to handle munging data in our bundle and architecture workflows.

## Usage

**Validate a JSON Schema**:

Useful for:

* Validating _our_ bundles' schemas in _our_ CI before release
* Validating _user_ input to **packages** at the beginning of a workflow

```shell
xo schema validate --schema=cmd/testdata/valid-schema.json --document=cmd/testdata/valid-document.json

# or
cd massdriver-bundles
xo schema validate -s path/to/draft-07-schema.json -i ./bundles/$BUNDLE_NAME/schema-{inputs,connections,artifacts}.json
```

**Compiling variable definitions**:

[Terraform Variable Types](https://www.terraform.io/docs/configuration/expressions/types.html#types)

Output to STDOUT:

```shell
xo provisioner compile terraform -s examples/compiling-schemas/variables.schema.json -o -
```

Output to file:

```shell
xo provisioner compile terraform -s examples/compiling-schemas/variables.schema.json -o variables.tf.json
```

## Development

### Adding Commands

Add commands using the [Cobra Generator](https://github.com/spf13/cobra/blob/master/cobra/README.md).

Commands should be scoped (subcommand) under a parent "command" to facilitate organization.

Blogs on writing Cobra commands:

* https://towardsdatascience.com/how-to-create-a-cli-in-golang-with-cobra-d729641c7177

### Dereferencing $refs (hydration)

Hydrating schemas is important for a few reasons:

* storybook can't resolve json schema
* provides us w/ schema "snapshots" of a release
  * if clients are making http requests for schemas, artifacts, and specs, technically a deployment could change one of those files while they are requesting their "batch" to build the front-end UIs. By dereferencing everything into one file, it negates this changes
* less HTTP requests for the front-end clients
* much easier to debug,
* tfcompiler doesnt understand refs

Base on [compiletojsonschema](https://compiletojsonschema.readthedocs.io/en/latest/index.html)