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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:resource:shortName=pnk,scope=Cluster
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Group",type=string,JSONPath=`.spec.implementationType.group`
// +kubebuilder:printcolumn:name="Kind",type=string,JSONPath=`.spec.implementationType.kind`
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodNetworkKind describes a kind of pod networks implemented
// by a specific controller and represented by implementation-defined objects.
type PodNetworkKind struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// spec is the desired state of the PodNetworkKind.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Spec PodNetworkKindSpec `json:"spec"`
	// status is the current state of the PodNetworkKind.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Status PodNetworkKindStatus `json:"status"`
}

// PodNetworkKindSpec describes how pod network objects of this kind are identified.
type PodNetworkKindSpec struct {
	// ImplementationType identifies the API type of the pod network objects
	// belonging to this PodNetworkKind.
	// The ImplementationType may reference either a namespace-scoped or a cluster-scoped resource type.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="spec.implementationType is immutable"
	ImplementationType metav1.GroupKind `json:"implementationType"`
}

// PodNetworkKindStatus describes the observed state of the PodNetworkKind.
type PodNetworkKindStatus struct {
	// conditions is the list of conditions for this PodNetworkKind.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// Well-known condition types for PodNetworkKinds.
const (
	// PodNetworkKindConditionImplementationTypeReady indicates whether the resource type referenced by the
	// PodNetworkKind's ImplementationType exists and is ready in the cluster.
	PodNetworkKindConditionImplementationTypeReady = "ImplementationTypeReady"
)

// Well-known condition reasons for PodNetworkKinds.
const (
	// PodNetworkKindReasonResourceTypeNotFound in the ImplementationTypeReady condition indicates
	// that the resource type referenced by the PodNetworkKind's ImplementationType does not exist
	// in the cluster.
	PodNetworkKindReasonResourceTypeNotFound string = "ResourceTypeNotFound"
	// PodNetworkKindReasonMissingRBAC in the ImplementationTypeReady condition indicates
	// that the necessary RBAC permissions are missing for the resource type referenced by the
	// PodNetworkKind's ImplementationType.
	PodNetworkKindReasonMissingRBAC string = "MissingRBAC"
	// PodNetworkKindReasonTypeConflict in the ImplementationTypeReady condition indicates
	// that another PodNetworkKind already references the same ImplementationType (Group/Kind).
	// A PodNetworkKind is considered valid only if it is strictly the oldest object by
	// metadata.creationTimestamp; any newer PodNetworkKind receives this reason.
	// If multiple PodNetworkKinds share the same creationTimestamp, all of them receive this reason.
	PodNetworkKindReasonTypeConflict string = "TypeConflict"
	// PodNetworkKindReasonCompliant in the ImplementationTypeReady condition indicates
	// that the resource type referenced by the PodNetworkKind's ImplementationType exists and is ready.
	PodNetworkKindReasonCompliant string = "Compliant"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodNetworkKindList is a list of PodNetworkKind resources.
type PodNetworkKindList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []PodNetworkKind `json:"items"`
}

// PodNetwork identifies a specific pod network instance.
type PodNetwork struct {
	// Kind identifies the PodNetworkKind responsible for the pod network.
	// +optional
	Kind string `json:"kind,omitempty"`

	// Name identifies the pod network object.
	Name string `json:"name"`

	// Namespace identifies the namespace of the pod network object.
	// Not set if the pod network object is a non-namespaced resource.
	// +optional
	Namespace *string `json:"namespace,omitempty"`
}

// DeviceStatus represents the structure pointing to the pod network in ResourceClaim.status.devices[].data.
type DeviceStatus struct {
	// PodNetwork identifies the pod network instance attached to this device.
	PodNetwork *PodNetwork `json:"podNetwork"`
}
