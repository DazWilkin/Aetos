[Aëtos](https://en.wikipedia.org/wiki/A%C3%ABtos)

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

## Publish

```bash
DATA="
{
    \"labels\":[
        \"a\",
        \"b\",
        \"c\"
    ],
    \"metrics\":[
        \"foo\",
        \"bar\",
        \"baz\"
    ]
}
"

curl \
--request POST \
--data "${DATA}" \
http://localhost:8080/publish
```

## Metrics

Metrics will be empty until a configuration is `/publish`'ed (see [Publish](#Publish))

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

