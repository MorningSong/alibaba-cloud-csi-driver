package pov

import (
	"context"
	"os"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-sigs/alibaba-cloud-csi-driver/pkg/cloud/metadata"
	"github.com/kubernetes-sigs/alibaba-cloud-csi-driver/pkg/common"
	"github.com/kubernetes-sigs/alibaba-cloud-csi-driver/pkg/pov/internal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
	mountutils "k8s.io/mount-utils"
)

type nodeService struct {
	inFlight *internal.InFlight
	mounter  mountutils.Interface
	common.GenericNodeServer
}

func newNodeService(meta metadata.MetadataProvider) nodeService {
	return nodeService{
		inFlight: internal.NewInFlight(),
		mounter:  mountutils.NewWithoutSystemd(""),
		GenericNodeServer: common.GenericNodeServer{
			NodeID: metadata.MustGet(meta, metadata.InstanceID),
		},
	}
}

func (d *nodeService) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	return &csi.NodeStageVolumeResponse{}, nil
}

func (d *nodeService) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (d *nodeService) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	return &csi.NodeExpandVolumeResponse{}, nil
}

func (d *nodeService) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	if !d.inFlight.Insert(req.TargetPath) {
		return nil, status.Errorf(codes.Aborted, "There is already an operation for %s", req.TargetPath)
	}
	defer d.inFlight.Delete(req.TargetPath)

	notMnt, err := d.mounter.IsLikelyNotMountPoint(req.TargetPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(req.TargetPath, os.ModePerm); err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			notMnt = true
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	if !notMnt {
		klog.Infof("NodePublishVolume: %s already mounted", req.TargetPath)
		return &csi.NodePublishVolumeResponse{}, nil
	}

	err = d.mounter.Mount("tmpfs", req.TargetPath, "tmpfs", []string{"ro"})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "mount readonly tmpfs: %v", err)
	}

	klog.V(2).InfoS("NodePublishVolume succeeded", "volumeId", req.VolumeId, "targetPath", req.TargetPath)
	return &csi.NodePublishVolumeResponse{}, nil
}

func (d *nodeService) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	if !d.inFlight.Insert(req.TargetPath) {
		return nil, status.Errorf(codes.Aborted, "There is already an operation for %s", req.TargetPath)
	}
	defer d.inFlight.Delete(req.TargetPath)

	var err error
	forceUnmounter, ok := d.mounter.(mountutils.MounterForceUnmounter)
	if ok {
		err = mountutils.CleanupMountWithForce(req.TargetPath, forceUnmounter, false, time.Minute)
	} else {
		err = mountutils.CleanupMountPoint(req.TargetPath, d.mounter, false)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "umount %s: %v", req.TargetPath, err)
	}

	klog.V(2).InfoS("NodeUnpublishVolume succeeded", "volumeId", req.VolumeId, "targetPath", req.TargetPath)
	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (d *nodeService) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	segments := map[string]string{
		TopologyKey: GlobalConfigVar.regionID,
	}
	topology := &csi.Topology{Segments: segments}

	return &csi.NodeGetInfoResponse{
		NodeId:             d.NodeID,
		MaxVolumesPerNode:  d.getVolumesLimit(),
		AccessibleTopology: topology,
	}, nil
}

func (d *nodeService) getVolumesLimit() int64 {
	return 0
}

func (d *nodeService) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	return &csi.NodeGetCapabilitiesResponse{}, nil
}

func (d *nodeService) NodeGetVolumeStats(ctx context.Context, req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	return &csi.NodeGetVolumeStatsResponse{}, nil
}
