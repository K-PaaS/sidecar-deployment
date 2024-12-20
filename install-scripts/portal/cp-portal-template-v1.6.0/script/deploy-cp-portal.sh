#!/bin/bash
source cp-portal-vars.sh
declare -A DEPLOY_CONFIG
DEPLOY_CONFIG[IPV6_ENABLED]=true
DEPLOY_CONFIG[INGRESS_ENABLED]=true
DEPLOY_CONFIG[EXPOSE_TYPE]="ingress"
CMD_CREATE_TLS_SECRET="kubectl create secret tls $TLS_SECRET --cert=../certs/${HOST_DOMAIN}.crt  --key=../certs/${HOST_DOMAIN}.key"

# Create cluster-admin token
kubectl create sa $K8S_CLUSTER_ADMIN -n $K8S_CLUSTER_ADMIN_NAMESPACE
kubectl create clusterrolebinding $K8S_CLUSTER_ADMIN --clusterrole=cluster-admin --serviceaccount=$K8S_CLUSTER_ADMIN_NAMESPACE:$K8S_CLUSTER_ADMIN
K8S_CLUSTER_ADMIN_TOKEN=$(kubectl create token $K8S_CLUSTER_ADMIN --duration=999999h -n $K8S_CLUSTER_ADMIN_NAMESPACE)

# Create a vault bound cidr
VAULT_BOUND_CIDR_ARR=($(kubectl get pods -n $INGRESS_NAMESPACE --selector=$INGRESS_CONTROLLER_SELECTOR --field-selector=status.phase=Running -o jsonpath='{range .items[*]}{@.status.podIP}{"/16"}{"\t"}{end}'))
printf -v VAULT_BOUND_CIDR '"%s",' "${VAULT_BOUND_CIDR_ARR[@]}"
VAULT_BOUND_CIDR="${VAULT_BOUND_CIDR%,}"

# Copy the directory
cp -r ../vault_orig ../vault
cp -r ../values_orig ../values

