IMAGE = ghcr.io/xenitab/sonobuoy-plugins/xks:latest
E2E_SKIP =
E2E_FOCUS =
E2E_EXTRA_ARGS =

all: test

test:
	go test -timeout=0 -v ./e2e -ginkgo.v -ginkgo.skip="${E2E_SKIP}" -ginkgo.focus="${E2E_FOCUS}" ${E2E_EXTRA_ARGS}

docker-build:
	docker build -t ${IMAGE} .

docker-push:
	docker push ${IMAGE}
