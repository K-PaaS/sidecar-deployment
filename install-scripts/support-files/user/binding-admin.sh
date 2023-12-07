#!/bin/bash

source ../../variables.yml

set -e 

main() {

  if [[ $# -ne 2 ]]; then
    cat <<EOF >&2
Usage:
  $(basename "$0") <account_kind (ua, sa)> <username>

EOF
    exit 1
  fi
  if [[ $1 -ne ua ]] || [[ $1 -ne sa ]]; then
    cat <<EOF >&2
Usage:
  $(basename "$0") <account_kind (ua, sa)> <username>

EOF
    exit 1
  fi

  sa_namespace=$root_namespace
  if [[ $1 == sa ]]; then
    echo $1
    if [[ -z $(kubectl get sa -n $sa_namespace $2) ]]; then
      echo "plz username name"
      exit 1
    fi
  fi


  if [[ $1 == ua ]]; then
    rolebinding_name="admin-${1}-${2}" yq e '.metadata.name = env(rolebinding_name)' rolebinding-template/rolebinding-admin.yaml | username="${2}" yq e '.subjects[0].name = env(username)'  | yq e '.subjects[0].kind = "User"'  | yq e 'del(.subjects[0].namespace)' | kubectl apply -f -
  elif [[ $1 == sa ]]; then
    rolebinding_name="admin-${1}-${2}" yq e  '.metadata.name = env(rolebinding_name)' rolebinding-template/rolebinding-admin.yaml | username="${2}" yq e '.subjects[0].name = env(username)' | sa_namespace="${sa_namespace}" yq e '.subjects[0].namespace = env(sa_namespace)' | kubectl apply -f -
  fi
}

main $@

