IMAGE=mugdhaadhav/quickstart-sample-csi-driver
TAG=v0.0.1
build:
	docker build -f Dockerfile -t ${IMAGE}:${TAG} .

run: 
	docker run -it ${IMAGE}:${TAG}
