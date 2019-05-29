# QDR Operator

A Kubernetes operator for managing Layer 7 (e.g. Application Layer) addressing and routing within and across clusters. The operator manages *interior* and *edge* QDR deployments automating resource creation and administration.

## Introduction

This operator provides an `Interconnect` [Custom Resource Definition](https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/)
(CRD) that describes a deployment of [Apache Qpid Dispatch Router](https://qpid.apache.org/components/dispatch-router/index.html) allowing developers to expertly
deploy routers for their infrastructure and middle-ware oriented  messaging requirements. The number of messaging routers, the deployment topology, address
semantics and other desired options can be specified through the CRD.

## Usage

Deploy the QDR Operator into the Kubernetes cluster where it will manage requests for the `Interconnect` resource. The QDR Operator will watch for create, update and delete resource requests and perform the necessary steps to ensure the present cluster state matches the desired state.

### Deploy QDR Operator

The `deploy` directory contains the manifests needed to properly install the
Operator.

Create the service account for the operator.

```
$ kubectl create -f deploy/service_account.yaml
```

Create the RBAC role and role-binding that grants the permissions
necessary for the operator to function.

```
$ kubectl create -f deploy/role.yaml
$ kubectl create -f deploy/role_binding.yaml
```

Deploy the CRD to the cluster that defines the Interconnect resource.

```
$ kubectl create -f deploy/crds/interconnectedcloud_v1alpha1_interconnect_crd.yaml
```

Next, deploy the operator into the cluster.

```
$ kubectl create -f deploy/operator.yaml
```

This step will create a pod on the Kubernetes cluster for the QDR Operator.
Observe the `qdr-operator` pod and verify it is in the running state.

```
$ kubectl get pods -l name=qdr-operator
```

If for some reason, the pod does not get to the running state, look at the
pod details to review any event that prohibited the pod from starting.

```
$ kubectl describe pod -l name=qdr-operator
```

You will be able to confirm that the new CRD has been registered in the cluster and you can review its details.

```
$ kubectl get crd
$ kubectl describe crd interconnects.interconnectedcloud.github.io
```

To create a router deployment, you must create a `Interconnect` resource representing the desired specification of the deployment. For example, to create a 3-node router mesh deployment you may run:

```console
$ cat <<EOF | kubectl create -f -
apiVersion: interconnectedcloud.github.io/v1alpha1
kind: Interconnect
metadata:
  name: example-interconnect
spec:
  # Add fields here
  deploymentPlan:
    image: quay.io/interconnectedcloud/qdrouterd:1.7.0
    role: interior
    size: 3
    placement: Any
EOF
```

The operator will create a deployment of three router instances, all connected together with default address semantics. It will also create a service through which the *interior* router mesh can be accessed. It will configure a default set of *listeners* and *connectors* as described below. You will be able to confirm that the instance has been created in the cluster and you can review its details. To view the Interconnect instance, the deployment it manages and the associated pods that are deployed:

```
$ kubectl describe interconnect example-interconnect
$ kubectl describe deploy example-interconnect
$ kubectl describe svc example-interconnect
$ kubectl get pods -o yaml
```

### Deployment Plan

The CRD *Deployment Plan* defines the attributes for an Interconnect instance.

#### Role and Placement

The *Deployment Plan* **Role** defines the mode of operation for the routers in a topology.

  * **interior** - This role creates an interconnect of auto-meshed routers for concurrent connection capacity and resiliency.
    Connectivity between the routers will be defined by *InterRouterListeners* and *InterRouterConnectors*. Downlink
    connectivity with *edge* routers will be via *EdgeListeners*.

  * **edge** -  This role creates a set of stand-alone routers. The connectivity from the *edge* to *interior* routers
    will be via *EdgeConnectors*.

The *Deployment Plan* **Placement** defines the deployment resource and the associated scheduling of the pods in the cluster.

  * **Any** - There is no constraint on pod placement. The operator will manage a *Deployment* resource where the number of
    pods scheduled will be up to *Deployment Plan Size*.

  * **Every** - Router pods will be placed on each node in the cluster. The operator will manage a *DaemonSet* resource where the
    number of pods scheduled will correspond to the number of nodes in the cluster. Note the *Deployment Plan Size* is
    not used.

  * **Anti-Affinity** - This constrains scheduling and prevents multiple router pods from running on the same node in the
    cluster. The operator will manage a *Deployment* resource with number of pods up to *Deployment Plan Size*.
    Note if *Deployment Plan Size* is greater than the number of nodes in the cluster, the excess pods that cannot be
    scheduled will remain in the *pending* state.

### Connectivity

The connectivity between routers in a deployment is via the declared *listeners* and *connectors*. There are three types of *listeners* supported by the operator.

  * **Listeners** - A listener for accepting normal messaging client connections. The operator supports this listener for
    both *interior* and *edge* routers.

  * **InterRouterListeners** - A listener for accepting connections from peer *interior* routers. The operator
    support this listener for *interior* routers **only**.

  * **EdgeListeners** - A listener for accepting connections from downlink *edge* routers. The operator supports this
    listener for *interior* routers **only**.

There are three types of *connectors* supported by the operator.

  * **Connectors** - A connector for connecting to an external messaging intermediary. The operator supports this connector
    for both *interior* and *edge* routers.

  * **InterRouterConnectors** - A connector for establishing connectivity to peer *interior* routers. The operator
    supports this connector for *interior* routers **only**.

  * **EdgeConnector** - A connector for establishing up-link connectivity from *edge* to *interior* routers. The operator
    supports this connector for *edge* routers **only**.

## Development

This Operator is built using the [Operator SDK](https://github.com/operator-framework/operator-sdk). Follow the [Quick Start](https://github.com/operator-framework/operator-sdk) instructions to checkout and install the operator-sdk CLI.

Local development may be done with [minikube](https://github.com/kubernetes/minikube) or [minishift](https://www.okd.io/minishift/).

#### Source Code

Clone this repository to a location on your workstation such as `$GOPATH/src/github.com/ORG/REPO`. Navigate to the repository and install the dependencies.

```
$ cd $GOPATH/src/github.com/ORG/REPO/qdr-operator
$ dep ensure && dep status
```

#### Run Operator Locally

Ensure the service account, role, role bindings and CRD are added to  the local cluster.

```
$ kubectl create -f deploy/service_account.yaml
$ kubectl create -f deploy/role.yaml
$ kubectl create -f deploy/role_binding.yaml
$ kubectl create -f deploy/crds/interconnectedcloud_v1alpha1_interconnect_crd.yaml
```

Start the operator locally for development.

```
$ operator-sdk up local
```

Create a minimal Interconnect resource to observe and test your changes.

```console
$ cat <<EOF | kubectl create -f -
apiVersion: interconnectedcloud.github.io/v1alpha1
kind: Interconnect
metadata:
  name: example-interconnect
spec:
  deploymentPlan:
    image: quay.io/interconnectedcloud/qdrouterd:1.7.0
    role: interior
    size: 3
    placement: Any
EOF
```

As you make local changes to the code, restart the operator to enact the changes.

#### Build

The Makefile will do the dependency check, operator-sdk generate k8s, run local test, and finally the operator-sdk build. Please ensure any local docker server is running.

```
make
```

#### Test

Before submitting PR, please test your code. 

File or local validation.
```
$ make test
```

Cluster-based test. 
Ensure there is a cluster running before running the test.

```
$ make cluster-test
```

## Manage the operator using the Operator Lifecycle Manager in OpenShift 4.0 or above
To install this operator on OpenShift 4 for end-to-end testing, make sure you have access to a quay.io account to create an application repository. Follow the [authentication](https://github.com/operator-framework/operator-courier/#authentication) instructions for Operator Courier to obtain an account token. This token is in the form of "basic XXXXXXXXX" and both words are required for the command.
Mainly you need to install `operator-courier` and generate the quay token before you proceed to the next steps.

Edit file `deploy/olm-catalog/courier/qdr-operatorsource.yaml`, Remember to replace `registryNamespace` with your quay username id. The name, display name and publisher of the operator are the only other attributes that may be modified.

Edit file `deploy/olm-catalog/courier/bundle_dir/<version>/interconnectedcloud.package.yaml`, and point the `packagename` to your quay application namespace.

Push the operator bundle to your quay application repository as seen below:

```bash
operator-courier --verbose push deploy/olm-catalog/courier/bundle_dir/0.1.0 <quay_username_id> qdrapp-operator 0.1.0 "basic XXXXXXXXX"
```

Note that the push command does not overwrite an existing repository, and it needs to be deleted before a new version can be built and uploaded. 

Once the bundle has been uploaded, create an [Operator Source](https://github.com/operator-framework/community-operators/blob/master/docs/testing-operators.md#linking-the-quay-application-repository-to-your-openshift-40-cluster) or use the command below to load your operator bundle in OpenShift.

```bash
oc create -f deploy/olm-catalog/courier/qdr-operatorsource.yaml
```
It will take a few minutes for the operator to become visible under the _OperatorHub_ section of the OpenShift console _Catalog_. It can be easily found by filtering the provider type to _Custom_.

## Manage the operator using the Operator Lifecycle Manager below Openshift 4.0

Ensure the Operator Lifecycle Manager is installed in the local cluster.  By default, the `catalog-source.sh` will install the operator catalog resources in `operator-lifecycle-manager` namespace.  You may also specify different namespace where you have the Operator Lifecycle Manager installed.

```
$ ./hack/catalog-source.sh <namespace>
$ oc apply -f deploy/olm-catalog/qdr-operator/0.1.0/catalog-source.yaml
```
