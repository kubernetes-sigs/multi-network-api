# Multi-Network API

The Multi-Network API is a sub-project of Kubernetes [SIG Network](https://github.com/kubernetes/community/tree/master/sig-network). This repository contains the API definitions, Custom Resource Definitions (CRDs), specifications, and generated Go client code for Kubernetes Multi-Network API.

Kubernetes historically attaches all pods to a single cluster-wide network. The Multi-Network project standardizes how pods connect to multiple distinct networks and how those networks are identified, categorized, and selected. This repository defines the `PodNetworkKind` API, standard Dynamic Resource Allocation (DRA) attributes, and discovery mechanisms used by pod network implementations.

For more details, see:
- [Requirements and Use Cases](docs/requirements.md)
- [PodNetworkKind Design Specification](docs/pod-network-kind.md)

## Getting Started

To install the `PodNetworkKind` CRD in your cluster:
```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/multi-network-api/refs/heads/main/deployment/multinetwork.networking.x-k8s.io_podnetworkkinds.yaml
```

To install the controller for the CRD, visit https://github.com/kubernetes-sigs/multi-network.

Go package documentation:
- https://pkg.go.dev/sigs.k8s.io/multi-network-api
- https://pkg.go.dev/sigs.k8s.io/multi-network

## Community, Discussion, and Support

This project is maintained by Kubernetes [SIG Network](https://github.com/kubernetes/community/tree/master/sig-network) under the Multi-Network sub-project.

- **Controller Implementation:** [kubernetes-sigs/multi-network](https://github.com/kubernetes-sigs/multi-network)
- **Meeting Notes:** [Multi-Network Meeting Notes](https://docs.google.com/document/d/1pe_0aOsI35BEsQJ-FhFH9Z_pWQcU2uqwAnOx2NIx6OY/edit?usp=sharing)
- **Bi-weekly Meeting (Tuesdays at 06:00 PT (Pacific Time)):** [Zoom Link](https://zoom.us/j/95680858961?pwd=M1c2TTdMZHpMUUtIYXRpbjRobkNJZz09)
- **Meeting Recordings:** [YouTube Kubernetes Community](https://www.youtube.com/@KubernetesCommunity/search?query=multi-network)
- **Slack:** [#sig-network-multi-network](https://kubernetes.slack.com/archives/C03UT5H9KDZ) on [Kubernetes Slack](https://slack.k8s.io/)
- **Mailing List:** [SIG Network Google Group](https://groups.google.com/a/kubernetes.io/g/sig-network)

### Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for instructions on how to contribute to this repository.

### Code of Conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).

### License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
