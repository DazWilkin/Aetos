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

### Publish

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

### Varz

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
aetos_collector_bar{a="737c191aed52450b9c655083c9971bdd",b="b1d10db2016c2f83c13b25fcb170cdeb",c="298a664f356e310dbaf9117a0d108b1e"} 0.14934113401237586
aetos_collector_bar{a="a165efd196e17ba195ad4dc50028b39a",b="34f25f6f596e0e4a471136e00726093b",c="63e3dc58db6926e5fd33177aa05336f9"} 0.5069931251554308
aetos_collector_bar{a="b39baf03412f39006635c8da36237ff0",b="f474ce9df880f0a1f5d810a7ab7a539d",c="caf9334ca1325a0ff28ec4b7c88aa06e"} 0.6938566189478629
# HELP aetos_collector_baz a randomly generated metric with 3 labels with cardinality 3
# TYPE aetos_collector_baz gauge
aetos_collector_baz{a="737c191aed52450b9c655083c9971bdd",b="b1d10db2016c2f83c13b25fcb170cdeb",c="298a664f356e310dbaf9117a0d108b1e"} 0.32286605378434446
aetos_collector_baz{a="a165efd196e17ba195ad4dc50028b39a",b="34f25f6f596e0e4a471136e00726093b",c="63e3dc58db6926e5fd33177aa05336f9"} 0.8045848769581659
aetos_collector_baz{a="b39baf03412f39006635c8da36237ff0",b="f474ce9df880f0a1f5d810a7ab7a539d",c="caf9334ca1325a0ff28ec4b7c88aa06e"} 0.8665926183342777
# HELP aetos_collector_foo a randomly generated metric with 3 labels with cardinality 3
# TYPE aetos_collector_foo gauge
aetos_collector_foo{a="737c191aed52450b9c655083c9971bdd",b="b1d10db2016c2f83c13b25fcb170cdeb",c="298a664f356e310dbaf9117a0d108b1e"} 0.7990841616320065
aetos_collector_foo{a="a165efd196e17ba195ad4dc50028b39a",b="34f25f6f596e0e4a471136e00726093b",c="63e3dc58db6926e5fd33177aa05336f9"} 0.6949752807660148
aetos_collector_foo{a="b39baf03412f39006635c8da36237ff0",b="f474ce9df880f0a1f5d810a7ab7a539d",c="caf9334ca1325a0ff28ec4b7c88aa06e"} 0.43778715739069657
```

## Kubernetes

### Jsonnet

```bash
IMAGE="ghcr.io/dazwilkin/aetos:168150eca910b4707b75da5302df217bbb43b12e"
PORT="8080"

NAMESPACE="aetos"

jsonnet \
  ./kubernetes.jsonnet \
  --ext-str image=${IMAGE} \
  --ext-str port=${PORT} \
| kubectl apply \
  --filename=- \
  --namespace=${NAMESPACE}
```

## Prometheus

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
ghcr.io/dazwilkin/aetos:168150eca910b4707b75da5302df217bbb43b12e
```

> **NOTE** `cosign.pub` may be downloaded [here](./cosign.pub)

To install `cosign`, e.g.:

```bash
go install github.com/sigstore/cosign/cmd/cosign@latest
```


<hr/>
<br/>
<a href="https://www.buymeacoffee.com/dazwilkin" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/default-orange.png" alt="Buy Me A Coffee" height="41" width="174"></a>