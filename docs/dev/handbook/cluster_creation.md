Cluster Creation with Cluster-class
========================================

This developer handbook describe how to create a cluster using a clusterclass based CRD. Note: it only works with
any tanzu-framework commit after the [PR](https://github.com/vmware-tanzu/tanzu-framework/pull/2581).

First, make sure you have the push access to an OCI registry. Export its URL as `OCI_REGISTRY` environment variable.

### OCI Registry with GCP

For example, let's suppose we have a project with ID `my-project-1527816345739` on GCP. Then we can set the registry URL
as `gcr.io/my-project-1527816345739/tkg/management`. Please refer to the
[GCP Registry documentation](https://cloud.google.com/container-registry/docs/overview) to decide the actual
URL.

Export the registry URL as `OCI_REGISTRY`
```shell
export OCI_REGISTRY=gcr.io/my-project-1527816345739/tkg/management
```







