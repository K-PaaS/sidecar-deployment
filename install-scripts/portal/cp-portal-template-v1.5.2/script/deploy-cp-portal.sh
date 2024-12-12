#!/bin/bash

source cp-portal-vars.sh

# create cluster-admin token
kubectl create sa $K8S_CLUSTER_ADMIN -n $K8S_CLUSTER_ADMIN_NAMESPACE
kubectl create clusterrolebinding $K8S_CLUSTER_ADMIN --clusterrole=cluster-admin --serviceaccount=$K8S_CLUSTER_ADMIN_NAMESPACE:$K8S_CLUSTER_ADMIN
K8S_CLUSTER_ADMIN_TOKEN=$(kubectl create token $K8S_CLUSTER_ADMIN --duration=999999h -n $K8S_CLUSTER_ADMIN_NAMESPACE)

# create a vault bound cidr
VAULT_BOUND_CIDR_ARR=($(kubectl get pods -n $CP_DEFAULT_INGRESS_NAMESPACE --selector=$CP_DEFAULT_INGRESS_CONTROLLER_SELECTOR --field-selector=status.phase=Running -o jsonpath='{range .items[*]}{@.status.podIP}{"/16"}{"\t"}{end}'))
printf -v VAULT_BOUND_CIDR '"%s",' "${VAULT_BOUND_CIDR_ARR[@]}"
VAULT_BOUND_CIDR="${VAULT_BOUND_CIDR%,}"

# copy the directory
cp -r ../vault_orig ../vault
cp -r ../keycloak_orig ../keycloak
cp -r ../values_orig ../values

