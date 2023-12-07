#!/bin/bash
source variables.yml

# namespace create
for ns in $sidecar_namespace $root_namespace; do
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Namespace
metadata:
  labels:
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/enforce: restricted
  name: $ns
EOF
done

# registry secret
if [[ ${use_dockerhub} = true ]]; then
	kubectl create secret -n $root_namespace docker-registry image-registry-credentials \
		--docker-username=$registry_id \
		--docker-password=$registry_password
elif [[ ${use_dockerhub} = false ]]; then
	kubectl create secret -n $root_namespace docker-registry image-registry-credentials \
		--docker-username=$registry_id \
		--docker-password=$registry_password \
		--docker-server=$registry_address
else
  echo "plz check variable use_dockerhub"
  return
fi
