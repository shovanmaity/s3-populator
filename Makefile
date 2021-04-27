IMAGE_NAME?=quay.io/shovanmaity/s3-populator:latest

build:
	CGO_ENABLED=0 go build -o s3-populator .

image:
	docker build -f Dockerfile -t $(IMAGE_NAME) .
	docker image prune -f --filter label=type=build-container

push:
	docker push $(IMAGE_NAME)

.PHONY: build image push