# set a iaas type
IPV6_ENABLED=true
if [[ $HOST_CLUSTER_IAAS_TYPE -lt 1 ]] || [[ $HOST_CLUSTER_IAAS_TYPE -gt ${#IAAS_TYPE[@]} ]]
then
  HOST_CLUSTER_IAAS_TYPE="1"
fi
# ipv6Enabled set to false if iaas type is NAVER
if [[ $((HOST_CLUSTER_IAAS_TYPE -1)) -eq 2 ]]
then
  IPV6_ENABLED=false
fi

# replace values
REPOSITORY_HOST=$(echo $REPOSITORY_URL | awk -F[/:] '{print $4}')
find ../vault -name "payload.json" -exec sed -i "s@{VAULT_BOUND_CIDR}@$VAULT_BOUND_CIDR@g" {} \;
find ../keycloak -name "*.json" -exec sed -i "s@{CP_PORTAL_URL}@$CP_PORTAL_URL@g" {} \;
find ../keycloak -name "*.json" -exec sed -i "s@{CP_SERVICE_PIPELINE_URL}@$CP_SERVICE_PIPELINE_URL@g" {} \;
find ../keycloak -name "*.json" -exec sed -i "s@{CP_SERVICE_SOURCE_CONTROL_URL}@$CP_SERVICE_SOURCE_CONTROL_URL@g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{K8S_MASTER_NODE_IP}/$K8S_MASTER_NODE_IP/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{HOST_CLUSTER_IAAS_TYPE}/${IAAS_TYPE[$HOST_CLUSTER_IAAS_TYPE -1]}/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{HOST_DOMAIN}/$HOST_DOMAIN/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{IMAGE_TAGS}/$IMAGE_TAGS/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{IMAGE_PULL_POLICY}/$IMAGE_PULL_POLICY/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{IMAGE_PULL_SECRET}/$IMAGE_PULL_SECRET/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{SERVICE_TYPE}/$SERVICE_TYPE/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{SERVICE_PROTOCOL}/$SERVICE_PROTOCOL/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{CP_DEFAULT_INGRESS_CLASS_NAME}/$CP_DEFAULT_INGRESS_CLASS_NAME/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s@{VAULT_URL}@$VAULT_URL@g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{VAULT_INGRESS_HOST}/$(echo $VAULT_URL | awk -F[/:] '{print $4}')/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{VAULT_ROLE_NAME}/$VAULT_ROLE_NAME/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{VAULT_STORAGECLASS}/$K8S_STORAGECLASS/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{REPOSITORY_NAMESPACE}/${NAMESPACE[1]}/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s@{REPOSITORY_URL}@$REPOSITORY_URL@g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{REPOSITORY_HOST}/$REPOSITORY_HOST/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{REPOSITORY_PASSWORD}/$REPOSITORY_PASSWORD/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{REPOSITORY_PROJECT_NAME}/$REPOSITORY_PROJECT_NAME/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{REPOSITORY_STORAGECLASS}/$K8S_STORAGECLASS/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{DATABASE_NAMESPACE}/${NAMESPACE[2]}/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{DATABASE_URL}/$DATABASE_URL/g" {} \;
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
find ../values -name "*.yaml" -exec sed -i "s/{KEYCLOAK_LOG_LEVEL}/$KEYCLOAK_LOG_LEVEL/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{KEYCLOAK_CP_REALM}/$KEYCLOAK_CP_REALM/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{KEYCLOAK_CP_CLIENT_ID}/$KEYCLOAK_CP_CLIENT_ID/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{KEYCLOAK_CP_CLIENT_SECRET}/$KEYCLOAK_CP_CLIENT_SECRET/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{KEYCLOAK_INGRESS_TLS_ENABLED}/$KEYCLOAK_INGRESS_TLS_ENABLED/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{KEYCLOAK_INGRESS_HOST}/$(echo $KEYCLOAK_URL | awk -F[/:] '{print $4}')/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{KEYCLOAK_TLS_SECRET}/$KEYCLOAK_TLS_SECRET/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{NAMESPACE}/${NAMESPACE[4]}/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{HOST_CLUSTER_NAME}/$HOST_CLUSTER_NAME/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s@{CP_PORTAL_URL}@$CP_PORTAL_URL@g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{CP_PORTAL_INGRESS_HOST}/$(echo $CP_PORTAL_URL | awk -F[/:] '{print $4}')/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{CP_PORTAL_STORAGECLASS}/$K8S_STORAGECLASS/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{CP_SERVICE_PIPELINE_NAMESPACE}/$CP_SERVICE_PIPELINE_NAMESPACE/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{CP_SERVICE_SOURCE_CONTROL_NAMESPACE}/$CP_SERVICE_SOURCE_CONTROL_NAMESPACE/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{CP_DEFAULT_INGRESS_NAMESPACE}/$CP_DEFAULT_INGRESS_NAMESPACE/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{CP_CERT_SETUP_NAME}/${CHART_NAME[5]}/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{CP_CERT_SETUP_NAMESPACE}/$CP_CERT_SETUP_NAMESPACE/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{IPV6_ENABLED}/$IPV6_ENABLED/g" {} \;

CMD_CREATE_REGISTRY_SECRET="kubectl create secret docker-registry $IMAGE_PULL_SECRET --docker-server=$REPOSITORY_URL --docker-username=$REPOSITORY_USERNAME --docker-password=$REPOSITORY_PASSWORD"

# 1. deploy the vault
chmod +x ../vault/deploy-vault.sh
. ../vault/deploy-vault.sh
find ../values -name "*.yaml" -exec sed -i "s/{VAULT_ROLE_ID}/$VAULT_ROLE_ID/g" {} \;
find ../values -name "*.yaml" -exec sed -i "s/{VAULT_SECRET_ID}/$VAULT_SECRET_ID/g" {} \;

# 2. deploy the harbor
kubectl create namespace ${NAMESPACE[1]}
$CMD_CREATE_REGISTRY_SECRET -n ${NAMESPACE[1]}
helm install -f ../values/${CHART_NAME[1]}.yaml ${CHART_NAME[1]} ../charts/${CHART_NAME[1]}.tgz -n ${NAMESPACE[1]}

while :
do
  REPOSITORY_HTTP_CODE=$(curl -L -k -s -o /dev/null -w "%{http_code}\n" $REPOSITORY_URL/api/v2.0/projects)
  echo "[$REPOSITORY_HTTP_CODE] Please wait for several minutes for Harbor deployment to complete..."
  if [ $REPOSITORY_HTTP_CODE -eq 200 ]; then
    break
  fi
  sleep 10
done

# 3. create a project in harbor
REPOSITORY_TLS_SECRET=$(kubectl get secrets -n ${NAMESPACE[1]} --field-selector type=kubernetes.io/tls --no-headers -o custom-columns=":metadata.name" | grep ingress)
kubectl get secret $REPOSITORY_TLS_SECRET -n ${NAMESPACE[1]} -o jsonpath="{['data']['ca\.crt']}" | base64 --decode > cp-harbor-ca.crt
REPOSITORY_CERT=$(<./cp-harbor-ca.crt)
# update harbor cert to each node
echo "[Setup the harbor cert in cluster]..."
helm install -f ../values/${CHART_NAME[5]}.yaml ${CHART_NAME[5]} ../charts/${CHART_NAME[5]}.tgz -n $CP_CERT_SETUP_NAMESPACE --set-literal data.target.cert="$REPOSITORY_CERT"

while :
do
  POD_COUNT=$((kubectl get pods -n $CP_CERT_SETUP_NAMESPACE -l $CP_CERT_SETUP_SELECTOR --field-selector status.phase!=Running --no-headers | wc -l) 2> /dev/null)
  echo "[remaining: $POD_COUNT] Adding a harbor cert to each node's container runtime..."
  if [[ $POD_COUNT -lt 1 ]]; then
    echo "Completed..."
    break
  fi
  sleep 5
done

sudo mv cp-harbor-ca.crt /usr/local/share/ca-certificates/
sudo update-ca-certificates

curl -u $REPOSITORY_USERNAME:$REPOSITORY_PASSWORD -k $REPOSITORY_URL/api/v2.0/projects -XPOST --data-binary "{\"project_name\": \"$REPOSITORY_PROJECT_NAME\", \"public\": false}" -H "Content-Type: application/json" -i
sudo podman login $REPOSITORY_HOST --username $REPOSITORY_USERNAME --password $REPOSITORY_PASSWORD
helm repo add $REPOSITORY_PROJECT_NAME --username $REPOSITORY_USERNAME --password $REPOSITORY_PASSWORD $REPOSITORY_URL/chartrepo/$REPOSITORY_PROJECT_NAME
helm plugin install https://github.com/chartmuseum/helm-push.git

# 4. push the keycloak image in harbor
cd ../keycloak
sudo podman build --tag $REPOSITORY_HOST/$REPOSITORY_PROJECT_NAME/${CHART_NAME[3]}:latest .
sudo podman push $REPOSITORY_HOST/$REPOSITORY_PROJECT_NAME/${CHART_NAME[3]}
cd ../script

# 5. push the container platform portal image in harbor
for IMAGE in ${IMAGE_NAME[@]}
do
    sudo podman load -i ../images/$IMAGE.tar.gz
    sudo podman tag localhost:5000/container-platform/$IMAGE $REPOSITORY_HOST/$REPOSITORY_PROJECT_NAME/$IMAGE
    sudo podman push $REPOSITORY_HOST/$REPOSITORY_PROJECT_NAME/$IMAGE
done

# 6. push the chart in harbor
for CHART in ${CHART_NAME[@]}
do
  helm cm-push --username $REPOSITORY_USERNAME --password $REPOSITORY_PASSWORD ../charts/$CHART.tgz $REPOSITORY_PROJECT_NAME
done
helm repo update $REPOSITORY_PROJECT_NAME

# 7. deploy the mariadb
kubectl create namespace ${NAMESPACE[2]}
$CMD_CREATE_REGISTRY_SECRET -n ${NAMESPACE[2]}
kubectl apply -f ../values/${CHART_NAME[2]}-configmap.yaml -n ${NAMESPACE[2]}
helm install -f ../values/${CHART_NAME[2]}.yaml ${CHART_NAME[2]} $REPOSITORY_PROJECT_NAME/mariadb -n ${NAMESPACE[2]}

# 8. deploy the keycloak
kubectl create namespace ${NAMESPACE[3]}
$CMD_CREATE_REGISTRY_SECRET -n ${NAMESPACE[3]}
if [ $KEYCLOAK_INGRESS_TLS_ENABLED = "true" ]; then
  kubectl create secret tls $KEYCLOAK_TLS_SECRET --key $KEYCLOAK_TLS_KEY_PATH --cert $KEYCLOAK_TLS_CERT_PATH -n ${NAMESPACE[3]}
fi
helm install -f ../values/${CHART_NAME[3]}.yaml ${CHART_NAME[3]} $REPOSITORY_PROJECT_NAME/keycloak -n ${NAMESPACE[3]}

# 9. deploy the container platform portal
kubectl create namespace ${NAMESPACE[4]}
$CMD_CREATE_REGISTRY_SECRET -n ${NAMESPACE[4]}
helm install -f ../values/${RELEASE_NAME}.yaml ${RELEASE_NAME} $REPOSITORY_PROJECT_NAME/${CHART_NAME[4]} -n ${NAMESPACE[4]}

# uninstall cp-cert-setup
helm uninstall ${CHART_NAME[5]} -n $CP_CERT_SETUP_NAMESPACE
