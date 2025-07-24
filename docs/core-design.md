# Multi-Network design

## Goal

- Define the network handler for Kubernetes clusters.
- Define behaviours for the handler.
- Define integration and interactions with [Dynamic Resource Allocation](https://kubernetes.io/docs/concepts/scheduling-eviction/dynamic-resource-allocation) (DRA).
- Standardize APIs for driver implementations.

## Non-Goal

- Provide implementation for the driver. See [https://github.com/kubernetes-sigs/cni-dra-driver](https://github.com/kubernetes-sigs/cni-dra-driver) for an example implementation.
- Standardize network configuration.

## Terminology

**Default PodNetwork** \- This is the initial cluster-wide PodNetwork provided during cluster creation that is available to the Pod when no additional networking configuration is provided in Pod spec.
**Primary PodNetwork** \- This is the PodNetwork inside the Pod’s network namespace whose interface is used for the default gateway (0.0.0.0) in the default virtual routing and forwarding table (VRF).
**DRA** \- Dynamic Resource Allocation is an API allowing device sharing across multiple workloads (see the official [page](https://kubernetes.io/docs/concepts/scheduling-eviction/dynamic-resource-allocation/) for more details)

## Background

### API and Implementation

This design focuses mainly on the user-facing APIs definitions and how these are going to interact with DRA.
The Implementation component is the DRA driver delivered by the provider that handles the network interface for the pod, based on the DRA configuration. It also integrates with the DRA scheduling mechanism based on the implementation specific logic. This area is fully owned by the provider for the Multi-Network capability. This proposal will try to standardize the API contact points, but not the implementation itself.

### Network

“Network” word is a very overloaded term, and for many might have different meaning: it might be represented as a VLAN, or a specific interface in a Node, or identified as a unique IP subnet (a CIDR). In this design we do not want to limit to one definition, and ensure that we are flexible enough to satisfy most definitions.

## Overview

This design introduces a new object called “PodNetwork”. This object will be a handle representing a network, Pods can connect to. PodNetwork will hold centralized configuration for these endpoints.
The PodNetwork object will be an input to the DRA driver, for it to provide the “podNetwork” attributed devices.
Application developers (end users) will leverage DRA’s APIs (ResourceClaim, ResourceClaimTemplate) to reference these devices in a common way.

<p align="center">
  <img src="images/image7.png?raw=true"/>
</p>

### Positioning in DRA

To better understand how Multi-Network will function with DRA we want to define 3 main focus areas:

1. Direct DRA device selection \- here are the AI/ML use cases, or simple device handling for which users do not need centralized configuration.
2. PodNetwork as Cluster-Level construct to enable orchestration of a common set of configurations \- this gives cluster administrators means to easily maintain uniformity. This design focuses on this area.
3. PodNetwork as a handle to enable Kubernetes networking primitives: Services, NetworkPolicies, etc.

<p align="center">
  <img src="images/image8.png?raw=true"/>
</p>

## Detailed design

### PodNetwork CRD

This design defines a PodNetwork object. This will be a CustomResourceDefinition (CRD). PodNetwork’s form is generic and vendor-neutral. This API object is what we wish to use for the future integration with the rest of the Kubernetes ecosystem.

#### Motivation

As of Today, there is no core object that represents the network that the Pod connects to. It represents the administrative resource at the cluster scope for configuration governance for the Pod connection, which provides the cluster administrators centralized control over what is used by the cluster workloads. Furthermore, this will allow workloads to select more than 1 such network configuration to attach to. In such situations, we want to be able to differentiate where specific Kubernetes functionality needs to be applied, e.g. serviceability which uses labels and namespace placement as the identifiers of what workloads to apply to, this object will be another indicator which interfaces of a Pod to use.

Lastly, this object does not try to standardize the configuration of PodNetwork on how it is attaching interfaces to the Pod. This is fully owned by the implementation. In the future, if we will see common fields or functionalities occuring across the implementers, we will try to introduce them as a common part of the API.

##### Alternatives

DRA DeviceClass \- one of the considerations was to re-use the DeviceClass object that comes with the DRA. The problem of reusing DeviceClass is that it is is not strictly a network-specific object. DeviceClass can be used for any resource, e.g. a GPU. This will introduce confusion in integrating it with other Kubernetes functionality related to networking.

#### API

PodNetwork object is described as follows:

```go
// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// PodNetwork represents a logical network in Kubernetes Cluster.
type PodNetwork struct {
        metav1.TypeMeta   `json:",inline"`
        metav1.ObjectMeta `json:"metadata,omitempty"`

        Spec   PodNetworkSpec   `json:"spec"`
        Status PodNetworkStatus `json:"status,omitempty"`
}

// PodNetworkSpec contains the specifications for podNetwork object
type PodNetworkSpec struct {

        // Enabled is used to administratively enable/disable a PodNetwork.
        // When set to false, PodNetwork Ready condition will be set to False.
        // Defaults to True.
        //
        // +kubebuilder:default=true
        Enabled bool `json:"enabled,"`

        // Provider is the name of the handler of this PodNetwork.
        // The value of this field MUST be a domain prefixed path.
        //
        // Example: "example.net/multi-network".
        //
        // This field is not mutable and cannot be empty.
        //
        //
        // +kubebuilder:validation:MinLength=1
        // +kubebuilder:validation:MaxLength=253
        // +kubebuilder:validation:XValidation:message="Value is immutable",rule="self == oldSelf"
        // +required
        Provider string `json:"provider,"`

        // Parameters define provider implementation specific configuration.
        //
        // +optional
        Parameters runtime.RawExtension `json:"parameters,omitempty"`
}

// PodNetworkStatus contains the status information related to the PodNetwork.
type PodNetworkStatus struct{
        // Conditions describe the current conditions of the PodNetwork.
        //
        // Known condition types are:
        // * "Ready"
        //
        // +optional
        // +listType=map
        // +listMapKey=type
        // +kubebuilder:validation:MaxItems=5
        // +kubebuilder:default={{type: "Ready", status: "Unknown", reason:"Pending", message:"Waiting for controller", lastTransitionTime: "1970-01-01T00:00:00Z"}}
        Conditions []metav1.Condition `json:"conditions,omitempty"`
}
```

##### Group

We will create a new api group called `multinetwork`. Currently it will be placed in the `networking.x-k8s.io` API group.

Eventually the goal will be to place this API under `networking.k8s.io` after securing the appropriate approvals.

#### Examples

##### Generic

```yaml
apiVersion: multinetwork.networking.k8s.io/v1
kind: PodNetwork
metadata:
  name: net-dataplane
spec:
  enabled: true
  provider: "foo.io/bar"
status:
  conditions:
  - lastProbeTime: null
    lastTransitionTime: "2025-01-31T18:38:01Z"
    status: "True"
    type: Ready
```

##### With implementation-specific config

```yaml
apiVersion: multinetwork.networking.k8s.io/v1
kind: PodNetwork
metadata:
  name: dataplane
spec:
  enabled: true
  provider: "example.dra.x-k8s.io/network"
  parameters:
    master: "eth0"
    mode: "bridge"
    ipam:
      type: "host-local"
      subnet: "192.168.1.0/24"
      rangeStart: "192.168.1.200"
      rangeEnd: "192.168.1.216"
      routes:
      - dst: "0.0.0.0/0"
      gateway: "192.168.1.1"
status:
  conditions:
  - lastProbeTime: null
    lastTransitionTime: "2022-11-17T18:38:01Z"
    status: "True"
    type: Ready
```

#### Status Conditions

The PodNetwork object will use following conditions:

* **Ready** \- indicates that the PodNetwork object is correct (validated) and can be consumed by the DRA Driver. New Pods should not be attached to a PodNetwork that is not Ready. This condition does not indicate readiness of specific PodNetwork on a per Node-basis. Following are the error reasons for this condition:


| Reason name | Description |
| ----- | ----- |
| ConditionsNotReady | Other conditions on the object have “false” value. |
| AdministrativelyDisabled | The PodNetwork’s Enabled field is set to false. |

The provider’s implementation is responsible for the conditions life-cycle.

##### Other Conditions

We will not restrict adding other conditions to this object. This way implementations can further control the Readiness of PodNetwork by adding their own Condition and logic for it.

##### Implementation consideration

Implementations must enforce the Ready Condition from PodNetwork. It is up to the implementation for the solution, e.g. withdraw resource availability for the specific PodNetwork across Nodes.

#### Provider

The provider field gives the ability to uniquely identify what implementation (provider) is going to handle specific instances of PodNetwork objects. The value has to be a domain prefixed path.

#### Enabled

The Enabled field is created to allow administrative control for a given PodNetwork (e.g. for migration). When set to false no new Pods can be attached to that PodNetwork.

#### Parameters

Parameters field is a free form configuration for the implementations to use and define how these fields are used.

#### InUse state

It is recommended to put PodNetwork objects that are referenced by at least one Pod in an InUse state. It is up to the implementation to do it and how that is implemented.

#### Lifecycle

A PodNetwork will be in following phases:

1. Created/NotReady \- when PodNetwork is just created its Ready condition is false. Pods can reference such a PodNetwork, but will be in Pending state until the PodNetwork becomes Ready.
2. Ready \- when validation of the PodNetwork succeeded and Ready condition is set to true. Here Pods can start attaching to a PodNetwork.
3. InUse \- when there is a Pod that references a given PodNetwork.
4. Disabled \- when the user sets the Enabled field to false. Such PodNetwork should be NotReady, and no new Pods will be able to attach to such PodNetwork.

### DRA integration

This design defines DRA driver implementation standardization. This serves to allow unification on the usability of Multi-Network within Kubernetes Cluster and will allow for future features integrations.

#### Device PodNetwork attribute

This design standardizes for all drivers supporting PodNetwork to provide the “\<driver name\>/podNetwork” attribute for the corresponding device listed in DRA’s ResourceSlice object.

##### Example

Below is an example of a device belonging to PodNetwork “blue-net”.

```yaml
apiVersion: resource.k8s.io/v1beta1
kind: ResourceSlice
metadata:
  name: node1-example.dra.x-k8s.io-79gzs
spec:
  devices:
  - name: eno1
    basic:
      attributes:
        example.dra.x-k8s.io/podNetwork:
          string: blue-net
```

#### ResourceClaim Status

This design standardizes for all drivers supporting PodNetwork to properly set the “status.devices.\[\].data” with the podNetwork field. The structure for the data should include following definition:

```go
type podStatus struct {
	PodNetwork string `json:"podNetwork"`
}
```

##### Example

Below is an example of a ResourceClaim with the correct status set for a device belonging to PodNetwork “blue-net”.

```yaml
apiVersion: resource.k8s.io/v1beta1
kind: ResourceClaim
metadata:
  name: my-rc
spec:
…
status:
  devices:
  - device: eno1
    driver: example.dra.x-k8s.io
    data:
      podNetwork: blue-net
```

### Attaching PodNetwork to a Pod

With above standardization, users can use following example on how to reference specific PodNetwork in the DRA ResourceClaim:

```yaml
apiVersion: resource.k8s.io/v1beta1
kind: ResourceClaim
metadata:
  name: my-rc
spec:
  devices:
    requests:
    - name: req-blue
      deviceClassName: example.dra.x-k8s.io
      selectors:
      - cel:
          expression: |-
            has(device.attributes["example.dra.x-k8s.io"].podNetwork) &&
            device.attributes["example.dra.x-k8s.io"].podNetwork == "blue-net"
```

Such configuration can be pushed to DRA’s DeviceClass as well, to simplify the user-facing claim objects.

### Conformance tests

We will add conformance tests for all the behaviours mentioned in this design.
