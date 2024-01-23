local host = std.extVar("host");
local port = std.parseInt(std.extVar("port"));

{
    "kind": "Ingress",
    "apiVersion": "networking.k8s.io/v1",
    "metadata": {
        "labels": {
        "app": "aetos"
        },
        "name": "aetos",
    },
    "spec": {
        "defaultBackend": {
        "service": {
            "name": "aetos",
            "port": {
            "number": port
            }
        }
        },
        "ingressClassName": "tailscale",
        "tls": [
        {
            "hosts": [
            host
            ]
        }
        ]
    }
}