#!/usr/bin/env bash

if [[ -z ${1} ]]; then
    CATALOG_NS="operator-lifecycle-manager"
else
    CATALOG_NS=${1}
fi

CSV=`cat deploy/olm-catalog/qdr-operator/0.2.0/qdr-operator.v0.2.0.clusterserviceversion.yaml | sed -e 's/^/          /' | sed '0,/ /{s/          /        - /}'`
CRD=`cat deploy/crds/interconnectedcloud_v1alpha1_interconnect_crd.yaml  | sed -e 's/^/          /' | sed '0,/ /{s/          /        - /}'`
PKG=`cat deploy/olm-catalog/qdr-operator/0.2.0/interconnectedcloud.package.yaml | sed -e 's/^/          /' | sed '0,/ /{s/          /        - /}'`

cat << EOF > deploy/olm-catalog/qdr-operator/0.2.0/catalog-source.yaml
apiVersion: v1
kind: List
items:
  - apiVersion: v1
    kind: ConfigMap
    metadata:
      name: qdr-resources
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
      name: qdr-resources
      namespace: ${CATALOG_NS}
    spec:
      configMap: qdr-resources
      displayName: Qdr Operators
      publisher: Red Hat
      sourceType: internal
    status:
      configMapReference:
        name: qdr-resources
        namespace: ${CATALOG_NS}
EOF
