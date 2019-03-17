# Qdrouterd Operator

A Kubernetes operator to manage Qdrouterd deployments, automating mesh creation and administration

## Introduction

Qdrouterd Operator proves a `Qdrouterd` [Custom Resource Definition](https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/) (CRD) that models a Qdrouterd deployment. This CRD allows for specifying the number of messaging routers, the deployment topology as well as other options for the interconnect operation:

## Usage

Deploy the Qdrouterd Operator into the Kubernetes cluster where it will manage requests for the `Qdrouterd` resource. The Qdrouterd Operator will watch for create, update and delete resource requests and perform the necessary steps to ensure the present cluster state matches the desired state.

#### Deploy Qdrouterd Operator

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

Deploy the CRD to the cluster that defines the Qdrouterd resource.

```
$ kubectl create -f deploy/crds/interconnectedcloud_v1alpha1_qdrouterd_crd.yaml
```

Next, deploy the operator into the cluster.

```
$ kubectl create -f deploy/operator.yaml
```

This step will create a pod on the Kubernetes cluster for the Qdrouterd Operator.
Observe the Qdrouterd Operator pod and verify it is in the running state.

```
$ kubectl get pods -l name=qdrouterd-operator
```

If for some reason, the pod does not get to the running state, look at the
pod details to review any event that probited the pod from starting.

```
$ kubectl describe pod -l name=qdrouterd-operator
```

You will be able to confirm that the new CRD has been registerd in the cluster and you can review its details.

```
$ kubectl get crd
$ kubectl describd crd qdrouterd-operator
```

To create a Qdrouterd deployment, you must create a `Qdrouterd` resource representing the desired specfication of the deployment. For example, to creat a 3-node Qdrouterd mesh deployment you may run:

```console
$ cat <<EOF | kubectl create -f -
apiVersion: interconnectedcloud.github.io/v1alpha1
kind: Qdrouterd
metadata:
  name: amq-interconnect
spec:
  # Add fields here
  count: 3
  image: quay.io/interconnectedcloud/qdrouterd:1.5.0
  deploymentMode: lbfrontend
  addresses:
    - prefix: balanced
      distribution: balanced
    - prefix: closest
      distribution: closest
    - prefix: multicast
      distribution: multicast
EOF
```

The Qdrouterd Operator will act upon the creation of the resource by
creating the necessary Kubernetes resources for the desired deployment.
These resources will be monitored by the Qdrouterd Operator and will maintain
the desired state as long as the `Qdrouterd` resource exists. 

## Development

This Operator is built using the [Operator SDK](https://github.com/operator-framework/operator-sdk). Follow the [Quick Start](https://github.com/operator-framework/operator-sdk) instructions to checkout and install the operator-sdk CLI.

Local development may be done with [minikube](https://github.com/kubernetes/minikube) or [minishift](https://www.okd.io/minishift/).

#### Source Code

Clone this repository to a location on your workstation such as `$GOPATH/src/github.com/ORG/REPO`. Navigate to the repository and install the dependencies.

```
$ cd $GOPATH/src/github.com/ORG/REPO/qdrouterd-operator
$ dep ensure && dep status
```

#### Run Operator Locally

Ensure the service account, role, role bindings and CRD are added to  the local cluster.

```
$ kubectl create -f deploy/service_account.yaml
$ kubectl create -f deploy/role.yaml
$ kubectl create -f deploy/role_binding.yaml
$ kubectl create -f deploy/crds/interconnectedcloud_v1alpha1_qdrouterd_crd.yaml
```

Start the operator locally for development.

```
$ operator-sdk up local
```

Create a minimal Qdrouterd resource to observe and test your changes.

```console
$ cat <<EOF | kubectl create -f -
apiVersion: interconnectedcloud.github.io/v1alpha1
kind: Qdrouterd
metadata:
  name: amq-interconnect
spec:
  # Add fields here
  count: 3
  image: quay.io/interconnectedcloud/qdrouterd:1.5.0
  deploymentMode: lbfrontend
EOF
```

As you make local changes to the code, restart the operator to enact the changes.

