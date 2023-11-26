package server

import (
	"log"
	"net"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
)

type config struct {
	name    string
	version string
}

type driver struct {
	config config
}

func Start() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{}

	d := &driver{
		config: config{
			name:    "quickstart-sample",
			version: "v0.0.1",
		},
	}

	s := grpc.NewServer(opts...)
	csi.RegisterIdentityServer(s, d)
	csi.RegisterControllerServer(s, d)
	csi.RegisterNodeServer(s, d)
	s.Serve(lis)
}
