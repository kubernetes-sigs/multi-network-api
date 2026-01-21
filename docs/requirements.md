# Multi-Network overview

This doc is trying to provide motivation, use cases and requirements for the Multi-Network project. This doc is not final and any future changes are acceptable.

## Terminology

**Default PodNetwork** \- This is the initial cluster-wide network that pods are, by default, attached to.

### Network

Pod's network interfaces are attached to a Network. A single Network can have many attached Pods (1 to many).

Network is an abstract entity, but we expect:

* Interfaces attached to the same Network within a cluster have unique IP addresses but interfaces attached to different networks could have overlapping IP addresses.
* Interfaces attached to the same Network are by default mutually reachable (modulo policy). This is analogous to the standard Pod reachability.
* Users had an ability to access functionality similar to Kubernetes ClusterIP services and Network Policy but scoped to a given Network. Note that the APIs used will not be the core Service and NetworkPolicy APIs. Functionally, traffic that is routed (via the interface) to the Network can have such features applied to it.

Most important is what is not explicitly specified:

* Kind of Network
* Routability between Networks

## Motivation

As of Today, there is no core object that represents the Network that the Pod attaches to, especially in a multihomed (with multiple network interfaces) Pods. Multi-Network projects main goal is to standardize the way we are able to identify specific network interfaces in a Pod. We want to build on that capability and develop future features that require such specific selection of the Pod's network interfaces.

Lastly, this effort does not try to standardize the configuration of Pod's network interface or how it is attached to the Pod. This part is fully owned by the Multi-Network providers.

### Serviceability

One such feature is serviceability which Today is achieved via the Kubernetes Service object. We use labels and namespace placement as the means to select the correct set of Pods. When Pod has multiple network interfaces, we would like to have the ability to select specific interface of the Pod for the Service. Considering the overwhelming set of capabilities the Service object already holds, we will be looking at the [Gateway API](https://gateway-api.sigs.k8s.io/) as the API to drive this.

### Network Policies

Another feature where we would like to leverage Multi-Network is Network Policies. Here, for a multihomed Pod, we would want to have the ability to apply the policy to the whole Pod or just specific network interfaces.

### Default PodNetwork

In the long term, with this project, we want to be able to create Pods that still have just a single network interface, but it is NOT the Default PodNetwork. This might require targeted changes in core Kubernetes that would allow such functionality, but without pushing the notion of Multi-Network itself there.

## User Stories

Below is the list of use cases that Multi-Network is trying to address. We do not list basic use cases that just add network interfaces to a pod, since those are currently handled by Dynamic Resource Allocation.

### Story \#1

I have implemented my Kubernetes cluster networking using a virtual switch. In this implementation I am capable of creating isolated Networks. I need a means to express to which Network my workloads connect to.

<p align="center">
  <img src="images/vswitch-mn.png?raw=true"/>
</p>

### Story \#2

As a Virtual Machine \-based compute platform provider that I run on top of Kubernetes and Kubevirt I require multi-tenancy with ability to connect my VMs to different VPC-like networks. The isolation has to be achieved on Layer-2 for security reasons.

<p align="center">
  <img src="images/vm-vpc.png?raw=true"/>
</p>


### Story \#3

As a Kubernetes cluster administrator I wish to isolate workloads based on namespaces and network access via assigning different default Network to a Namespace. I do not want the tenants to change their manifests for that purpose. Those workloads should have the same level of Kubernetes functionality: Services, NetworkPolicies, access to Kubernetes API.

<p align="center">
  <img src="images/pod-vpc.png?raw=true"/>
</p>


## Requirements

1. Multi-Network will not require Multi-Network-specific changes in core Kubernetes.
1. Multi-Network should not change the behavior of existing clusters where no multi-network constructs are configured/used.
1. Multi-Network must have the ability to group and identify network interfaces configured in Pods.
1. Multi-Network must provide ability for other features (e.g. serviceability, firewall) to specify which Pod’s network interfaces to target/select.
1. Multi-Network shall provide ability to create Pods without Default PodNetwork (which might need separate effort in core Kubernetes to achieve it).
