package server

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/utils/mount"
)

const (
	baseVolumeDir string = "/data/volumes"
)

func (d *driver) NodeGetCapabilities(context.Context, *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: []*csi.NodeServiceCapability{},
	}, nil
}

func (d *driver) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	targetPath := req.GetTargetPath()
	log.Printf("NodePublishVolume: request received for volume %s with target-path: %s\n", req.VolumeId, targetPath)
	hostPath := fmt.Sprintf("%s/%s", baseVolumeDir, req.VolumeId)

	// Create directory on the host
	if err := os.MkdirAll(hostPath, 0777); err != nil {
		return nil, fmt.Errorf("error creating host path: %w", err)
	}
	log.Printf("NodePublishVolume: created host-path directory: %s \n", hostPath)

	// Check if targetPath exists and is a mount point
	mounter := mount.New("")
	notMnt, err := mount.IsNotMountPoint(mounter, targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(targetPath, 0750); err != nil {
				return nil, fmt.Errorf("error creating target path: %w", err)
			}
			notMnt = true
		} else {
			return nil, fmt.Errorf("error while checking if target path is a mountpoint: %w", err)
		}
	}

	if !notMnt {
		return nil, status.Error(codes.Aborted, "target-path is not a mountpoint")
	}

	// Mount hostPath on the targetPath
	options := []string{"bind"}
	if err := mounter.Mount(hostPath, targetPath, "", options); err != nil {
		return nil, fmt.Errorf("failed to mount device: %s at %s", hostPath, targetPath)
	}
	log.Printf("NodePublishVolume: host path %s successfully mounted at %s for volume: %s", hostPath, targetPath, req.VolumeId)

	return &csi.NodePublishVolumeResponse{}, nil
}

func (d *driver) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	targetPath := req.GetTargetPath()
	hostPath := fmt.Sprintf("%s/%s", baseVolumeDir, req.VolumeId)
	log.Printf("NodeUnpublishVolume: request received for volume %s with target-path %s\n", req.VolumeId, targetPath)

	// Unmount only if the target path exists and is really a mount point.
	if notMountPoint, err := mount.IsNotMountPoint(mount.New(""), targetPath); err != nil {
		if !os.IsNotExist(err) {
			return nil, status.Error(codes.NotFound, "couldn't find target-path")
		}
	} else if !notMountPoint {
		// Unmounting the volume/directory.
		err = mount.New("").Unmount(targetPath)
		if err != nil {
			return nil, fmt.Errorf("unmount target path: %w", err)
		}

		log.Println("NodeUnpublishVolume: unmount successful for target path", targetPath)
	}

	// Delete the mount point.
	if err := os.RemoveAll(targetPath); err != nil {
		return nil, fmt.Errorf("remove target path: %w", err)
	}
	log.Println("NodeUnpublishVolume: attempted deletion of target path", targetPath)

	// Delete the hostpath.
	if err := os.RemoveAll(hostPath); err != nil {
		return nil, fmt.Errorf("remove host path: %w", err)
	}
	log.Println("NodeUnpublishVolume: attempted deletion of host path", hostPath)

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (d *driver) NodeStageVolume(context.Context, *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	return &csi.NodeStageVolumeResponse{}, nil
}

func (d *driver) NodeUnstageVolume(context.Context, *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (d *driver) NodeGetVolumeStats(context.Context, *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	return &csi.NodeGetVolumeStatsResponse{}, nil
}

func (d *driver) NodeExpandVolume(context.Context, *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	return &csi.NodeExpandVolumeResponse{}, nil
}

func (d *driver) NodeGetInfo(context.Context, *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	return &csi.NodeGetInfoResponse{
		NodeId: d.config.nodeID,
	}, nil
}
