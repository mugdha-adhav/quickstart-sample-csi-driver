package server

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

func (d *driver) GetPluginInfo(context.Context, *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	return &csi.GetPluginInfoResponse{
		Name:          d.config.name,
		VendorVersion: d.config.version,
	}, nil
}

func (d *driver) GetPluginCapabilities(context.Context, *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	return &csi.GetPluginCapabilitiesResponse{
		Capabilities: []*csi.PluginCapability{
			// {
			// 	Type: &csi.PluginCapability_Service_{
			// 		Service: &csi.PluginCapability_Service{
			// 			Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
			// 		},
			// 	},
			// },
		},
	}, nil
}

func (d *driver) Probe(context.Context, *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	return &csi.ProbeResponse{}, nil
}
