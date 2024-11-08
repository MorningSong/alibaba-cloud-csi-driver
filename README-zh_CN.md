# 阿里云Kubernetes CSI插件
[![GoReportCard Widget]][GoReportCardResult]

[English](./README.md) | 简体中文

## 插件介绍

阿里云CSI插件实现了在Kubernetes中对阿里云云存储卷的生命周期管理，支持动态创建、挂载、使用云数据卷。 当前的CSI实现基于K8S 1.14以上的版本；

支持的阿里云存储：***云盘、NAS、CPFS、OSS、LVM***

警告：不建议在 ACK 集群中手动部署该驱动程序。用户应使用 ACK 组件中心自动部署和管理 [Alibaba Cloud CSI Driver](https://help.aliyun.com/zh/ack/product-overview/csi-plugin)。 flexvolume 迁移场景除外, 请按照[迁移文档](https://help.aliyun.com/zh/ack/ack-managed-and-ack-dedicated/user-guide/use-csi-compatible-controller-to-migrate-from-flexvolume-to-csi)进行迁移。

免责声明：ACK 官方不支持在集群中手动部署该驱动程序。


## 版本说明

| Feature         | Stage | Min Kubernetes Version | Min Driver Version |
|-----------------|-------|------------------------|--------------------|
| Topology        | GA    | 1.17                   | v1.0.2             |
| Resize (Expand) | GA    | 1.16                   | v1.0.5             |
| Snapshots       | GA    | 1.20                   | v1.1.2             |


### 云盘CSI插件

云盘CSI插件支持动态创建云盘数据卷、挂载数据卷。云盘是一种块存储类型，只能同时被一个负载使用(ReadWriteOnce)。

云盘CSI插件更多详细说明请参考[云盘](./docs/disk.md)


### NAS CSI插件

NAS CSI插件支持为应用负载挂载阿里云NAS存储卷，也支持动态创建NAS卷。NAS存储是一种共享存储，可以同时被多个应用负载使用(ReadWriteMany)。

NAS CSI插件更多详细说明请参考[NAS](./docs/nas.md)


### CPFS CSI插件

**已删除：挂载 CPFS 2.0 请使用 NAS CSI 插件**

### OSS CSI插件

OSS CSI插件支持为应用负载挂载阿里云OSS Bucket，目前不支持动态创建OSS Bucket。OSS存储是一种共享存储，可以同时被多个应用负载使用(ReadWriteMany)。

OSS CSI插件更多详细说明请参考[OSS](./docs/oss.md)

## 社区, 贡献, 讨论, 支持

可以到 [Kubernetes](https://kubernetes.io/community/) 社区学到如何获取支持；

可以到 [Cloud Provider SIG](https://github.com/kubernetes/community/tree/master/sig-cloud-provider) 联系到项目管理者；

可以加入钉钉群(群ID：33936810)和我们一起讨论遇到的问题；


### 行为准则

参与Kubernetes社区参考[Kubernetes行为准则](code-of-conduct.md)；

可以向社区提交 [Issue](https://github.com/kubernetes-sigs/alibaba-cloud-csi-driver/issues)；


[GoReportCard Widget]: https://goreportcard.com/badge/github.com/kubernetes-sigs/alibaba-cloud-csi-driver
[GoReportCardResult]: https://goreportcard.com/report/github.com/kubernetes-sigs/alibaba-cloud-csi-driver


## 安全
对于发现的安全漏洞，请邮件发送至kubernetes-security@service.aliyun.com，您可在[SECURITY.md](./SECURITY.md)文件中找到更多信息。

## 链接

- [安装](./docs/install.md)
