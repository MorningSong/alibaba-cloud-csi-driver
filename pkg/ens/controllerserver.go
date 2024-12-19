package ens

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-sigs/alibaba-cloud-csi-driver/pkg/common"
	"github.com/kubernetes-sigs/alibaba-cloud-csi-driver/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
)

// controller server try to create/delete volumes/snapshots
type controllerServer struct {
	recorder record.EventRecorder

	createdVolumeMap map[string]*csi.Volume
	common.GenericControllerServer
}

func NewControllerServer() csi.ControllerServer {

	c := &controllerServer{
		recorder:         utils.NewEventRecorder(),
		createdVolumeMap: map[string]*csi.Volume{},
	}

	return c
}

func (cs *controllerServer) ControllerGetCapabilities(context.Context, *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: common.ControllerRPCCapabilities(
			csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
			csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
		),
	}, nil
}

func (cs *controllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {

	klog.Infof("CreateVolume: Starting CreateVolume: %+v", req)

	if value, ok := cs.createdVolumeMap[req.Name]; ok {
		klog.Infof("CreateVolume: volume already be created pvName: %s, VolumeId: %s, volumeContext: %v", req.Name, value.VolumeId, value.VolumeContext)
		return &csi.CreateVolumeResponse{Volume: value}, nil
	}
	reqParams := req.GetParameters()
	diskParams, err := ValidateCreateVolumeParams(reqParams)
	if err != nil {
		klog.Errorf("CreateVolume: invalidate params err: %+v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// every ens node is a independent region, we should get target region by access ens server
	nodeID := diskParams.NodeSelected
	if nodeID != "" {
		instance, err := GlobalConfigVar.ENSCli.DescribeInstance(nodeID)
		if err != nil {
			err = fmt.Errorf("CreateVolume: failed to get region from target nodeID: %+v, err: %+v", nodeID, err)
			klog.Error(err.Error())
			return nil, status.Error(codes.Aborted, err.Error())
		}
		if len(instance) == 1 {
			klog.Infof("CreateVolume: success to get region: %s from target node: %s", *instance[0].EnsRegionId, nodeID)
			diskParams.RegionID = *instance[0].EnsRegionId
		} else {
			err = fmt.Errorf("CreateVolume: get wrong number of instances with instanceid: %s, resp: %+v", nodeID, instance)
			klog.Error(err)
			return nil, status.Error(codes.Aborted, err.Error())
		}
	}

	volSizeBytes := int64(req.GetCapacityRange().GetRequiredBytes())
	requestGB := strconv.FormatInt((volSizeBytes+1024*1024*1024-1)/(1024*1024*1024), 10)
	actualDiskType := ""
	actualDiskID := ""
	for _, dt := range strings.Split(diskParams.DiskType, ",") {
		diskID, err := GlobalConfigVar.ENSCli.CreateVolume(diskParams.RegionID, dt, requestGB)
		if err == nil {
			actualDiskType = dt
			actualDiskID = diskID
			break
		} else {
			klog.Errorf("CreateVolume: failed to create volume: %+v", err)
			if strings.Contains(err.Error(), DiskNotAvailable) || strings.Contains(err.Error(), DiskNotAvailableVer2) {
				klog.Infof("CreateVolume: Create Disk for volume %s with diskCatalog: %s is not supported in region: %s", req.Name, dt, diskParams.RegionID)
				continue
			}
		}
	}
	if actualDiskType == "" || actualDiskID == "" {
		klog.Errorf("CreateVolume: no disk created")
		return nil, status.Errorf(codes.InvalidArgument, "no disk created by regionID: %s, diskType: %s, requestGB: %s", diskParams.RegionID, diskParams.DiskType, requestGB)
	}
	volumeContext := updateVolumeContext(req.GetParameters())
	volumeContext["type"] = actualDiskType
	volumeContext["region"] = diskParams.RegionID

	tmpVol := formatVolumeCreated(actualDiskType, actualDiskID, volSizeBytes, volumeContext)
	cs.createdVolumeMap[req.Name] = tmpVol

	return &csi.CreateVolumeResponse{Volume: tmpVol}, nil
}

func (cs *controllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	return &csi.DeleteVolumeResponse{}, nil
}

func (cs *controllerServer) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {

	klog.Infof("ControllerPublishVolume: Starting Publish Volume: %+v", req)
	if GlobalConfigVar.EnableAttachDetachController == "false" {
		klog.Infof("ControllerPublishVolume: ADController Disable to attach disk: %s to node: %s", req.VolumeId, req.NodeId)
		return &csi.ControllerPublishVolumeResponse{}, nil
	}
	_, err := attachDisk(req.VolumeId, req.NodeId)
	if err != nil {
		klog.Errorf("ControllerPublishVolume: attach disk: %s to node: %s with error: %s", req.VolumeId, req.NodeId, err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	klog.Infof("ControllerPublishVolume: Successful attach disk: %s to node: %s", req.VolumeId, req.NodeId)
	return &csi.ControllerPublishVolumeResponse{}, nil
}

func (cs *controllerServer) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {

	klog.Infof("ControllerUnpublishVolume: Starting Unpublish Volume: %+v", req)

	if GlobalConfigVar.EnableAttachDetachController == "false" {
		klog.Infof("ControllerPublishVolume: ADController Disable to attach disk: %s to node: %s", req.VolumeId, req.NodeId)
		return &csi.ControllerUnpublishVolumeResponse{}, nil
	}

	klog.Infof("ControllerUnpublishVolume: detach disk: %s from node: %s", req.VolumeId, req.NodeId)
	err := detachDisk(req.VolumeId, req.NodeId)
	if err != nil {
		klog.Errorf("ControllerUnpublishVolume: detach disk: %s from node: %s with error: %s", req.VolumeId, req.NodeId, err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	klog.Infof("ControllerUnpublishVolume: Successful detach disk: %s from node: %s", req.VolumeId, req.NodeId)
	return &csi.ControllerUnpublishVolumeResponse{}, nil
}

func (cs *controllerServer) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	return &csi.ControllerExpandVolumeResponse{}, nil
}

func (cs *controllerServer) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	for _, cap := range req.VolumeCapabilities {
		if cap.GetAccessMode().GetMode() != csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER {
			return &csi.ValidateVolumeCapabilitiesResponse{Message: ""}, nil
		}
	}
	return &csi.ValidateVolumeCapabilitiesResponse{
		Confirmed: &csi.ValidateVolumeCapabilitiesResponse_Confirmed{
			VolumeCapabilities: req.VolumeCapabilities,
		},
	}, nil
}
