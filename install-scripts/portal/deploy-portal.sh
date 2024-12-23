#!/bin/bash

source portal-deploy-variables.yml

# SCRIPT START

## VARIABLES SETTING
### SIDECAR VARIABLES
SIDECAR_API_URL=$(helm get values $HELM_SIDECAR_NAME -n $HELM_SIDECAR_NAMESPACE | yq e '.api.apiServer.url')
SIDECAR_TOKEN_KIND=bearer
SIDECAR_ROOTNAMESPACE=$(helm get values $HELM_SIDECAR_NAME -n $HELM_SIDECAR_NAMESPACE | yq e '.rootNamespace')
SIDECAR_ROLEADMIN=korifi-controllers-admin
SIDECAR_ROLEUSER=korifi-controllers-root-namespace-user
SIDECAR_ROLEORGMANAGER=korifi-controllers-organization-manager
SIDECAR_ROLEORGUSER=korifi-controllers-organization-user
SIDECAR_PORTAL_API_URL=http://$PORTAL_API_NAME.$(helm get values $HELM_SIDECAR_NAME -n $HELM_SIDECAR_NAMESPACE | yq e '.defaultAppDomainName')
### CP-PORTAL VARIABLES
CP_PORTAL_API_URI=http://$PORTAL_API_NAME.$(helm get values $HELM_SIDECAR_NAME -n $HELM_SIDECAR_NAMESPACE | yq e '.defaultAppDomainName')
CP_PORTAL_COMMON_API_URI=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.CP_PORTAL_COMMON_API_URI')
CP_PORTAL_METRIC_COLLECTOR_API_URI=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.CP_PORTAL_METRIC_COLLECTOR_API_URI')
CP_PORTAL_TERRAMAN_API_URI=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.CP_PORTAL_TERRAMAN_API_URI')
CP_PORTAL_UI_URI=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.CP_PORTAL_UI_URI')
DATABASE_TERRAMAN_ID=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.DATABASE_TERRAMAN_ID')
DATABASE_TERRAMAN_PASSWORD=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.DATABASE_TERRAMAN_PASSWORD')
DATABASE_URL=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.DATABASE_URL')
DATABASE_USER_ID=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.DATABASE_USER_ID')
DATABASE_USER_PASSWORD=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.DATABASE_USER_PASSWORD')
K8S_MASTER_NODE_IP=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.K8S_MASTER_NODE_IP')
KEYCLOAK_ADMIN_PASSWORD=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.KEYCLOAK_ADMIN_PASSWORD')
KEYCLOAK_ADMIN_USERNAME=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.KEYCLOAK_ADMIN_USERNAME')
KEYCLOAK_CP_CLIENT_ID=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.KEYCLOAK_CP_CLIENT_ID')
KEYCLOAK_CP_CLIENT_SECRET=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.KEYCLOAK_CP_CLIENT_SECRET')
KEYCLOAK_CP_REALM=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.KEYCLOAK_CP_REALM')
KEYCLOAK_DB_SCHEMA=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.KEYCLOAK_DB_SCHEMA')
KEYCLOAK_URI=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.KEYCLOAK_URI')
REPOSITORY_URL=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.REPOSITORY_URL')
VAULT_ROLE_ID=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.VAULT_ROLE_ID')
VAULT_ROLE_NAME=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.VAULT_ROLE_NAME')
VAULT_SECRET_ID=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.VAULT_SECRET_ID')
VAULT_URL=$(helm get values $HELM_CP_PORTAL_RESOURCE_NAME -n $HELM_CP_PORTAL_NAMESPACE | yq e '.configmap.data.VAULT_URL')


echo "SIDECAR_API_URL: $SIDECAR_API_URL"
echo "SIDECAR_TOKEN_KIND: $SIDECAR_TOKEN_KIND"
echo "SIDECAR_ROOTNAMESPACE: $SIDECAR_ROOTNAMESPACE"
echo "SIDECAR_ROLEADMIN: $SIDECAR_ROLEADMIN"
echo "SIDECAR_ROLEUSER: $SIDECAR_ROLEUSER"
echo "SIDECAR_PORTAL_API_URL: $SIDECAR_PORTAL_API_URL"
echo "CP_PORTAL_API_URI: $CP_PORTAL_API_URI"
echo "CP_PORTAL_COMMON_API_URI: $CP_PORTAL_COMMON_API_URI"
echo "CP_PORTAL_METRIC_COLLECTOR_API_URI: $CP_PORTAL_METRIC_COLLECTOR_API_URI"
echo "CP_PORTAL_TERRAMAN_API_URI: $CP_PORTAL_TERRAMAN_API_URI"
echo "CP_PORTAL_UI_URI: $CP_PORTAL_UI_URI"
echo "DATABASE_TERRAMAN_ID: $DATABASE_TERRAMAN_ID"
echo "DATABASE_TERRAMAN_PASSWORD: $DATABASE_TERRAMAN_PASSWORD"
echo "DATABASE_URL: $DATABASE_URL"
echo "DATABASE_USER_ID: $DATABASE_USER_ID"
echo "DATABASE_USER_PASSWORD: $DATABASE_USER_PASSWORD"
echo "K8S_MASTER_NODE_IP: $K8S_MASTER_NODE_IP"
echo "KEYCLOAK_ADMIN_PASSWORD: $KEYCLOAK_ADMIN_PASSWORD"
echo "KEYCLOAK_ADMIN_USERNAME: $KEYCLOAK_ADMIN_USERNAME"
echo "KEYCLOAK_CP_CLIENT_ID: $KEYCLOAK_CP_CLIENT_ID"
echo "KEYCLOAK_CP_CLIENT_SECRET: $KEYCLOAK_CP_CLIENT_SECRET"
echo "KEYCLOAK_CP_REALM: $KEYCLOAK_CP_REALM"
echo "KEYCLOAK_DB_SCHEMA: $KEYCLOAK_DB_SCHEMA"
echo "KEYCLOAK_URI: $KEYCLOAK_URI"
echo "REPOSITORY_URL: $REPOSITORY_URL"
echo "VAULT_ROLE_ID: $VAULT_ROLE_ID"
echo "VAULT_ROLE_NAME: $VAULT_ROLE_NAME"
echo "VAULT_SECRET_ID: $VAULT_SECRET_ID"
echo "VAULT_URL: $VAULT_URL"
echo "TARGET_CLUSTER: $TARGET_CLUSTER"
echo "CP_PORTAL_ADMIN_NAME: $CP_PORTAL_ADMIN_NAME"
echo "SIDECAR_ROLEORGMANAGER: $SIDECAR_ROLEORGMANAGER"
echo "SIDECAR_ROLEORGUSER: $SIDECAR_ROLEORGUSER"


