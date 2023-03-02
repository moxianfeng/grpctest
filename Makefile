IMG=hub.expvent.com.cn:1111/public/grpctest
DOCKER_VERSION=v0.3

.PHONY: proto
proto:
	buf generate

proto-update:
	buf mod update

build:
	go build -o bin/grpctest

docker:
	docker -D buildx build -f Dockerfile --pull -t $(IMG):$(DOCKER_VERSION) --sbom=false --provenance=false --platform=linux/arm64,linux/amd64 . --push
