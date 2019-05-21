#!/usr/bin/env bash

#
# Copyright 2019, Interconnectedcloud authors.
# License: Apache License 2.0 (see the file LICENSE or http://apache.org/licenses/LICENSE-2.0.html).
#

set -o errexit
set -o nounset
set -o pipefail

SCRIPTPATH="$(cd "$(dirname "$0")" && pwd -P)"
GENERATOR_BASE=${SCRIPTPATH}/../vendor/k8s.io/code-generator

"$GENERATOR_BASE/generate-groups.sh" "client,informer,lister" \
    github.com/interconnectedcloud/qdr-operator/pkg/client \
    github.com/interconnectedcloud/qdr-operator/pkg/apis \
    interconnectedcloud:v1alpha1 \
    --go-header-file "${SCRIPTPATH}/boilerplate.go.txt" \
    "$@"
