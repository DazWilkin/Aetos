#!/usr/bin/env bash

IMAGE="ghcr.io/dazwilkin/aetos:a941df2172122d8840866a126345a2e2bbf23489"
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