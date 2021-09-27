# xo - eXecution Orchestrator

`xo` is a small tool to provisioning and developing Massdriver bundles.

## Development

### Building

XO is built using go:

```shell
go build
```

If you encounter an error: `fatal: could not read Username for 'https://github.com': terminal prompts disabled`, you can globally configure git:

```shell
git config --global --add url."git@github.com:".insteadOf "https://github.com/"
go build
```

### Adding Commands

Add commands using the [Cobra Generator](https://github.com/spf13/cobra/blob/master/cobra/README.md).

Commands should be scoped (subcommand) under a parent "command" to facilitate organization.

Blogs on writing Cobra commands:

* https://towardsdatascience.com/how-to-create-a-cli-in-golang-with-cobra-d729641c7177
