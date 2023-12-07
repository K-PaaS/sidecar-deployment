#!/bin/bash
source variables.yml

################################################################################################################
#################################     inject-self-signed-cert 활용 가이드      #################################
################################################################################################################
# Private Repository (e.g. Harbor, distribution....) 사용 시 HTTPS가 지원되야 K-PaaS Sidecar 설치가 가능하다.
# HTTPS 설정 시 자체 서명된 인증서를 사용하려면 Private Repository를 사용하는 POD에 HTTPS 설정 시 사용한 인증서를 저장해야 한다.
# deploy-inject-self-signed-cert.sh에서 사용되는 cert-injection-webhook(https://github.com/vmware-tanzu/cert-injection-webhook)는 
# labels 또는 annotations을 설정하여 해당 labels 또는 annotations을 가진 POD가 배포될 시 컨테이너 내부에 인증서를 삽입한다.
################################################################################################################
# 사용 방법
# 1. Kubernetes의 사용할 Registry에 대한 insecure-registry 설정을 진행한다.
# 2. 이하의 과정은 Sidecar를 이용하여 Application을 배포하기 전 어느 시점에 해도 무방하다. 
# 3. variables.yml 설정을 진행한다.
# 3-1. is_self_signed_certificate=true
# 3-2. app_registry_cert_path=support-files/private-repository.ca
# 4. app_registry_cert_path에 위치한 파일에 Private Repository에 사용된 인증서 CA를 넣는다.
# 5. deploy-inject-self-signed-cert.sh를 실행한다.
################################################################################################################



if [[ ${use_dockerhub} = false ]] && [[ ${is_self_signed_certificate} = "true" ]]; then
  if [[ -e $registry_cert_path ]]; then
    echo "registry_cert_path file exists"
  else
    echo "plz check registry_cert_path"
    return
  fi
else
  echo "plz check variable registry_kind or is_self_signed_certificate"
  return
fi

ytt -f ./support-files/cert-injection-webhook-config \
      --data-value-file ca_cert_data=${registry_cert_path} \
      --data-value-yaml labels="[kpack.io/build, private-repo-cert-injection]"  > support-files/cert-injection-webhook/manifest.yaml


kapp deploy -a cert-injection-webhook -f ./support-files/cert-injection-webhook/manifest.yaml -y

