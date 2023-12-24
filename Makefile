IMAGE=mugdhaadhav/quickstart-sample-csi-driver
TAG=latest
build:
	docker build -f Dockerfile -t ${IMAGE}:$(TAG) .

push:
	docker build -f Dockerfile -t ${IMAGE}:$(TAG) . --push

load:
	kind load docker-image ${IMAGE}:$(TAG)

run: 
	docker run -it -p 50051:50051 ${IMAGE}:$(TAG)

deploy-kind: build load
	kubectl apply -f deploy

remove-kind:
	kubectl delete -f deploy
