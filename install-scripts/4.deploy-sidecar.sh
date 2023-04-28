#!/bin/bash

source variables.yml

kapp deploy -a sidecar -f manifest/sidecar-rendered.yml

if [[ ${use_logging_service} = "true" ]]; then
        source enable-logging-service.sh
fi