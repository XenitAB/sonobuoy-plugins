IMAGE = ghcr.io/xenitab/sonobuoy-plugins/xks:latest

all: test

test:
	go test -timeout=0 -v ./e2e -ginkgo.v

docker-build:
	docker build -t ${IMAGE} .

docker-push:
	docker push ${IMAGE}
