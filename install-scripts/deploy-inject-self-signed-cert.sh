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
# 1. insecure-registry 설정을 진행한다. (CRI-O 기준 설정 가이드 : https://github.com/K-PaaS/container-platform/blob/master/install-guide/bosh/cp-bosh-deployment-spray-guide-v1.1.md#3.1)
# 2. variables.yml 설정을 진행한다.
# 2-1. is_self_signed_certificate=true
# 2-2. app_registry_cert_path=support-files/private-repository.ca
# 3. app_registry_cert_path에 위치한 파일에 Private Repository에 사용된 인증서 CA를 넣는다.
# 4. deploy-inject-self-signed-cert.sh를 실행한다.
# 5. Sidecar 설치 과정에 따라 2.generate-values.sh, 3.rendering-values.sh 스크립트를 실행한다.
# 6. 4.deploy-sidecar.sh 스크립트를 실행하여 Sidecar를 설치한다.
################################################################################################################



if [[ ${app_registry_kind} = "private" ]] && [[ ${is_self_signed_certificate} = "true" ]]; then
  if [[ -e $app_registry_cert_path ]]; then
    echo "app_registry_cert_path file exists"
  else
    echo "plz check app_registry_cert_path"
    return
  fi
else
  echo "plz check variable app_registry_kind or is_self_signed_certificate"
  return
fi

ytt -f ./support-files/cert-injection-webhook-config \
      --data-value-file ca_cert_data=${app_registry_cert_path} \
      --data-value-yaml labels="[kpack.io/build, private-repo-cert-injection]"  > support-files/cert-injection-webhook/manifest.yaml


kapp deploy -a cert-injection-webhook -f ./support-files/cert-injection-webhook/manifest.yaml -y

