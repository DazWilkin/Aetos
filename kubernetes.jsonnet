local image = std.extVar("image");
local port = std.parseInt(std.extVar("port"));

local cardinality = "--cardinality=" + std.extVar("cardinality");
local endpoint = "--endpoint=0.0.0.0:" + std.extVar("port");
local labels = "--labels=" + std.extVar("labels");
local metrics = "--metrics=" + std.extVar("metrics");

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
                    cardinality,
                    endpoint,
                    labels,
                    metrics
                  ],
                  "image": image,
                  "imagePullPolicy": "IfNotPresent",
                  "name": "aetos",
                  "ports": [
                    {
                      "name": "metrics",
                      "containerPort": port,
                      "protocol": "TCP"
                    }
                  ],
                  "resources": {
                    "limits": {
                      "memory": "500Mi"
                    },
                    "requests":{
                      "cpu": "250m",
                      "memory": "250Mi"
                    }
                  },
                  "securityContext": {
                    "allowPrivilegeEscalation": false,
                    "privileged": false,
                    "readOnlyRootFilesystem": true,
                    "runAsGroup": 1000,
                    "runAsNonRoot": true,
                    "runAsUser": 1000,
                  }
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
              "port": port,
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
      },
      {
        "kind": "VerticalPodAutoscaler",
        "apiVersion": "autoscaling.k8s.io/v1",
        "metadata": {
          "name": "aetos",
          "labels": {
            "app": "aetos",
          }
        },
        "spec": {
          "targetRef": {
            "kind": "Deployment",
            "apiVersion": "apps/v1",
            "name": "aetos"
          },
          "updatePolicy": {
            "updateMode": "Off"
          }
        }
      }
    ]
  }