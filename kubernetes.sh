#!/usr/bin/env bash

IMAGE="ghcr.io/dazwilkin/aetos:3629797944b7c178749ffbaf47087d81070a64ef"
PORT="8080"

# Revise
CARDINALITY="3" # Max  10
LABELS="3"      # Max   5
METRICS="3"     # Max 250

NAMESPACE="aetos"

kubectl create namespace ${NAMESPACE}

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

kubectl get deployment/aetos \
--namespace=${NAMESPACE} \
--output=jsonpath="{.spec.template.spec.containers[0].image}"
