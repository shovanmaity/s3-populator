IMAGE_NAME?=quay.io/shovanmaity/s3-populator:latest

build:
	CGO_ENABLED=0 go build -o s3-populator .

image: build
	docker build -f Dockerfile -t $(IMAGE_NAME) .

push:
	docker push $(IMAGE_NAME)

.PHONY: build image push
