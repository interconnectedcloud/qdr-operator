#!/bin/sh

if [[ -z ${1} ]]; then
    CATALOG_NS="operator-lifecycle-manager"
else
    CATALOG_NS=${1}
fi

CSV=`cat deploy/olm-catalog/qdrouterd-operator/0.1.0/qdrouterd-operator.v0.1.0.clusterserviceversion.yaml | sed -e 's/^/          /' | sed '0,/ /{s/          /        - /}'`
CRD=`cat deploy/crds/interconnectedcloud_v1alpha1_qdrouterd_crd.yaml  | sed -e 's/^/          /' | sed '0,/ /{s/          /        - /}'`
PKG=`cat deploy/olm-catalog/qdrouterd-operator/0.1.0/interconnectedcloud.package.yaml | sed -e 's/^/          /' | sed '0,/ /{s/          /        - /}'`

cat << EOF > deploy/olm-catalog/qdrouterd-operator/0.1.0/catalog-source.yaml
apiVersion: v1
kind: List
items:
  - apiVersion: v1
    kind: ConfigMap
    metadata:
      name: qdrouterd-resources
      namespace: ${CATALOG_NS}
    data:
      clusterServiceVersions: |
${CSV}
      customResourceDefinitions: |
${CRD}
      packages: >
${PKG}

  - apiVersion: operators.coreos.com/v1alpha1
    kind: CatalogSource
    metadata:
      name: qdrouterd-resources
      namespace: ${CATALOG_NS}
    spec:
      configMap: qdrouterd-resources
      displayName: Qdrouterd Operators
      publisher: Red Hat
      sourceType: internal
    status:
      configMapReference:
        name: qdrouterd-resources
        namespace: ${CATALOG_NS}
EOF
