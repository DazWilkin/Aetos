ARG GOLANG_VERSION=1.21

ARG GOOS=linux
ARG GOARCH=amd64

ARG COMMIT
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

ARG VERSION
ARG COMMIT

RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} \
    go build \
    -ldflags "-X main.OSVersion=${VERSION} -X main.GitCommit=${COMMIT}" \
    -a -installsuffix cgo \
    -o /go/bin/aetos \
    ./cmd

FROM gcr.io/distroless/static-debian11:latest

LABEL org.opencontainers.image.description "Prometheus Exporter of random metrics|labels"
LABEL org.opencontainers.image.source https://github.com/DazWilkin/aetos

COPY --from=build /go/bin/aetos /

ENTRYPOINT ["/aetos"]
CMD ["--cardinality=3","--labels=3","--metrics=3","--path=metrics"]