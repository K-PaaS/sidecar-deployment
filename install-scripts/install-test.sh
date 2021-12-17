#!/bin/bash
TEST_SPACE=temp-test-space
TEST_APP_NAME=temp-test-app

cf login -a api.$(cat manifest/sidecar-values.yml | grep "system_domain" | cut -d ' ' -f 2 | sed -e 's/\"//g') \
-u admin \
-p $(cat manifest/sidecar-values.yml | grep "cf_admin_password" | cut -d ' ' -f 2) \
-o system \
--skip-ssl-validation << EOF
  

EOF

cf create-space $TEST_SPACE
cf target -s $TEST_SPACE
cf push $TEST_APP_NAME -p ../tests/smoke/assets/test-node-app

echo "=============================="
echo "check output 'Hello World'"
curl https://$(cf app $TEST_APP_NAME | grep "routes" | sed -e 's/ //g' | cut -d ':' -f 2) -k
echo "=============================="
cf delete $TEST_APP_NAME -f
cf delete-space $TEST_SPACE -f
