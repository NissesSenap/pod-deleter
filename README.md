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

- Do we need tracing/metrics?
- How should we handle input value for in-cluster or not?
  - Could do it the other way and try if we can run it inside the cluster, if not just go with the other way...
- Verify that pod Id is correct vs falco event information, would be nice to match pod ID vs the output of k8s when getting the pod

### block list

If anything is defined in allow list it will automatically make block list not needed since it's only those namespaces that should have pods
that you can delete from.

Block list is defined by default but a allow list isn't even created unless needed.

### Implementations

Should we parse cloudEvents as well?

- OpenFaas
- Lambda

## Simple e2e test

Set test event

```shell
export BODY='{"output":"14:49:49.264147779: Notice A shell was spawned in a container with an attached terminal (user=root user_loginuid=-1 k8s.ns=default k8s.pod=alpine container=a15057582acc shell=sh parent=runc cmdline=sh -c uptime terminal=34816 container_id=a15057582acc image=alpine) k8s.ns=default k8s.pod=alpine container=a15057582acc k8s.ns=default k8s.pod=alpine container=a15057582acc","priority":"Notice","rule":"Terminal shell in container","time":"2021-05-01T14:49:49.264147779Z", "output_fields": {"container.id":"a15057582acc","container.image.repository":"alpine","evt.time":1619880589264147779,"k8s.ns.name":"default","k8s.pod.name":"alpine","proc.cmdline":"sh -c uptime","proc.name":"sh","proc.pname":"runc","proc.tty":34816,"user.loginuid":-1,"user.name":"root"}}'
```

Start alpine pod and test against it:

```shell
kubectl run alpine --namespace default --image=alpine --restart='Never' -- sh -c "sleep 600"
kubectl annotate pod alpine -n default 'falco.org/protected=true'

# Verify the annotation
kubectl get pods -n alpine -n default -o yaml |grep annotations -A 3
# Run the application, it shouldn't delete anything
go run main.go


kubectl annotate pod alpine -n default 'falco.org/protected=True' --overwrite
# Run the application, it still shouldn't delete anything, notice that it will say true and not True
go run main.go

kubectl annotate pod alpine -n default 'falco.org/protected=false' --overwrite
# Run the application, it SHOULD delete the pod alpine in namespace default
go run main.go
```

### e2e allow

```shell
export ALLOW_NAMESPACE='default'
kubectl run alpine --namespace default --image=alpine --restart='Never' -- sh -c "sleep 600"


```

### critical namespaces

Pick one of them, you can't have both.

```shell
export BLOCK_NAMESPACE='block-ns1 block-ns2'
export ALLOW_NAMESPACE='app-dev-* service-x'

```
