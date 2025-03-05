ARG GOLANG_VERSION=1.24

ARG TARGETOS
ARG TARGETARCH

ARG COMMIT
ARG VERSION

FROM --platform=${TARGETARCH} docker.io/golang:${GOLANG_VERSION} AS build

WORKDIR /aetos

COPY go.* ./

COPY api/ ./api
COPY cmd/ ./cmd
COPY collector ./collector
COPY handler ./handler
COPY opts ./opts

ARG TARGETOS
ARG TARGETARCH

ARG COMMIT
ARG VERSION

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -ldflags "-X main.GitCommit=${COMMIT} -X main.OSVersion=${VERSION}" \
    -a -installsuffix cgo \
    -o /go/bin/aetos \
    ./cmd

FROM --platform=${TARGETARCH} gcr.io/distroless/static-debian12:latest

ARG COMMIT
ARG VERSION

LABEL org.opencontainers.image.description="Prometheus Exporter of random metrics|labels"
LABEL org.opencontainers.image.source=https://github.com/DazWilkin/aetos
LABEL org.opencontainers.image.commit=${COMMIT}
LABEL org.opencontainers.image.version=${VERSION}

COPY --from=build /go/bin/aetos /

ENTRYPOINT ["/aetos"]
CMD ["--cardinality=3","--labels=3","--metrics=3","--path=metrics"]
