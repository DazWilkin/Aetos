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

