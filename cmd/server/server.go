package server

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
)

type config struct {
	name     string
	version  string
	nodeID   string
	endpoint string
}

type driver struct {
	config config
}

func Start(endpoint, name, nodeID string) {
	d := &driver{
		config: config{
			name:     name,
			version:  "v0.0.1",
			endpoint: endpoint,
			nodeID:   nodeID,
		},
	}

	log.Printf("\nDriver config is: %v", d)

	proto, addr := parseEndpoint(d.config.endpoint)

	lis, err := net.Listen(proto, addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer([]grpc.ServerOption{}...)
	csi.RegisterIdentityServer(s, d)
	// csi.RegisterControllerServer(s, d)
	csi.RegisterNodeServer(s, d)
	s.Serve(lis)
}

func parseEndpoint(endpoint string) (string, string) {
	s := strings.SplitN(endpoint, "://", 2)
	os.Remove(s[1])
	return s[0], s[1]
}
