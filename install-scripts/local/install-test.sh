#!/bin/bash

## install test

# target
cf api https://localhost --skip-ssl-validation

# login
cf auth kind-sidecar

# create org & space
cf create-org org && cf create-space -o org space && cf target -o org

# sample app push
cf push sample-app -p ../support-files/sample-app/

# sample app check
curl -k https://sample-app.apps-127-0-0-1.nip.io


