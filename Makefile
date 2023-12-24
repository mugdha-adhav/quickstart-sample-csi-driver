IMAGE=mugdhaadhav/quickstart-sample-csi-driver
TAG=v0.0.1
CLUSTER_NAME=kind
build:
	docker build -f Dockerfile -t ${IMAGE}:$(TAG) .

push:
	docker build -f Dockerfile -t ${IMAGE}:$(TAG) . --push

load:
	kind load docker-image ${IMAGE}:$(TAG) --name $(CLUSTER_NAME)

run: 
	docker run -it -p 50051:50051 ${IMAGE}:$(TAG)

deploy-kind:
	kubectl apply -f deploy

remove-kind:
	kubectl delete -f deploy
