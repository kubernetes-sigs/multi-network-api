#!/usr/bin/env bash

# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
OUTPUT_PKG="sigs.k8s.io/multi-network-api"

go run sigs.k8s.io/controller-tools/cmd/controller-gen \
	crd \
	paths="${OUTPUT_PKG}/apis/..." \
	output:crd:artifacts:config=deployment

go run sigs.k8s.io/controller-tools/cmd/controller-gen \
	object:headerFile="${SCRIPT_ROOT}/hack/boilerplate.go.txt" \
	paths="./..."

go run k8s.io/code-generator/cmd/client-gen \
	--clientset-name "versioned" \
	--input-base "" \
	--input "${OUTPUT_PKG}/apis/v1alpha1" \
	--output-dir "pkg/client/clientset" \
	--output-pkg "${OUTPUT_PKG}/pkg/client/clientset" \
	--go-header-file ${SCRIPT_ROOT}/hack/boilerplate.go.txt

go run k8s.io/code-generator/cmd/lister-gen \
	--output-dir "pkg/client/listers" \
	--output-pkg "${OUTPUT_PKG}/pkg/client/listers" \
	--go-header-file ${SCRIPT_ROOT}/hack/boilerplate.go.txt \
	"${OUTPUT_PKG}/apis/v1alpha1"

go run k8s.io/code-generator/cmd/informer-gen \
	--versioned-clientset-package "${OUTPUT_PKG}/pkg/client/clientset/versioned" \
	--listers-package "${OUTPUT_PKG}/pkg/client/listers" \
	--output-dir "pkg/client/informers" \
	--output-pkg "${OUTPUT_PKG}/pkg/client/informers" \
	--go-header-file ${SCRIPT_ROOT}/hack/boilerplate.go.txt \
	"${OUTPUT_PKG}/apis/v1alpha1"

go run k8s.io/code-generator/cmd/register-gen \
	--output-file "zz_generated.register.go" \
	"${OUTPUT_PKG}/apis/v1alpha1"
