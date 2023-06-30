#!/bin/bash

source logging-service-variable.yml
source infra/infra-variable.yml

CONFIGMAP_INPUT_PATH="support-files/fluentd-plugin-${LOGGING_OUTPUT_PLUGIN}.yml"
CONFIGMAP_ORIGIN_PATH="manifest/fluentd-configmap-origin.yml"


## Modify logging input data for fluentd ConfigMap
if [[ ${LOGGING_OUTPUT_PLUGIN} = "influxdb" ]]; then
        sed -i "s/<INFLUXDB_IP>/$INFLUXDB_IP/" ${CONFIGMAP_INPUT_PATH}
        sed -i "s/<INFLUXDB_HTTP_PORT>/$INFLUXDB_HTTP_PORT/" ${CONFIGMAP_INPUT_PATH}
        sed -i "s/<INFLUXDB_HTTPS_ENABLED>/$INFLUXDB_HTTPS_ENABLED/" ${CONFIGMAP_INPUT_PATH}
else
        INFLUXDB_URL=""

        if [[ ${INFLUXDB_HTTPS_ENABLED} = "true" ]]; then
                INFLUXDB_URL="https:\/\/"
                OPT_LINE_NUM=$(grep -n "http_method" ${CONFIGMAP_INPUT_PATH} | cut -d: -f1)
                sed -i "${OPT_LINE_NUM} i\          tls_verify_mode none" ${CONFIGMAP_INPUT_PATH}
        else
                INFLUXDB_URL="http:\/\/"
        fi

        INFLUXDB_URL+=${INFLUXDB_IP}:${INFLUXDB_HTTP_PORT}
        sed -i "s/<INFLUXDB_URL>/$INFLUXDB_URL/" ${CONFIGMAP_INPUT_PATH}
fi

sed -i "s/<INFLUXDB_MEASUREMENT>/$INFLUXDB_MEASUREMENT/" ${CONFIGMAP_INPUT_PATH}
sed -i "s/<INFLUXDB_DATABASE>/$INFLUXDB_DATABASE/" ${CONFIGMAP_INPUT_PATH}
sed -i "s/<INFLUXDB_USERNAME>/$INFLUXDB_USERNAME/" ${CONFIGMAP_INPUT_PATH}
sed -i "s/<INFLUXDB_PASSWORD>/$INFLUXDB_PASSWORD/" ${CONFIGMAP_INPUT_PATH}
sed -i "s/<INFLUXDB_TIME_PRECISION>/$INFLUXDB_TIME_PRECISION/" ${CONFIGMAP_INPUT_PATH}


## Modify fluentd-config-ver-1
mkdir manifest -p

kubectl get cm fluentd-config-ver-1 -n ${FLUENTD_NAMESPACE} -o yaml > ${CONFIGMAP_ORIGIN_PATH}

CONFIGMAP_HEAD_LINE_NUM=$(grep -n "<match \*\*>" ${CONFIGMAP_ORIGIN_PATH} | cut -d: -f1)
CONFIGMAP_HEAD_DATA=$(sed -n "1,`expr ${CONFIGMAP_HEAD_LINE_NUM} "-" "1"`p" ${CONFIGMAP_ORIGIN_PATH})
CONFIGMAP_LOGGING_DATA=$(cat ${CONFIGMAP_INPUT_PATH})
CONFIGMAP_TAIL_LINE_NUM=$(grep -n "kind: ConfigMap" ${CONFIGMAP_ORIGIN_PATH} | cut -d: -f1)
CONFIGMAP_TAIL_DATA=$(sed -n "${CONFIGMAP_TAIL_LINE_NUM},\$p" ${CONFIGMAP_ORIGIN_PATH})
CONFIGMAP_NEW_PATH="manifest/fluentd-configmap-logging.yml"

cat << EOF >> ${CONFIGMAP_NEW_PATH}
${CONFIGMAP_HEAD_DATA}
${CONFIGMAP_LOGGING_DATA}
${CONFIGMAP_TAIL_DATA}
EOF

rm ${CONFIGMAP_ORIGIN_PATH}

kubectl replace -f ${CONFIGMAP_NEW_PATH}
kubectl rollout restart ds/fluentd -n ${FLUENTD_NAMESPACE}
