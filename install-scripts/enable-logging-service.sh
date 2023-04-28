#!/bin/bash

source variables.yml

## Modify logging input data for fluentd ConfigMap
workspace=$(pwd)
fluentd_logging_input_path="${workspace}/support-files/logging-configmap-plugin-${logging_output_plugin}.yml"
fluentd_namespace=cf-system

if [[ ${logging_output_plugin} = "influxdb" ]]; then
        sed -i "s/<influxdb_ip>/$influxdb_ip/" ${fluentd_logging_input_path}
        sed -i "s/<influxdb_http_port>/$influxdb_http_port/" ${fluentd_logging_input_path}
        sed -i "s/<influxdb_https_enabled>/$influxdb_https_enabled/" ${fluentd_logging_input_path}
else
        influxdb_url=""

        if [[ ${influxdb_https_enabled} = "true" ]]; then
                influxdb_url="https:\/\/"
                opt_line_num=$(grep -n "http_method" ${fluentd_logging_input_path} | cut -d: -f1)
                sed -i "${opt_line_num} i\          tls_verify_mode none" ${fluentd_logging_input_path}
        else
                influxdb_url="http:\/\/"
        fi

        influxdb_url+=${influxdb_ip}:${influxdb_http_port}
        sed -i "s/<influxdb_url>/$influxdb_url/" ${fluentd_logging_input_path}
fi

sed -i "s/<influxdb_measurement>/$influxdb_measurement/" ${fluentd_logging_input_path}
sed -i "s/<influxdb_database>/$influxdb_database/" ${fluentd_logging_input_path}
sed -i "s/<influxdb_username>/$influxdb_username/" ${fluentd_logging_input_path}
sed -i "s/<influxdb_password>/$influxdb_password/" ${fluentd_logging_input_path}
sed -i "s/<influxdb_time_precision>/$influxdb_time_precision/" ${fluentd_logging_input_path}


## Modify fluentd-config-ver-1
fluentd_cm_origin_path="${workspace}/support-files/fluentd-configmap-origin.yml"

kubectl get cm fluentd-config-ver-1 -n ${fluentd_namespace} -o yaml > ${fluentd_cm_origin_path}

fluentd_cm_head_line_num=$(grep -n "<match \*\*>" ${fluentd_cm_origin_path} | cut -d: -f1)
fluentd_cm_head_data=$(sed -n "1,`expr ${fluentd_cm_head_line_num} "-" "1"`p" ${fluentd_cm_origin_path})
fluentd_cm_logging_data=$(cat ${fluentd_logging_input_path})
fluentd_cm_tail_line_num=$(grep -n "kind: ConfigMap" ${fluentd_cm_origin_path} | cut -d: -f1)
fluentd_cm_tail_data=$(sed -n "${fluentd_cm_tail_line_num},\$p" ${fluentd_cm_origin_path})
fluentd_cm_new_path="${workspace}/manifest/fluentd-configmap-new.yml"

cat << EOF >> ${fluentd_cm_new_path}
${fluentd_cm_head_data}
${fluentd_cm_logging_data}
${fluentd_cm_tail_data}
EOF

rm ${fluentd_cm_origin_path}

kubectl replace -f ${fluentd_cm_new_path}
kubectl rollout restart ds/fluentd -n ${fluentd_namespace}