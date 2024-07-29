#!/bin/bash

set -e 

source ../../variables.yml

# create kubeconfig
server=$(kubectl config view -o jsonpath='{.clusters[0].cluster.server}')
username=$1
namespace=$root_namespace
token_name=sidecar-$username-sa-token

  if [[ $# -ne 1 ]]; then
    cat <<EOF >&2
Usage:
  $(basename "$0") <username>

EOF
    exit 1
  fi

# service account add
kubectl create sa $username -n $root_namespace

# create token
cat << EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: $token_name
  namespace: $namespace
  annotations:
    kubernetes.io/service-account.name: $username
type: kubernetes.io/service-account-token
EOF

# create kubeconfig
kubectl -n $namespace get secret/$token_name -o jsonpath='{.data.token}' | base64 --decode

ca=$(kubectl -n $namespace get secret/$token_name -o jsonpath='{.data.ca\.crt}')
token=$(kubectl -n $namespace get secret/$token_name -o jsonpath='{.data.token}' | base64 --decode)
namespace=$(kubectl -n $namespace get secret/$token_name -o jsonpath='{.data.namespace}' | base64 --decode)

echo "
apiVersion: v1
kind: Config
clusters:
- name: sidecar-cluster
  cluster:
    certificate-authority-data: ${ca}
    server: ${server}
contexts:
- name: sidecar-context
  context:
    cluster: sidecar-cluster
    namespace: ${namespace}
    user: ${username}
current-context: sidecar-context
users:
- name: ${username}
  user:
    token: ${token}
" > sidecar-$username.sa.kubeconfig
