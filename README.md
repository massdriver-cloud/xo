# xo - eXecution Orchestrator

`xo` is a small tool to handle munging data in our bundle and architecture workflows.


## Usage

**Validate a JSON Schema**:

Useful for:

* Validating _our_ bundles' `manifest.json` in _our_ CI before release
* Validating _user_ input to **payloads** at the beginning of a workflow

```bash
xo schema validate --schema=path/to/schema.json --input=path/to/input.json

# or
cd massdriver-bundles
xo schema validate -s ./definitions/bundle-metadata.json -i ./bundles/$BUNDLE_NAME/metadata.json
```

**Generating variable definitions and variable files**:

[Terraform Variable Types](https://www.terraform.io/docs/configuration/expressions/types.html#types)

TBD

## Development

### Adding Commands

Add commands using the [Cobra Generator](https://github.com/spf13/cobra/blob/master/cobra/README.md).

Commands should be scoped (subcommand) under a parent "command" to facilitate organization.

Blogs on writing Cobra commands:

* https://towardsdatascience.com/how-to-create-a-cli-in-golang-with-cobra-d729641c7177