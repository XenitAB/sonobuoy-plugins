# XKS

E2E tests to validate XKS deployments

## Usage

Run the xks test suite by using the plugin manifest.
```shell
sonobuoy run --plugin xks.yaml
```

The latest manifest can also be referenced remotely.
```shell
sonobuoy run --plugin https://raw.githubusercontent.com/XenitAB/sonobuoy-plugins/main/xks/xks.yaml
```

## Tests

The tests are split up into separat

### FluxV2

Verifies that the FluxV2 components are installed and that the bootstrap `GitRepository` and `Kustomize` are in a healthy state.

### Ingress

Tests that public DNS records and valid certificates can be provisioned in a cluster.

### Pod Security

Cheks best practice security hardening is enabled in the cluster.

### High Availabiltiy

Checks critical configuration for high availibiltiy which could affect uptime of an application.

## Development

When developing it might be easier to run the tests locally against
a remote cluster. The tests are configured to run agains the current
context configure in `KUBECONFIG`.
```shell
make test
```

Running all the tests may take a while, so it may be better to only
run specific tests. The filtering is done on the test description
names with Regex.
```
make test E2E_FOCUS=FluxV2
```


It is also possible skip specific tests.
```shell
make test E2E_SKIP=FluxV2
```
