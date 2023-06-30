#!/bin/bash

source variables.yml
source logging/logging-service-variable.yml

kapp deploy -a sidecar -f manifest/sidecar-rendered.yml

if [[ ${ENABLE_LOGGING_SERVICE} = "true" ]]; then
        cd logging/infra
        source deploy-logging-infra.sh

        cd ..
        source enable-logging-service.sh

        cd ..
fi
