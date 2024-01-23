ARG GOLANG_VERSION=1.21.6

ARG GOOS=linux
ARG GOARCH=amd64

ARG CHECKSUM
ARG VERSION

FROM docker.io/golang:${GOLANG_VERSION} as build

WORKDIR /aetos

COPY go.* ./

COPY api/ ./api
COPY cmd/ ./cmd
COPY collector ./collector
COPY handler ./handler
COPY xxx ./xxx

ARG GOOS
ARG GOARCH

ARG CHECKSUM
ARG VERSION

RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} \
    go build \
    -ldflags "-X main.Checksum=${CHECKSUM} -X main.Version=${VERSION}" \
    -a -installsuffix cgo \
    -o /go/bin/aetos \
    ./cmd

FROM gcr.io/distroless/static-debian11:latest

LABEL org.opencontainers.image.description "Prometheus Exporter of random metrics|labels"
LABEL org.opencontainers.image.source https://github.com/DazWilkin/aetos

COPY --from=build /go/bin/aetos /

ENTRYPOINT ["/aetos"]
CMD ["--cardinality=3","--labels=3","--metrics=3","--path=metrics"]