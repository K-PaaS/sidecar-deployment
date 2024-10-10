#!/bin/bash

source ../../variables.yml

set -e 

randomGuidGenerate(){
        randomGuid=$(tr -dc a-f0-9 </dev/urandom | head -c 8 )-$(tr -dc a-f0-9 </dev/urandom | head -c 4 )-$(tr -dc a-f0-9 </dev/urandom | head -c 4 )-$(tr -dc a-f0-9 </dev/urandom | head -c 4 )-$(tr -dc a-f0-9 </dev/urandom | head -c 12 )
        echo ${randomGuid}
}

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

  if [[ $1 == sa ]]; then
    echo $1
    echo -n "plz input service acount namespace : "
    read sa_namespace
    if [[ -z $(kubectl get sa -n $sa_namespace $2) ]]; then
      echo "plz check username & namespace"
      exit 1
    fi
  fi



  if [[ $1 == ua ]]; then
    rolebinding_name="admin-${1}-${2}" yq e '.metadata.name = env(rolebinding_name)' rolebinding-template/rolebinding-admin.yaml | username="${2}" yq e '.subjects[0].name = env(username)'  | yq e '.subjects[0].kind = "User"'  | yq e 'del(.subjects[0].namespace)' | randomGuid=$(randomGuidGenerate) yq e '.metadata.labels."cloudfoundry.org/role-guid" = env(randomGuid)' | kubectl apply -f -
  elif [[ $1 == sa ]]; then
    rolebinding_name="admin-${1}-${2}" yq e  '.metadata.name = env(rolebinding_name)' rolebinding-template/rolebinding-admin.yaml | username="${2}" yq e '.subjects[0].name = env(username)' | sa_namespace="${sa_namespace}" yq e '.subjects[0].namespace = env(sa_namespace)' | randomGuid=$(randomGuidGenerate) yq e '.metadata.labels."cloudfoundry.org/role-guid" = env(randomGuid)' | kubectl apply -f -
  fi
}

main $@

