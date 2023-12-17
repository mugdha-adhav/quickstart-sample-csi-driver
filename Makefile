IMAGE=mugdhaadhav/quickstart-sample-csi-driver
TAG=v0.0.1
build:
	docker build -f Dockerfile -t ${IMAGE}:${TAG} .

kind-push:
	kind load docker-image ${IMAGE}:${TAG}

run: 
	docker run -it -p 50051:50051 ${IMAGE}:${TAG}
