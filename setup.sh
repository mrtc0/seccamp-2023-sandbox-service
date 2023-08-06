#!/bin/bash

kubectl apply -f ./manifests/namespace.yaml
kubectl apply -f ./manifests/api.yaml
kubectl apply -f ./manifests/backend.yaml
kubectl apply -f ./manifests/payments.yaml
kubectl apply -f ./manifests/worker.yaml
