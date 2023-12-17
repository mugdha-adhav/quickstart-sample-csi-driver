package main

import (
	"flag"

	"github.com/mugdha-adhav/quickstart-sample-csi-driver/cmd/server"
)

func main() {
	var endpoint, name, nodeID string
	flag.StringVar(&endpoint, "endpoint", "unix:///csi/csi.sock", "CSI endpoint")
	flag.StringVar(&name, "drivername", "quickstart.csi.k8s.io", "Name of the driver")
	flag.StringVar(&nodeID, "nodeid", "", "node id")

	flag.Parse()

	server.Start(endpoint, name, nodeID)
}
