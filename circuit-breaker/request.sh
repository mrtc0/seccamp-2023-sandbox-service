#!/bin/bash

while sleep 1; do curl -s 'http://localhost:9091/' | jq '.upstream_calls."http://localhost:9092" | .name, .code'; done
