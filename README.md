# xo - eXecution Orchestrator

`xo` is a small utility to facilitate communication between a Massdriver provisioner and the Massdriver platform.

Many of the capabilities here also exist in `mass`, the [Massdriver CLI](https://github.com/massdriver-cloud/mass). You should use `mass` locally for any interactions with the Massdriver platform. `xo` is only intended to be used within provisioners.

## Usage

### Pulling a Bundle

```bash
xo bundle pull
```

Used by the [initialization provisioner](https://github.com/massdriver-cloud/provisioner-init) to pull the bundle for all the subsequent steps.

### Managing Artifacts

This should be included in the provisioner script. Refer [Massdriver's offical provisioners](https://docs.massdriver.cloud/provisioners/overview) for examples.

```bash
xo artifact publish -d artifact_name -n "Artifact description" -f artifact.json
```

```bash
xo artifact delete -d artifact_name -n "Artifact description"
```

### Reporting Deployment Status

#### Deployment Started

This is handled automatically by the [initialization provisioner](https://github.com/massdriver-cloud/provisioner-init).

```bash
xo deployment provision start
```

```bash
xo deployment decommission start
```

#### Deployment Completed

This is handled automatically by Massdriver's workflow orchestrator.

```bash
xo deployment provision complete
```

```bash
xo deployment decommission complete
```

#### Deployment Failed

This is handled automatically by Massdriver's workflow orchestrator.

```bash
xo deployment provision fail
```

```bash
xo deployment decommission fail
```