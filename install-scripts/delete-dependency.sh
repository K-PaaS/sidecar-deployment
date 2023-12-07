#!/bin/bash
source variables.yml

# delete dependency

## private_registry ca secret (option)

if [[ ${use_dockerhub} = true ]]; then
	echo "" > /dev/null
elif [[ ${use_dockerhub} = false ]]; then
	if [[ ${is_self_signed_certificate} = true ]]; then
		## kpack ca secret
		kubectl --namespace kpack delete secret $cert_secret_name
		## core ca secret
		kubectl --namespace $sidecar_namespace delete secret $cert_secret_name
	fi
else
  echo "plz check variable use_dockerhub"
  return
fi

## service-binding
kubectl delete -f "../dependency/service-binding/servicebinding-workloadresourcemappings-v"*".yaml"
kubectl delete -f "../dependency/service-binding/servicebinding-runtime-v"*".yaml"

## cert-manager
kubectl delete -f ../dependency/cert-manager/cert-manager.yaml

## contour
kubectl delete -f ../dependency/contour/

## kpack
kubectl delete -f ../dependency/kpack/


# delete registry secret
kubectl delete secret image-registry-credentials -n $sidecar_namespace

# delete namespace 
kubectl delete namespace $sidecar_namespace
kubectl delete namespace $root_namespace
