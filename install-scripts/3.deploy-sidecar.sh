#!/bin/bash
source variables.yml

HELM_VERSION=0.12.0

# deploy sidecar

if [[ ${use_dockerhub} = true ]]; then
	registry_address=index.docker.io
	registry_repositry_name=$registry_id
elif [[ ${use_dockerhub} = false ]]; then
	echo "" > /dev/null
else
  echo "plz check variable use_dockerhub"
  return
fi

helm dependency update ../helm/korifi-$HELM_VERSION


if [[ ${use_dockerhub} = false ]] && [[ ${is_self_signed_certificate} = true ]] ; then
        helm upgrade --install sidecar ../helm/korifi-$HELM_VERSION \
                --namespace="$sidecar_namespace" \
                --set=generateIngressCertificates=true \
                --set=rootNamespace="$root_namespace" \
                --set=adminUserName="$admin_username" \
                --set=api.apiServer.url="api.$system_domain" \
                --set=defaultAppDomainName="apps.$system_domain" \
                --set=containerRepositoryPrefix=$registry_address/$registry_repositry_name/ \
                --set=containerRegistryCACertSecret="$cert_secret_name" \
                --set=kpackImageBuilder.builderRepository=$registry_address/$registry_repositry_name/kpack-builder \
                --set=api.userCertificateExpirationWarningDuration=$(($user_certificate_expiration_duration_days*24))"h" \
                --set=networking.gatewayClass="$sidecar_namespace-gateway" \
                --wait
elif [[ ${use_dockerhub} = true ]] || [[ ${use_dockerhub} = false && ${is_self_signed_certificate} = false ]] ; then
        helm upgrade --install sidecar ../helm/korifi-$HELM_VERSION \
                --namespace="$sidecar_namespace" \
                --set=generateIngressCertificates=true \
                --set=rootNamespace="$root_namespace" \
                --set=adminUserName="$admin_username" \
                --set=api.apiServer.url="api.$system_domain" \
                --set=defaultAppDomainName="apps.$system_domain" \
                --set=containerRepositoryPrefix=$registry_address/$registry_repositry_name/ \
                --set=kpackImageBuilder.builderRepository=$registry_address/$registry_repositry_name/kpack-builder \
                --set=api.userCertificateExpirationWarningDuration=$(($user_certificate_expiration_duration_days*24))"h" \
                --set=networking.gatewayClass="$sidecar_namespace-gateway" \
                --wait
else
	echo "plz check variable is_self_signed_certificate or use_dockerhub"
	return
fi


echo "create admin"
pushd "support-files/user" >/dev/null
{
  ./create-new-ua.sh $admin_username
}
popd >/dev/null

kubectl get all -n $sidecar_namespace
