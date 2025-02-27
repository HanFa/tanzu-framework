# tanzu-auth-controller-manager

This directory is home to the `tanzu-auth-controller-manager`.

To run this on a cluster:

```sh
./hack/run.sh
```

To run the tests:

```sh
./hack/test.sh
```

To run static checks (e.g., `go fmt`, `go vet`, `golangci-lint`) and the tests:

```sh
./hack/check.sh
```

To generate the default pinniped addon secret:

```sh
# The primary use case of this feature is through the top-level Tanzu Framework Makefile, via a command like:
# (consult the top-level Makefile for further details)
make generate-package-secret PACKAGE=pinniped tkr=v1.23.3---vmware.1-tkg.1 iaas=vsphere

# to generate the secret in the context of the /pinniped-components/tanzu-auth-controller-manager, do one of the following:

# manually via ytt
ytt -f ./hack/ytt -v tkr=v1.23.3---vmware.1-tkg.1 -v infrastructure_provider=vsphere

# using the generate script
# arguments are ytt args passed as: -v <arg> for example, something like:
./hack/generate-package-secret.sh -v tkr=v1.23.3---vmware.1-tkg.1 -v infrastructure_provider=vsphere
```

To change the log level, add the `--v=LOG_LEVEL` arg to the controller deployment.  LOG_LEVEL should
be a number.  Default log level is 0. Example:

```yaml
---
kind: Deployment
apiVersion: apps/v1
metadata:
  name: tanzu-auth-controller-manager
  namespace: tanzu-auth
spec:
  selector:
    matchLabels:
      app: tanzu-auth-controller-manager
  template:
    metadata:
      labels:
        app: tanzu-auth-controller-manager
    spec:
      serviceAccountName: tanzu-auth-controller-manager-sa
      containers:
      - args:
        - --v=1
        image: #@ data.values.image
        name: tanzu-auth-controller-manager
```