# Set a iaas type
if [[ $HOST_CLUSTER_IAAS_TYPE -lt 1 ]] || [[ $HOST_CLUSTER_IAAS_TYPE -gt ${#IAAS_TYPE[@]} ]]
then
  HOST_CLUSTER_IAAS_TYPE="1"
fi
# ipv6Enabled set to false if iaas type is NAVER
if [[ $((HOST_CLUSTER_IAAS_TYPE -1)) -eq 2 ]]
then
  DEPLOY_CONFIG[IPV6_ENABLED]=false
fi

# Replace values
REPOSITORY_HOST=$(echo $REPOSITORY_URL | awk -F[/:] '{print $4}')
find ../vault -name "payload.json" -exec sed -i "s@{VAULT_BOUND_CIDR}@$VAULT_BOUND_CIDR@g" {} \;
find ../values -name "*.json" -exec sed -i "s@{CP_PORTAL_URL}@$CP_PORTAL_URL@g" {} \;
find ../values -name "*.json" -exec sed -i "s@{CP_SERVICE_PIPELINE_URL}@$CP_SERVICE_PIPELINE_URL@g" {} \;
find ../values -name "*.json" -exec sed -i "s@{CP_SERVICE_SOURCE_CONTROL_URL}@$CP_SERVICE_SOURCE_CONTROL_URL@g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{K8S_MASTER_NODE_IP}/$K8S_MASTER_NODE_IP/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{HOST_CLUSTER_IAAS_TYPE}/${IAAS_TYPE[$HOST_CLUSTER_IAAS_TYPE -1]}/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{HOST_DOMAIN}/$HOST_DOMAIN/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{IMAGE_TAGS}/$IMAGE_TAGS/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{IMAGE_PULL_POLICY}/$IMAGE_PULL_POLICY/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{IMAGE_PULL_SECRET}/$IMAGE_PULL_SECRET/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{TLS_SECRET}/$TLS_SECRET/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{SERVICE_TYPE}/$SERVICE_TYPE/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{SERVICE_PROTOCOL}/$SERVICE_PROTOCOL/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{INGRESS_CLASS_NAME}/$INGRESS_CLASS_NAME/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{VAULT_NAMESPACE}/${NAMESPACE[0]}/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s@{VAULT_URL}@$VAULT_URL@g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{VAULT_HOST}/$(echo $VAULT_URL | awk -F[/:] '{print $4}')/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{VAULT_ROLE_NAME}/$VAULT_ROLE_NAME/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{VAULT_STORAGECLASS}/$K8S_STORAGECLASS/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{REPOSITORY_NAMESPACE}/${NAMESPACE[2]}/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s@{REPOSITORY_URL}@$REPOSITORY_URL@g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{REPOSITORY_HOST}/$REPOSITORY_HOST/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{REPOSITORY_USERNAME}/$REPOSITORY_USERNAME/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{REPOSITORY_PASSWORD}/$REPOSITORY_PASSWORD/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{REPOSITORY_PROJECT_NAME}/$REPOSITORY_PROJECT_NAME/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{REPOSITORY_STORAGECLASS}/$K8S_STORAGECLASS/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{DATABASE_NAMESPACE}/${NAMESPACE[1]}/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{DATABASE_URL}/$DATABASE_URL/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{DATABASE_HOST}/$(echo "${DATABASE_URL%:*}")/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{DATABASE_PORT}/$(echo "${DATABASE_URL#*:}")/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{DATABASE_USER_ID}/$DATABASE_USER_ID/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{DATABASE_USER_PASSWORD}/$DATABASE_USER_PASSWORD/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{DATABASE_TERRAMAN_ID}/$DATABASE_TERRAMAN_ID/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{DATABASE_TERRAMAN_PASSWORD}/$DATABASE_TERRAMAN_PASSWORD/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{DATABASE_STORAGECLASS}/$K8S_STORAGECLASS/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{KEYCLOAK_NAMESPACE}/${NAMESPACE[3]}/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s@{KEYCLOAK_URL}@$KEYCLOAK_URL@g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{KEYCLOAK_DB_VENDOR}/$KEYCLOAK_DB_VENDOR/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{KEYCLOAK_DB_SCHEMA}/$KEYCLOAK_DB_SCHEMA/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{KEYCLOAK_ADMIN_USERNAME}/$KEYCLOAK_ADMIN_USERNAME/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{KEYCLOAK_ADMIN_PASSWORD}/$KEYCLOAK_ADMIN_PASSWORD/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{KEYCLOAK_SESSIONS_COUNT}/$KEYCLOAK_SESSIONS_COUNT/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{KEYCLOAK_CP_REALM}/$KEYCLOAK_CP_REALM/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{KEYCLOAK_CP_REALM_ID}/$KEYCLOAK_CP_REALM_ID/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{KEYCLOAK_CP_CLIENT_ID}/$KEYCLOAK_CP_CLIENT_ID/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{KEYCLOAK_CP_CLIENT_SECRET}/$KEYCLOAK_CP_CLIENT_SECRET/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{KEYCLOAK_HOST}/$(echo $KEYCLOAK_URL | awk -F[/:] '{print $4}')/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{CHART_REPOSITORY_NAME}/$CHART_REPOSITORY_NAME/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s@{CHART_REPOSITORY_URL}@$CHART_REPOSITORY_URL@g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{CHART_REPOSITORY_HOST}/$(echo $CHART_REPOSITORY_URL | awk -F[/:] '{print $4}')/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{CHART_REPOSITORY_STORAGECLASS}/$K8S_STORAGECLASS/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{CHAOS_MESH_NAMESPACE}/${NAMESPACE[6]}/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{NAMESPACE}/${NAMESPACE[4]}/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{HOST_CLUSTER_NAME}/$HOST_CLUSTER_NAME/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s@{CP_PORTAL_URL}@$CP_PORTAL_URL@g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{CP_PORTAL_HOST}/$(echo $CP_PORTAL_URL | awk -F[/:] '{print $4}')/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{CP_PORTAL_STORAGECLASS}/$K8S_STORAGECLASS/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{CP_SERVICE_PIPELINE_NAMESPACE}/$CP_SERVICE_PIPELINE_NAMESPACE/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{CP_SERVICE_SOURCE_CONTROL_NAMESPACE}/$CP_SERVICE_SOURCE_CONTROL_NAMESPACE/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{CP_CERT_SETUP_NAME}/${CHART_NAME[8]}/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{CP_CERT_SETUP_NAMESPACE}/$CP_CERT_SETUP_NAMESPACE/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{IPV6_ENABLED}/${DEPLOY_CONFIG[IPV6_ENABLED]}/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{INGRESS_ENABLED}/${DEPLOY_CONFIG[INGRESS_ENABLED]}/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{EXPOSE_TYPE}/${DEPLOY_CONFIG[EXPOSE_TYPE]}/g" {} \;

# Generate self signed certificate
chmod +x ./gen-cert.sh
. ./gen-cert.sh

# Setup the certificate in cluster
echo "[Setup the cert in cluster]..."
helm install -f ../values/${CHART_NAME[8]}.yaml ${CHART_NAME[8]} ../charts/${CHART_NAME[8]}.tgz -n $CP_CERT_SETUP_NAMESPACE --set-literal data.target.cert="$(cat ../certs/${HOST_DOMAIN}.crt)"
while :
do
  POD_COUNT=$((kubectl get pods -n $CP_CERT_SETUP_NAMESPACE -l $CP_CERT_SETUP_SELECTOR --field-selector status.phase!=Running --no-headers | wc -l) 2> /dev/null)
  echo "[remaining: $POD_COUNT] Adding a cert to each node's container runtime..."
  if [[ $POD_COUNT -lt 1 ]]; then
    echo "Completed..."
    break
  fi
  sleep 5
done
sudo cp ../certs/${HOST_DOMAIN}.crt /usr/local/share/ca-certificates/
sudo update-ca-certificates

# Deploy the vault
chmod +x ../vault/deploy-vault.sh
. ../vault/deploy-vault.sh
find ../values -name "*.yaml" -exec sed -i "s/{VAULT_ROLE_ID}/$VAULT_ROLE_ID/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{VAULT_SECRET_ID}/$VAULT_SECRET_ID/g" {} \;

# Deploy the mariadb
kubectl create namespace ${NAMESPACE[1]}
kubectl apply -f ../values/${CHART_NAME[1]}-configmap.yaml -n ${NAMESPACE[1]}
helm install -f ../values/${CHART_NAME[1]}.yaml ${CHART_NAME[1]} ../charts/${CHART_NAME[1]}.tgz -n ${NAMESPACE[1]}

# Deploy the harbor
kubectl create namespace ${NAMESPACE[2]}
$CMD_CREATE_TLS_SECRET -n ${NAMESPACE[2]}
helm install -f ../values/${CHART_NAME[2]}.yaml ${CHART_NAME[2]} ../charts/${CHART_NAME[2]}.tgz -n ${NAMESPACE[2]}
while :
do
  REPOSITORY_HTTP_CODE=$(curl -L -k -s -o /dev/null -w "%{http_code}\n" $REPOSITORY_URL/api/v2.0/projects)
  echo "[$REPOSITORY_HTTP_CODE] Please wait a few minutes for the Harbor deployment to finish..."
  if [ $REPOSITORY_HTTP_CODE -eq 200 ]; then
    break
  fi
  sleep 10
done

curl -u $REPOSITORY_USERNAME:$REPOSITORY_PASSWORD -k $REPOSITORY_URL/api/v2.0/projects -XPOST --data-binary "{\"project_name\": \"$REPOSITORY_PROJECT_NAME\", \"public\": false}" -H "Content-Type: application/json" -i
sudo podman login $REPOSITORY_HOST --username $REPOSITORY_USERNAME --password $REPOSITORY_PASSWORD
# Push images and charts to harbor
for IMAGE in ${IMAGE_NAME[@]}
do
    sudo podman load -i ../images/$IMAGE.tar.gz
    sudo podman tag localhost:5000/container-platform/$IMAGE $REPOSITORY_HOST/$REPOSITORY_PROJECT_NAME/$IMAGE
    sudo podman push $REPOSITORY_HOST/$REPOSITORY_PROJECT_NAME/$IMAGE
done
helm registry login $REPOSITORY_HOST --username $REPOSITORY_USERNAME --password $REPOSITORY_PASSWORD
for CHART in ${CHART_NAME[@]}
do
  helm push ../charts/$CHART.tgz oci://$REPOSITORY_HOST/$REPOSITORY_PROJECT_NAME
done

# Deploy the keycloak
kubectl create namespace ${NAMESPACE[3]}
$CMD_CREATE_TLS_SECRET -n ${NAMESPACE[3]}
kubectl create configmap $KEYCLOAK_CP_REALM --from-file=../values/$KEYCLOAK_CP_REALM.json -n ${NAMESPACE[3]}
helm install -f ../values/${CHART_NAME[3]}.yaml ${CHART_NAME[3]} ../charts/${CHART_NAME[3]}.tgz -n ${NAMESPACE[3]}

# Deploy the chartmuseum
kubectl create namespace ${NAMESPACE[5]}
$CMD_CREATE_TLS_SECRET -n ${NAMESPACE[5]}
helm install -f ../values/${CHART_NAME[5]}.yaml ${CHART_NAME[5]} ../charts/${CHART_NAME[5]}.tgz -n ${NAMESPACE[5]}
while :
do
  CHART_REPOSITORY_HTTP_CODE=$(curl -L -k -s -o /dev/null -w "%{http_code}\n" $CHART_REPOSITORY_URL/index.yaml)
  echo "[$CHART_REPOSITORY_HTTP_CODE] Check the status of ChartMuseum..."
  if [ $CHART_REPOSITORY_HTTP_CODE -eq 200 ]; then
    break
  fi
  sleep 5
done
helm plugin install https://github.com/chartmuseum/helm-push.git
helm repo add $CHART_REPOSITORY_NAME $CHART_REPOSITORY_URL

# Deploy the chaos-mesh
kubectl create namespace ${NAMESPACE[6]}
helm install -f ../values/${CHART_NAME[6]}.yaml ${CHART_NAME[6]} ../charts/${CHART_NAME[6]}.tgz -n ${NAMESPACE[6]}

# Deploy the cp-portal
kubectl create namespace ${NAMESPACE[4]}
helm install -f ../values/${RELEASE_NAME}.yaml ${RELEASE_NAME} ../charts/${CHART_NAME[4]}.tgz -n ${NAMESPACE[4]} \
     --set-literal tlsSecret.tls.crt=$(cat ../certs/${HOST_DOMAIN}.crt | base64 -w 0) \
     --set-literal tlsSecret.tls.key=$(cat ../certs/${HOST_DOMAIN}.key | base64 -w 0) \
     --set-literal configmap.data.CHART_REPO_CRT=$(cat ../certs/${HOST_DOMAIN}.crt | base64 -w 0)

# Uninstall cp-cert-setup
helm uninstall ${CHART_NAME[8]} -n $CP_CERT_SETUP_NAMESPACE
