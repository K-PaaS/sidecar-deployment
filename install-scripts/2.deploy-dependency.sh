#!/bin/bash
source variables.yml

# deploy dependency

## cert-manager
kubectl apply -f ../dependency/cert-manager/cert-manager.yaml

kubectl -n cert-manager rollout status deployment/cert-manager --watch=true
kubectl -n cert-manager rollout status deployment/cert-manager-webhook --watch=true
kubectl -n cert-manager rollout status deployment/cert-manager-cainjector --watch=true

cmctl check api --wait=2m

## contour
kubectl apply -f ../dependency/contour/contour-gateway-provisioner.yaml

gateway_name="${sidecar_namespace}-gateway" yq e  '(select(.kind == "GatewayClass")| .metadata.name = env(gateway_name))' ../dependency/contour/gateway-class.yaml | kubectl apply -f -


if [[ ${use_lb} == true ]]; then
	yq e '(select(.kind == "ContourDeployment") | .spec.envoy.networkPublishing.type = "LoadBalancerService" )' ../dependency/contour/contour-deployment.yaml | kubectl apply -f -

elif [[ ${use_lb} == false ]]; then
	yq e  '(select(.kind == "ContourDeployment")| .spec.envoy.networkPublishing.type = "NodePortService" )' ../dependency/contour/contour-deployment.yaml | kubectl apply -f -

else
	echo "plz check variable use_lb"
	return
fi


## kpack
kubectl apply -f ../dependency/kpack/

if [[ ${use_dockerhub} = true ]]; then
	echo "" > /dev/null
elif [[ ${use_dockerhub} = false ]]; then
	if [[ ${is_self_signed_certificate} == true ]]; then
		echo "is_self_signed_certificate true"
		cert_secret_name="${cert_secret_name}" yq e  '(select(.kind == "Deployment" and .metadata.name =="kpack-controller" )| .spec.template.spec += {"volumes": [{"name": "korifi-registry-ca-cert", "secret": {"defaultMode": 420, "secretName": env(cert_secret_name)}}] })' ../dependency/kpack/release-0.13.4.yaml | yq e  '(select(.kind == "Deployment" and .metadata.name =="kpack-controller" )| .spec.template.spec.containers[0] += {"volumeMounts": [{"mountPath": "/etc/ssl/certs/registry-ca.crt", "name": "korifi-registry-ca-cert", "readOnly": true, "subPath": "ca.crt"}]})' | kubectl apply -f -
		## private_registry ca secret (option)
		### kpack ca secret
		kubectl --namespace kpack create secret generic $cert_secret_name \
			--from-file=ca.crt=$registry_cert_path
		### core ca secret
		kubectl --namespace $sidecar_namespace create secret generic $cert_secret_name \
			--from-file=ca.crt=$registry_cert_path
	elif [[ ${is_self_signed_certificate} == false ]]; then
		echo "is_self_signed_certificate false"
	else
		echo "plz check variable is_self_signed_certificate"
		return
	fi
else
  echo "plz check variable use_dockerhub"
  return
fi

## service-binding
kubectl apply -f "../dependency/service-binding/servicebinding-runtime-v"*".yaml"
kubectl -n servicebinding-system rollout status deployment/servicebinding-controller-manager --watch=true
kubectl apply -f "../dependency/service-binding/servicebinding-workloadresourcemappings-v"*".yaml"

echo ""
echo "====================cert-manager===================="
echo ""
kubectl get all -n cert-manager
echo ""

echo ""
echo "======================contour======================="
echo ""
kubectl get all -n projectcontour
echo ""

echo ""
echo "======================kpack========================="
echo ""
kubectl get all -n kpack
echo ""

echo ""
echo "===================service-binding=================="
echo ""
kubectl get all -n servicebinding-system
echo ""
