FROM golang:1.15-alpine as builder

RUN apk add gcc pkgconfig libc-dev

WORKDIR /workspace

# copy modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# cache modules
RUN go mod download

# copy source code
COPY e2e/ e2e/

RUN CGO_ENABLED=1 go test -c -o e2e.test ./e2e

FROM alpine:3.12

# link repo to the GitHub Container Registry image
LABEL org.opencontainers.image.source="https://github.com/xenitab/sonobuoy-plugins"

RUN apk add --no-cache ca-certificates tini

COPY --from=builder /workspace/e2e.test /usr/local/bin/

RUN addgroup -S e2e && adduser -S -g e2e e2e

USER e2e

ENTRYPOINT [ "/sbin/tini", "--", "e2e.test" ]
