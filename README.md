# clusterclaims-controller

## Summary
When using Cluster pools, if you load this controller, it will create ManagedCluster and KlusterletAddonConfig resources, so that all you have to do is apply a simple `clusterClaim` resource.

```yaml
---
apiVersion: hive.openshift.io/v1
kind: ClusterClaim
metadata:
  name: my-cluster
  namespace: aws-east
spec:
  clusterPoolName: aws-eas
```
_Note: See `./examples/clusterclaim.yaml`_

## Deploy
To install this controller, you must be kubeadmin. The controller will be deployed to the `open-cluster-management` namespace.

```bash
oc apply -k ./deploy
```
It will take 1-2min for the image to download the first time. The controller runs two pods, and chooses a leader to reduce the possibility of an outage.

## Using ClusterClaims in GitOps
1. Create one or more cluster pools in ACM
2. In Git, commit a clusterClaim resource pointing to the pool. You can create as many claims as you want (each one will create a cluster with ACM)
3. When your done, delete the claim in Git (cluster is removed)

You can use ACM Applicaiton management or OpenShift GitOps to deliver the clusterClaim resource to the ACM hub.