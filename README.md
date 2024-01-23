[Aëtos](https://en.wikipedia.org/wiki/A%C3%ABtos)

[![build](https://github.com/DazWilkin/Aetos/actions/workflows/build.yml/badge.svg)](https://github.com/DazWilkin/Aetos/actions/workflows/build.yml)

Named after the eagle that tested Prometheus.

An exporter intended for testing Prometheus that can be configured dynamically (by `POST`'ing to its `/publish` endpoint see [publish](#publish)) to generate metrics with labels with random values

## Proto

```bash
MODULE="github.com/DazWilkin/Aetos"
protoc \
--proto_path=${PWD}/protos \
--go_out=${PWD} \
--go_opt=module=${MODULE} \
${PWD}/protos/aetos.proto
```

## Run

```bash
go run github.com/DazWilkin/Aetos/cmd \
--cardinality=3 \
--endpoint=:8080 \
--labels=3 \
--metrics=3 \
--path=/metrics
```

## API

### Publish (`/publish`)

The following will generate 3 metrics each with 3 labels:

```
foo{a="...",b="...",c="..."}
bar{a="...",b="...",c="..."}
baz{a="...",b="...",c="..."}
```

```bash
DATA='
{
    "labels":[
        "a",
        "b",
        "c"
    ],
    "metrics":[
        "foo",
        "bar",
        "baz"
    ]
}
'

curl \
--request POST \
--data "${DATA}" \
http://localhost:8080/publish
```

### Healthz (`/healthz`)

Static | Not implemented

```bash
curl \
--silent \
--request GET \
http://localhost:8080/healthz \
| jq -r .
```
```console
ok
```


### Varz (`/varz`)

```bash
curl \
--silent \
--request GET \
http://localhost:8080/varz \
| jq -r .
```
```JSON
{
  "labels": [
    "a",
    "b",
    "c"
  ],
  "metrics": [
    "foo",
    "bar",
    "baz"
  ]
}
```

## Metrics

Metrics will be empty until a configuration is `/publish`'ed (see [Publish](#publish))

Then:

```console
# HELP aetos_collector_bar a randomly generated metric with 3 labels with cardinality 3
# TYPE aetos_collector_bar gauge
aetos_collector_bar{a="a-0",b="b-0",c="c-0"} 0.4370108094286117
aetos_collector_bar{a="a-1",b="b-1",c="c-1"} 0.164219079866473
aetos_collector_bar{a="a-2",b="b-2",c="c-2"} 0.14382614145278458
# HELP aetos_collector_baz a randomly generated metric with 3 labels with cardinality 3
# TYPE aetos_collector_baz gauge
aetos_collector_baz{a="a-0",b="b-0",c="c-0"} 0.5501234459761718
aetos_collector_baz{a="a-1",b="b-1",c="c-1"} 0.5444040787463975
aetos_collector_baz{a="a-2",b="b-2",c="c-2"} 0.7750368437739439
# HELP aetos_collector_foo a randomly generated metric with 3 labels with cardinality 3
# TYPE aetos_collector_foo gauge
aetos_collector_foo{a="a-0",b="b-0",c="c-0"} 0.87781623201731
aetos_collector_foo{a="a-1",b="b-1",c="c-1"} 0.9076335347948783
aetos_collector_foo{a="a-2",b="b-2",c="c-2"} 0.1317955488606491
```

## Kubernetes

### Jsonnet

```bash
IMAGE="ghcr.io/dazwilkin/aetos:07847ed9eb9fd22d1b50d1ba5f583359c6ced7b6"
PORT="8080"

CARDINALITY="3"
LABELS="3"
METRICS="3"

NAMESPACE="aetos"

jsonnet \
  ./kubernetes.jsonnet \
  --ext-str image="${IMAGE}" \
  --ext-str port="${PORT}" \
  --ext-str cardinality="${CARDINALITY}" \
  --ext-str labels="${LABELS}" \
  --ext-str metrics="${METRICS}" \
| kubectl apply \
  --filename=- \
  --namespace=${NAMESPACE}
```

## Prometheus

### Operator

A `ServiceMonitor` named `aetos` is included with the [`kubernetes.jsonnet`](./kubernetes.jsonnet) file.

### Local

```bash
PORT="9090"

podman run \
--interactive --tty --rm \
--net=host \
--volume=${PWD}/prometheus.yml:/etc/prometheus/prometheus.yml \
prom/prometheus \
--web.config.file=/etc/prometheus/prometheus.yml \
--web.listen-address="0.0.0.0:${PORT}
```

## Sigstore

`aetos` container images are signed by [Sigstore](https://www.sigstore.dev/) and may be verified:

```bash
cosign verify \
--key=./cosign.pub \
ghcr.io/dazwilkin/aetos:07847ed9eb9fd22d1b50d1ba5f583359c6ced7b6
```

> **NOTE** `cosign.pub` may be downloaded [here](./cosign.pub)

To install `cosign`, e.g.:

```bash
go install github.com/sigstore/cosign/cmd/cosign@latest
```


<hr/>
<br/>
<a href="https://www.buymeacoffee.com/dazwilkin" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/default-orange.png" alt="Buy Me A Coffee" height="41" width="174"></a>