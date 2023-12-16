package server

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/utils/mount"
)

func (d *driver) NodeStageVolume(context.Context, *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	return &csi.NodeStageVolumeResponse{}, nil
}
func (d *driver) NodeUnstageVolume(context.Context, *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	return &csi.NodeUnstageVolumeResponse{}, nil
}
func (d *driver) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	log.Printf("NodePublishVolume: request received for volume %s with target-path: %s", req.VolumeId, req.TargetPath)
	path := fmt.Sprintf("%s/%s", baseVolumeDir, req.VolumeId)

	// Check if the target path is already mounted, if yes prevent remounting.
	notMountPoint, err := mount.IsNotMountPoint(mount.New(""), req.TargetPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error checking path %s for mount: %w", req.TargetPath, err)
		}
		// Target Path is ready for mounting.
		notMountPoint = true
	}
	if !notMountPoint {
		// Target path is already mounted, return gracefully.
		return &csi.NodePublishVolumeResponse{}, nil
	}

	// Mounting the volume.
	options := []string{"bind"}
	if err := mount.New("").Mount(path, req.TargetPath, "", options); err != nil {
		return nil, fmt.Errorf("failed to mount block device: %s at %s: %w", path, req.TargetPath, err)
	}
	log.Printf("NodePublishVolume: volume succesesfully mounted at: %s for volume: %s", req.TargetPath, req.VolumeId)

	return &csi.NodePublishVolumeResponse{}, nil
}
func (d *driver) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	targetPath := req.GetTargetPath()

	// Unmount only if the target path is really a mount point.
	if notMountPoint, err := mount.IsNotMountPoint(mount.New(""), targetPath); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("check target path: %w", err)
		}
	} else if !notMountPoint {
		// Unmounting the image or filesystem.
		err = mount.New("").Unmount(targetPath)
		if err != nil {
			return nil, fmt.Errorf("unmount target path: %w", err)
		}
	}
	// Delete the mount point.
	if err := os.RemoveAll(targetPath); err != nil {
		return nil, fmt.Errorf("remove target path: %w", err)
	}
	return &csi.NodeUnpublishVolumeResponse{}, nil
}
func (d *driver) NodeGetVolumeStats(context.Context, *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	return &csi.NodeGetVolumeStatsResponse{}, nil
}
func (d *driver) NodeExpandVolume(context.Context, *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	return &csi.NodeExpandVolumeResponse{}, nil
}
func (d *driver) NodeGetCapabilities(context.Context, *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: []*csi.NodeServiceCapability{},
	}, nil
}
func (d *driver) NodeGetInfo(context.Context, *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	return &csi.NodeGetInfoResponse{
		NodeId: d.config.nodeID,
	}, nil
}
