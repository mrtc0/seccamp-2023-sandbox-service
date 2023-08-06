#!/bin/bash

kubectl apply -f ./manifests/namespace.yaml
kubectl apply -Rf ./manifests/
