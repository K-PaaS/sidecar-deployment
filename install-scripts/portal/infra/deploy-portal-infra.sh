#!/bin/bash
  
source infra-variable.yml
source ../portal-app-variable.yml

# COMMON
ENV=$NAMESPACE_NAME yq e -i '.metadata.name = env(ENV)' etc/namespace.yaml

ENV=$NAMESPACE_NAME yq e -i '.metadata.namespace = env(ENV)' mariadb/mariadb-secret.yaml
ENV=$NAMESPACE_NAME yq e -i '.metadata.namespace = env(ENV)' mariadb/mariadb-service.yaml
ENV=$NAMESPACE_NAME yq e -i '.metadata.namespace = env(ENV)' mariadb/mariadb-statefulset.yaml
sed -i '/namespace:/ c\  namespace: '$NAMESPACE_NAME'' mariadb/mariadb-initsql.yaml

ENV=$NAMESPACE_NAME yq e -i '.metadata.namespace = env(ENV)' saio/saio-service.yaml
ENV=$NAMESPACE_NAME yq e -i '.metadata.namespace = env(ENV)' saio/saio-statefulset.yaml


# mariadb-secret.yaml
## MARIADB PASSWORD
ENV=$(echo $MARIADB_PASSWORD | base64 ) yq e -i '(.data.password = env(ENV)' mariadb/mariadb-secret.yaml


# mariadb-service.yaml
## MARIADB_SERVICE_PORT
ENV=$MARIADB_SERVICE_PORT yq e -i '(.spec.ports[0].port = env(ENV)' mariadb/mariadb-service.yaml

## MARIADB_CONTAINER_PORT
ENV=$MARIADB_CONTAINER_PORT yq e -i '(.spec.ports[0].targetPort = env(ENV)' mariadb/mariadb-service.yaml


# mariadb-statefulset.yaml
## MARIADB_CONTAINER_PORT
ENV=$MARIADB_CONTAINER_PORT yq e -i '(.spec.template.spec.containers[0].ports[0].containerPort = env(ENV)' mariadb/mariadb-statefulset.yaml


SYSTEM_DOMAIN=$(yq e '.system_domain' $SIDECAR_VALUES_PATH)
APP_DOMAIN=$(yq e '.app_domains[0]' $SIDECAR_VALUES_PATH)
# mariadb-initsql.yaml
## PORTAL_NAME
sed -i -e 's/<%= p("portal_default.name") %>/'$PORTAL_NAME'/g' mariadb/mariadb-initsql.yaml

## PORTAL_GATEWAY_URL
sed -i -e 's/<%= p("portal_default.url") %>/http:\/\/portal-gateway.'$APP_DOMAIN'/g' mariadb/mariadb-initsql.yaml

## PORTAL_UAA_URL
sed -i -e 's/<%= p("portal_default.uaa_url") %>/https:\/\/uaa.'$SYSTEM_DOMAIN'/g' mariadb/mariadb-initsql.yaml

## PORTAL_HEADER_AUTH
sed -i -e "s/<%= p(\"portal_default.header_auth\") %>/$PORTAL_HEADER_AUTH/g" mariadb/mariadb-initsql.yaml

## PORTAL_DESC
sed -i -e "s/<%= p(\"portal_default.desc\") %>/$PORTAL_DESC/g" mariadb/mariadb-initsql.yaml

# saio-service.yaml
## KEYSTONE_SERVICE_PORT
ENV=$KEYSTONE_SERVICE_PORT yq e -i '(.spec.ports[0].port = env(ENV)' saio/saio-service.yaml

##KEYSTONE_CONTAINER_PORT
ENV=$KEYSTONE_CONTAINER_PORT yq e -i '(.spec.ports[0].targetPort = env(ENV)' saio/saio-service.yaml

##PROXY_SERVICE_PORT
ENV=$PROXY_SERVICE_PORT yq e -i '(.spec.ports[1].port = env(ENV)' saio/saio-service.yaml

##PROXY_TARGET_PORT
ENV=$PROXY_TARGET_PORT yq e -i '(.spec.ports[1].targetPort = env(ENV)' saio/saio-service.yaml


# saio-statefulset.yaml
## MARIADB_SERVICE_PORT
ENV=$MARIADB_SERVICE_PORT yq e -i '(.spec.template.spec.containers[0].env[1].value = env(ENV)' saio/saio-statefulset.yaml


## MARIADB_PASSWORD
ENV=$MARIADB_PASSWORD yq e -i '(.spec.template.spec.containers[0].env[2].value = env(ENV)' saio/saio-statefulset.yaml

## SWIFT_ADDRESS
ENV="openstack-swift-keystone-docker."$NAMESPACE_NAME".svc.cluster.local" yq e -i '(.spec.template.spec.containers[0].env[3].value = env(ENV)' saio/saio-statefulset.yaml

## KEYSTONE_CONTAINER_PORT
ENV=$KEYSTONE_CONTAINER_PORT yq e -i '(.spec.template.spec.containers[0].env[4].value = env(ENV)' saio/saio-statefulset.yaml

## PROXY_TARGET_PORT
ENV=$PROXY_TARGET_PORT yq e -i '(.spec.template.spec.containers[0].env[5].value = env(ENV)' saio/saio-statefulset.yaml

## PORTAL_OBJECTSTORAGE_TENANTNAME
ENV=$PORTAL_OBJECTSTORAGE_TENANTNAME yq e -i '(.spec.template.spec.containers[0].env[7].value = env(ENV)' saio/saio-statefulset.yaml

## PORTAL_OBJECTSTORAGE_USERNAME
ENV=$PORTAL_OBJECTSTORAGE_USERNAME yq e -i '(.spec.template.spec.containers[0].env[8].value = env(ENV)' saio/saio-statefulset.yaml

## PORTAL_OBJECTSTORAGE_PASSWORD
ENV=$PORTAL_OBJECTSTORAGE_PASSWORD yq e -i '(.spec.template.spec.containers[0].env[9].value = env(ENV)' saio/saio-statefulset.yaml


# allow-cf-db-ingress-from-cf-workloads.yaml
## ORG
sed -i '/cloudfoundry.org\/org_name/ c\          cloudfoundry.org\/org_name: '$PORTAL_ORG_NAME'' etc/allow-cf-db-ingress-from-cf-workloads.yaml

## SPACE
sed -i '/cloudfoundry.org\/space_name/ c\          cloudfoundry.org\/space_name: '$PORTAL_SPACE_NAME'' etc/allow-cf-db-ingress-from-cf-workloads.yaml


# portal-sidecar.yaml
## ORG
sed -i '/cloudfoundry.org\/org_name/ c\      cloudfoundry.org\/org_name: '$PORTAL_ORG_NAME'' etc/portal-sidecar.yaml

## SPACE
sed -i '/cloudfoundry.org\/space_name/ c\      cloudfoundry.org\/space_name: '$PORTAL_SPACE_NAME'' etc/portal-sidecar.yaml


kubectl apply -f etc
kubectl apply -f mariadb
sleep 5
kubectl apply -f saio
