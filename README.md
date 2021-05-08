# pod-deleter

A PoC to build a more generic pod deleter that we can use with falcoexporter.

It's meant to be a agnostic go application that can be used in and outside a kubernetes cluster.

For example:

- AWS lambda
- OpenFaas
- tekton
- and many more.

## TODO

Should we use another log solution like "github.com/go-logr/logr"?

- There should be a interface so we can use any k8s client.
- turn podDeleter in to a interface
- Use env variabels to read criticalNamespace

### Implementations

Should we parse cloudEvents as well?

- OpenFaas
- Lambda
