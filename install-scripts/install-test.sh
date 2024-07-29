#!/bin/bash
TEST_ORG=temp-test-org
TEST_SPACE=temp-test-space
TEST_APP_NAME=temp-test-app

source variables.yml

export KUBECONFIG=$(pwd)/support-files/user/sidecar-$admin_username.ua.kubeconfig

cf login -a api.$(cat variables.yml | grep "system_domain" | cut -d '=' -f 2 | cut -d ' ' -f 1 ) \
-u $(cat variables.yml | grep "admin_username" | cut -d '=' -f 2 | cut -d ' ' -f 1 ) \
--skip-ssl-validation << EOF
  
EOF


cf create-org $TEST_ORG
cf create-space -o $TEST_ORG $TEST_SPACE
cf target -o $TEST_ORG -s $TEST_SPACE
cf push $TEST_APP_NAME -p support-files/sample-app

echo "=============================="
echo "check output 'Hello World'"
curl https://$(cf app $TEST_APP_NAME | grep "routes" | sed -e 's/ //g' | cut -d ':' -f 2) -k
echo "=============================="
cf delete $TEST_APP_NAME -f
cf delete-org $TEST_ORG -f
