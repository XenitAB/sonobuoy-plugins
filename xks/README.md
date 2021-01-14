# XKS

E2E tests to validate XKS deployments

## Usage
```shell
sonobuoy run --plugin xks.yaml
```

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