if [ -e ./portal-app-variables.yml ]; then
    while true; do
        read -p "Do you want rewrite portal-app-variables.yml file? (y/n) " yn
        case $yn in
            [Yy]* ) break;;
            [Nn]* ) exit;;
            * ) echo "Please answer y or n.";;
        esac
    done
fi

cat <<EOF > portal-app-variables.yml
SIDECAR_API_URL: $SIDECAR_API_URL
SIDECAR_TOKEN_KIND: $SIDECAR_TOKEN_KIND
SIDECAR_ROOTNAMESPACE: $SIDECAR_ROOTNAMESPACE
SIDECAR_ROLEADMIN: $SIDECAR_ROLEADMIN
SIDECAR_ROLEUSER: $SIDECAR_ROLEUSER
SIDECAR_PORTAL_API_URL: $SIDECAR_PORTAL_API_URL
CP_PORTAL_API_URI: $CP_PORTAL_API_URI
CP_PORTAL_COMMON_API_URI: $CP_PORTAL_COMMON_API_URI
CP_PORTAL_METRIC_COLLECTOR_API_URI: $CP_PORTAL_METRIC_COLLECTOR_API_URI
CP_PORTAL_PROVIDER_TYPE: $CP_PORTAL_PROVIDER_TYPE
CP_PORTAL_TERRAMAN_API_URI: $CP_PORTAL_TERRAMAN_API_URI
CP_PORTAL_UI_URI: $CP_PORTAL_UI_URI
DATABASE_TERRAMAN_ID: $DATABASE_TERRAMAN_ID
DATABASE_TERRAMAN_PASSWORD: $DATABASE_TERRAMAN_PASSWORD
DATABASE_URL: $DATABASE_URL
DATABASE_USER_ID: $DATABASE_USER_ID
DATABASE_USER_PASSWORD: $DATABASE_USER_PASSWORD
K8S_MASTER_HOST_KEY: $K8S_MASTER_HOST_KEY
K8S_MASTER_NODE_IP: $K8S_MASTER_NODE_IP
KEYCLOAK_ADMIN_PASSWORD: $KEYCLOAK_ADMIN_PASSWORD
KEYCLOAK_ADMIN_USERNAME: $KEYCLOAK_ADMIN_USERNAME
KEYCLOAK_CP_CLIENT_ID: $KEYCLOAK_CP_CLIENT_ID
KEYCLOAK_CP_CLIENT_SECRET: $KEYCLOAK_CP_CLIENT_SECRET
KEYCLOAK_CP_REALM: $KEYCLOAK_CP_REALM
KEYCLOAK_DB_SCHEMA: $KEYCLOAK_DB_SCHEMA
KEYCLOAK_URI: $KEYCLOAK_URI
REPOSITORY_URL: $REPOSITORY_URL
VAULT_ROLE_ID: $VAULT_ROLE_ID
VAULT_ROLE_NAME: $VAULT_ROLE_NAME
VAULT_SECRET_ID: $VAULT_SECRET_ID
VAULT_URL: $VAULT_URL
TARGET_CLUSTER: $TARGET_CLUSTER
CP_PORTAL_ADMIN_NAME: $CP_PORTAL_ADMIN_NAME
SIDECAR_ROLEORGMANAGER: $SIDECAR_ROLEORGMANAGER
SIDECAR_ROLEORGUSER: $SIDECAR_ROLEORGUSER
EOF

# CP Portal Admin Sidecar Role Binding
cd ../support-files/user
echo $K8S_CLUSTER_ADMIN_NAMESPACE | source binding-admin.sh sa $K8S_CLUSTER_ADMIN
cd ../../portal


# git submodule
git submodule init
git submodule update

# git checkout
cd sidecar-portal-api
git checkout v1.0.0


cd ../sidecar-portal-ui
git checkout v1.0.1

cd ..

# sidecar config change
export KUBECONFIG=$SIDECAR_ADMIN_KUBECONFIG

cf create-org $ORG_NAME
cf create-space -o $ORG_NAME $SPACE_NAME

cf target -o $ORG_NAME -s $SPACE_NAME

# cf push
cd sidecar-portal-api
cf push --vars-file ../portal-app-variables.yml $PORTAL_API_NAME

cd ../sidecar-portal-ui
cf push --vars-file ../portal-app-variables.yml $PORTAL_UI_NAME

cd ..

cf apps
