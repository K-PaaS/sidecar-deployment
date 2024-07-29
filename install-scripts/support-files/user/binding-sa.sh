#!/bin/bash

source ../../variables.yml

main() {
  if [[ $# -ne 4 ]]; then
    cat <<EOF >&2
Usage:
  $(basename "$0") <namespace> <username> <org_name> <space_name>

EOF
    exit 1
  fi
  sa_namespace=$1
  username="$2"
  org_name="$3"
  space_name="$4"

  real_org_name=$(kubectl get cforgs -n $root_namespace | grep " $org_name " | cut -d ' ' -f 1)
  if [[ -z $real_org_name ]]; then
    echo "plz org name"
    exit 1
  fi

  real_space_name=$(kubectl get cfspaces -n $real_org_name | grep " $space_name " | cut -d ' ' -f 1)
  if [[ -z $real_space_name  ]]; then
    echo "plz space name"
    exit 1
  fi

  if [[ -z $(kubectl get sa -n $sa_namespace $username) ]]; then
    echo "plz username name"
    exit 1
  fi

  rolebinding_name="${real_org_name}-${username}" yq e '.metadata.name = env(rolebinding_name)' rolebinding-template/rolebinding-sa-orguser.yaml | real_org_name="${real_org_name}" yq e '.metadata.namespace = env(real_org_name)' |
username="${username}" yq e '.subjects[0].name = env(username)' | sa_namespace="${sa_namespace}" yq e '.subjects[0].namespace = env(sa_namespace)' | kubectl apply -f -

  rolebinding_name="${real_space_name}-${username}" yq e '.metadata.name = env(rolebinding_name)' rolebinding-template/rolebinding-sa-spacedeveloper.yaml |  real_space_name="${real_space_name}" yq e '.metadata.namespace = env(real_space_name)' | username="${username}" yq e '.subjects[0].name = env(username)' | sa_namespace="${sa_namespace}" yq e '.subjects[0].namespace = env(sa_namespace)' | kubectl apply -f -

  rolebinding_name="sidecar-root-namespace-user-${sa_namespace}-${username}" yq e '.metadata.name = env(rolebinding_name)' rolebinding-template/rolebinding-root-namespace-user.yaml |  root_namespace="${root_namespace}" yq e '.metadata.namespace = env(root_namespace)' | username="${username}" yq e '.subjects[0].name = env(username)' | sa_namespace="${sa_namespace}" yq e '.subjects[0].namespace = env(sa_namespace)' | kubectl apply -f -

}

main $@
