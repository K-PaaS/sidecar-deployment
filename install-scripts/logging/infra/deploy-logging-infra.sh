#!/bin/bash

source infra-variable.yml


## COMMON
ENV=${LOGGING_NAMESPACE} yq e -i '.metadata.name = env(ENV)' etc/namespace.yaml

ENV=${LOGGING_NAMESPACE} yq e -i '.metadata.namespace = env(ENV)' influxdb/influxdb-config.yaml
ENV=${LOGGING_NAMESPACE} yq e -i '.metadata.namespace = env(ENV)' influxdb/influxdb-secret.yaml
ENV=${LOGGING_NAMESPACE} yq e -i '.metadata.namespace = env(ENV)' influxdb/influxdb-cert.yaml
ENV=${LOGGING_NAMESPACE} yq e -i '.metadata.namespace = env(ENV)' influxdb/influxdb-service.yaml
ENV=${LOGGING_NAMESPACE} yq e -i '.metadata.namespace = env(ENV)' influxdb/influxdb-statefulset.yaml

## INFLUXDB
ENV=${INFLUXDB_USERNAME} yq e -i '.stringData.username = env(ENV)' influxdb/influxdb-secret.yaml
ENV=${INFLUXDB_PASSWORD} yq e -i '.stringData.password = env(ENV)' influxdb/influxdb-secret.yaml
ENV=${INFLUXDB_DATABASE} yq e -i '.stringData.database = env(ENV)' influxdb/influxdb-secret.yaml
ENV=${INFLUXDB_HTTP_PORT} yq e -i '.stringData.httpPort = strenv(ENV)' influxdb/influxdb-secret.yaml
ENV=${INFLUXDB_HTTPS_ENABLED} yq e -i '.stringData.httpsEnabled = strenv(ENV)' influxdb/influxdb-secret.yaml
ENV=${INFLUXDB_RETENTION_POLICY} yq e -i '.stringData.retentionPolicy = env(ENV)' influxdb/influxdb-secret.yaml

ENV=${INFLUXDB_HTTP_PORT} yq e -i '.spec.ports[0].port = env(ENV)' influxdb/influxdb-service.yaml
ENV=${INFLUXDB_HTTP_PORT} yq e -i '.spec.ports[0].targetPort = env(ENV)' influxdb/influxdb-service.yaml

sed -i -e 's/<%= p("influxdb.http.bind-address") %>/'${INFLUXDB_HTTP_PORT}'/g' influxdb/influxdb-config.yaml
sed -i -e 's/<%= p("influxdb.https_enabled") %>/'${INFLUXDB_HTTPS_ENABLED}'/g' influxdb/influxdb-config.yaml


kubectl apply -f etc
kubectl apply -f influxdb
