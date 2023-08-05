#!/bin/bash

while sleep 0.01; do wget -q -O- http://localhost:8000; done
