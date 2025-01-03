#!/usr/bin/env bash

IMAGE="ghcr.io/dazwilkin/aetos:cdcce4e943a2b55040aecb40d07e70b5ea5f7f84"
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
