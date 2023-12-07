#!/bin/bash

## deploy korifi

kubectl apply -f deploy-sidecar-kind.yaml

echo "track the job progress command : 'kubectl -n sidecar-installer logs --follow job/install-sidecar'"
