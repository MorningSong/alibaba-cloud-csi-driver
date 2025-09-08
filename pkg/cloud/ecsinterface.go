package cloud

import "github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

type ECSInterface interface {
	CreateDisk(request *ecs.CreateDiskRequest) (response *ecs.CreateDiskResponse, err error)
	DescribeAvailableResource(request *ecs.DescribeAvailableResourceRequest) (response *ecs.DescribeAvailableResourceResponse, err error)
	DescribeInstanceHistoryEvents(request *ecs.DescribeInstanceHistoryEventsRequest) (response *ecs.DescribeInstanceHistoryEventsResponse, err error)
	AttachDisk(request *ecs.AttachDiskRequest) (response *ecs.AttachDiskResponse, err error)
	DetachDisk(request *ecs.DetachDiskRequest) (response *ecs.DetachDiskResponse, err error)
	DeleteDisk(request *ecs.DeleteDiskRequest) (response *ecs.DeleteDiskResponse, err error)
	DescribeInstances(request *ecs.DescribeInstancesRequest) (response *ecs.DescribeInstancesResponse, err error)
	DescribeInstanceTypes(request *ecs.DescribeInstanceTypesRequest) (response *ecs.DescribeInstanceTypesResponse, err error)
	DescribeDisks(request *ecs.DescribeDisksRequest) (response *ecs.DescribeDisksResponse, err error)
	ResizeDisk(request *ecs.ResizeDiskRequest) (response *ecs.ResizeDiskResponse, err error)
	CreateSnapshot(request *ecs.CreateSnapshotRequest) (response *ecs.CreateSnapshotResponse, err error)
	DescribeSnapshots(request *ecs.DescribeSnapshotsRequest) (response *ecs.DescribeSnapshotsResponse, err error)
}
