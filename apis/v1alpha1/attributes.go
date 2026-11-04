/*
Copyright (c) 2026 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	resourceapi "k8s.io/api/resource/v1"
)

const (
	// StandardDeviceAttributePrefix is the prefix used for standard device attributes.
	StandardDeviceAttributePrefix = "multinetwork.networking.x-k8s.io"

	// StandardDeviceAttributePodNetwork is a standard device attribute name
	// which describes a pod network.
	// The value is a string value referring to the name of an object from the GK defined in the
	// StandardDeviceAttributePodNetworkKind attribute.
	StandardDeviceAttributePodNetwork resourceapi.QualifiedName = StandardDeviceAttributePrefix + "/" + "podNetwork"
	// StandardDeviceAttributePodNetworkNamespace is a standard device attribute name
	// which describes the namespace of a pod network.
	// The value is a string value referring to the namespace of a pod network object.
	// The attribute should not be set for the PodNetworkKind pointing to a non-namespaced GK.
	// The attribute is mandatory for the PodNetworkKind pointing to a namespaced GK.
	StandardDeviceAttributePodNetworkNamespace resourceapi.QualifiedName = StandardDeviceAttributePrefix + "/" + "podNetworkNamespace"
	// StandardDeviceAttributePodNetworkKind is a standard device attribute name
	// which describes a PodNetworkKind.
	// The value is a string value referring to an existing PodNetworkKind object.
	StandardDeviceAttributePodNetworkKind resourceapi.QualifiedName = StandardDeviceAttributePrefix + "/" + "podNetworkKind"
)
