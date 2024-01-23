local image = std.extVar("image");
local port = std.extVar("port");

local cardinality = std.extVar("cardinality");
local labels = std.extVar("labels");
local metrics = std.extVar("metrics");

{
    "kind": "List",
    "apiVersion": "v1",
    "metadata": null,
    "items": [
      {
        "kind": "Deployment",
        "apiVersion": "apps/v1",
        "metadata": {
          "labels": {
            "app": "aetos"
          },
          "name": "aetos"
        },
        "spec": {
          "selector": {
            "matchLabels": {
              "app": "aetos"
            }
          },
          "template": {
            "metadata": {
              "labels": {
                "app": "aetos"
              }
            },
            "spec": {
              "containers": [
                {
                  "args": [
                    "--cardinality="+ cardinality,
                    "--endpoint=0.0.0.0:" + port,
                    "--labels=" + labels,
                    "--metrics=" + metrics
                  ],
                  "image": image,
                  "imagePullPolicy": "IfNotPresent",
                  "name": "aetos",
                  "ports": [
                    {
                      "name": "metrics",
                      "containerPort": std.parseInt(port),
                      "protocol": "TCP"
                    }
                  ],
                  "resources": {}
                }
              ]
            }
          }
        }
      },
      {
        "kind": "Service",
        "apiVersion": "v1",
        "metadata": {
          "labels": {
            "app": "aetos"
          },
          "name": "aetos"
        },
        "spec": {
          "selector": {
            "app": "aetos"
          },
          "ports": [
            {
              "name": "metrics",
              "port": std.parseInt(port),
              "protocol": "TCP",
              "targetPort": "metrics"
            }
          ],
          "type": "NodePort"
        }
      },
      {
        "kind": "ServiceMonitor",
        "apiVersion": "monitoring.coreos.com/v1",
        "metadata": {
          "labels": {
            "app": "aetos"
          },
          "name": "aetos"
        },
        "spec": {
          "selector": {
            "matchLabels": {
              "app": "aetos"
            }
          },
          "endpoints": [
            {
              "port": "metrics",
            }
          ]
        }
      }
    ]
  }