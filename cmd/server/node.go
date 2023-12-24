package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/kubernetes/pkg/volume/util/volumepathhandler"
	"k8s.io/utils/mount"
)

const (
	baseVolumeDir string = "/tmp/quickstart/volumes"
)

func init() {
	if err := os.MkdirAll(baseVolumeDir, os.ModeSticky|os.ModePerm); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func (d *driver) NodeGetCapabilities(context.Context, *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: []*csi.NodeServiceCapability{},
	}, nil
}

func (d *driver) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	targetPath := req.GetTargetPath()
	log.Printf("NodePublishVolume: request received for volume %s with target-path: %s\n", req.VolumeId, targetPath)
	driverPath := fmt.Sprintf("%s/%s", baseVolumeDir, req.VolumeId)

	mounter := mount.New("")

	// Create a block file.
	_, err := os.Stat(driverPath)
	if err != nil {
		if os.IsNotExist(err) {
			cmd := exec.Command("fallocate", "-l", "100M", driverPath)
			output, err := cmd.CombinedOutput()
			if err != nil {
				return nil, fmt.Errorf("failed to create block device: %v, %v", err, string(output))
			}
		} else {
			return nil, fmt.Errorf("failed to stat block device: %v, %v", driverPath, err)
		}
	}

	// Associate block file with the loop device.
	volPathHandler := volumepathhandler.VolumePathHandler{}
	_, err = volPathHandler.AttachFileDevice(driverPath)
	if err != nil {
		// Remove the block file because it'll no longer be used again.
		if err2 := os.Remove(driverPath); err2 != nil {
			fmt.Printf("failed to cleanup block file %s: %v\n", driverPath, err2)
		}
		return nil, fmt.Errorf("failed to attach device %v: %v", driverPath, err)
	}

	// Get loop device from the volume path.
	loopDevice, err := volPathHandler.GetLoopDevice(driverPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get the loop device: %w", err)
	}

	log.Println("NodePublishVolume: loop device", loopDevice)

	// Check if the target path exists. Create if not present.
	_, err = os.Lstat(targetPath)
	if os.IsNotExist(err) {
		if err = makeFile(targetPath); err != nil {
			return nil, fmt.Errorf("failed to create target path: %w", err)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to check if the target block file exists: %w", err)
	}

	// Check if the target path is already mounted, if yes prevent remounting.
	notMountPoint, err := mount.IsNotMountPoint(mounter, targetPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error checking path %s for mount: %w", targetPath, err)
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
	if err := mounter.Mount(loopDevice, targetPath, "", options); err != nil {
		return nil, fmt.Errorf("failed to mount block device: %s at %s: %w", driverPath, targetPath, err)
	}
	log.Printf("NodePublishVolume: volume succesesfully mounted at: %s for volume: %s", targetPath, req.VolumeId)

	return &csi.NodePublishVolumeResponse{}, nil
}

func (d *driver) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	targetPath := req.GetTargetPath()
	driverPath := fmt.Sprintf("%s/%s", baseVolumeDir, req.VolumeId)

	fileExists := func(path string) {
		_, err := os.Stat(path)

		if errors.Is(err, os.ErrNotExist) {
			log.Println("NodeUnPublishVolume: Path doesn't exist: ", path)
			return
		}
		log.Println("NodeUnPublishVolume: Path exists: ", path)
	}

	fileExists(targetPath)
	fileExists(driverPath)

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

	// Delete the volume
	if err := os.RemoveAll(driverPath); err != nil {
		return nil, fmt.Errorf("remove host path: %w", err)
	}
	log.Printf("NodeUnPublishVolume: host path deleted %s with target-path: %s\n", req.VolumeId, targetPath)

	fileExists(targetPath)
	fileExists(driverPath)
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

// makeFile ensures that the file exists, creating it if necessary.
// The parent directory must exist.
func makeFile(pathname string) error {
	f, err := os.OpenFile(pathname, os.O_CREATE, os.FileMode(0o644))
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	defer f.Close()

	return nil
}
