# Alibaba Cloud Kubernetes CSI Plugin
[![GoReportCard Widget]][GoReportCardResult]

English | [简体中文](./README-zh_CN.md)

WARNING: Deploying this driver to your ACK cluster manually is not recommended. Instead users should use ACK Add-ons to automatically deploy and manage the [Alibaba Cloud CSI Driver](https://www.alibabacloud.com/help/en/ack/product-overview/csi-plugin).

DISCLAIMER: Manual deployment of the driver in your ACK cluster is not officially supported by Alibaba Cloud.


## Introduction
Alibaba Cloud CSI plugins implement an interface between CSI enabled Container
Orchestrator and Alibaba Cloud Storage. It allows dynamically provision Disk
volumes and attach it to workloads.

Current implementation of CSI plugins has been tested in Kubernetes environment (requires Kubernetes 1.14+).

Current Support: ***Cloud Disk, NAS, CPFS, OSS***;

## Features in Development

| Feature         | Stage | Min Kubernetes Version | Min Driver Version |
|-----------------|-------|------------------------|--------------------|
| Topology        | GA    | 1.17                   | v1.0.2             |
| Resize (Expand) | GA    | 1.16                   | v1.0.5             |
| Snapshots       | GA    | 1.20                   | v1.1.2             |


### Cloud Disk CSI Plugin

Disk CSI Plugin support Cloud disk provision and attachment. And Cloud disk is type of block storage, can only used as ReadWriteOnce mode. Only be attached to one node at the same time.

More detail information please refer to [Cloud Disk](./docs/disk.md).


### NAS CSI Plugin

NAS CSI Plugin can support NAS volume provision and mount. Alibaba Cloud Network Attached Storage (NAS) storage is type of network storage which compatible with multiple standard protocols, such as NFS and SMB, and can be mount by multi nodes at the same time.

More detail information please refer to [NAS](./docs/nas.md).


### CPFS CSI Plugin

**Removed: Switch to NAS CSI plugin for CPFS 2.0**

### OSS CSI Plugin

OSS CSI Plugin support OSS bucket mount, but does not support provision volume. OSS storage is type of object storage and can be mount by multi nodes at the same time.

More detail information pls refer to [OSS](./docs/oss.md).


## Community, discussion, contribution, and support

Learn how to engage with the Kubernetes community on the [community page](https://kubernetes.io/community/).

You can reach the maintainers of this project at the [Cloud Provider SIG](https://github.com/kubernetes/community/tree/master/sig-cloud-provider).

You can join the DingDing Talking (GroupID: 33936810) to talk with us.

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).

Please submit an issue at: [Issues](https://github.com/kubernetes-sigs/alibaba-cloud-csi-driver/issues)


[GoReportCard Widget]: https://goreportcard.com/badge/github.com/kubernetes-sigs/alibaba-cloud-csi-driver
[GoReportCardResult]: https://goreportcard.com/report/github.com/kubernetes-sigs/alibaba-cloud-csi-driver

## Security
Please report vulnerabilities by email to kubernetes-security@service.aliyun.com. Also see our [SECURITY.md](./SECURITY.md) file for details.

## Links

- [Installation](./docs/install.md)
