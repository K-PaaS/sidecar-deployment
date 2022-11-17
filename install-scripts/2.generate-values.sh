#!/bin/bash
source variables.yml

mkdir manifest -p

../hack/generate-values.sh -d ${system_domain} > ./manifest/sidecar-values.yml

cat << EOF >> ./manifest/sidecar-values.yml
use_first_party_jwt_tokens: true
enable_automount_service_account_token: true
EOF


## LoadBalancer Setting
if [[ ${use_lb} = "true" ]]; then
cat << EOF >> ./manifest/sidecar-values.yml
load_balancer:
   enable: true
EOF

elif [[ ${use_lb} = "false" ]]; then
cat << EOF >> ./manifest/sidecar-values.yml
load_balancer:
   enable: false
EOF
else
        echo "plz check variables.yml : use_lb"
        return
fi

if [[ ${iaas} = "openstack" ]] && [[ ${use_lb} = "true" ]]; then
cat << EOF >> ./manifest/sidecar-values.yml
   static_ip: ${public_ip}
EOF
fi


## App Registry Setting
if [[ ${app_registry_kind} = "dockerhub" ]]; then
cat << EOF >> ./manifest/sidecar-values.yml
app_registry:
  hostname: https://index.docker.io/v1/
  repository_prefix: "${app_registry_repository}"
  username: "${app_registry_id}"
  password: "${app_registry_password}"
EOF

elif [[ ${app_registry_kind} = "private" ]]; then
cat << EOF >> ./manifest/sidecar-values.yml
app_registry:
  hostname: https://${app_registry_address}/v2/
  repository_prefix: "${app_registry_address}/${app_registry_repository}"
  username: "${app_registry_id}"
  password: "${app_registry_password}"
EOF
else
        echo "plz check variables.yml : app_registry_kind"
        return
fi

## App Private Registry Setting
if [[ ${is_self_signed_certificate} = "true" ]]; then
cat << EOF >> ./manifest/sidecar-values.yml
is_self_signed_certificate: true
EOF
elif [[ ${is_self_signed_certificate} = "false" ]]; then
cat << EOF >> ./manifest/sidecar-values.yml
is_self_signed_certificate: false
EOF
else
        echo "plz check variables.yml : is_self_signed_certificate"
        return
fi

## Portal Web User Setting
cat << EOF >> ./manifest/sidecar-values.yml
webuser_name: ${webuser_name}
EOF

## External Blobstore Setting
if [[ ${use_external_blobstore} = "true" ]]; then
  cp support-files/external-blobstore-values.yml manifest/ -f
  sed -i "s/<external_blobstore_ip>/$external_blobstore_ip/" manifest/external-blobstore-values.yml
  sed -i "s/<external_blobstore_port>/$external_blobstore_port/" manifest/external-blobstore-values.yml
  sed -i "s/<external_blobstore_id>/$external_blobstore_id/" manifest/external-blobstore-values.yml
  sed -i "s/<external_blobstore_password>/$external_blobstore_password/" manifest/external-blobstore-values.yml
  sed -i "s/<external_blobstore_package_directory>/$external_blobstore_package_directory/" manifest/external-blobstore-values.yml
  sed -i "s/<external_blobstore_droplet_directory>/$external_blobstore_droplet_directory/" manifest/external-blobstore-values.yml
  sed -i "s/<external_blobstore_resource_directory>/$external_blobstore_resource_directory/" manifest/external-blobstore-values.yml
  sed -i "s/<external_blobstore_buildpack_directory>/$external_blobstore_buildpack_directory/" manifest/external-blobstore-values.yml
fi


## External Database Setting
if [[ ${use_external_db} = "true" ]]; then
  if [[ ${external_db_kind} = "postgres" ]]; then
    cp support-files/external-db-values-postgresql.yml manifest/ -f
    sed -i "s/<external_db_ip>/$external_db_ip/" manifest/external-db-values-postgresql.yml
    sed -i "s/<external_db_port>/$external_db_port/" manifest/external-db-values-postgresql.yml
    sed -i "s/<external_cc_db_id>/$external_cc_db_id/" manifest/external-db-values-postgresql.yml
    sed -i "s/<external_cc_db_password>/$external_cc_db_password/" manifest/external-db-values-postgresql.yml
    sed -i "s/<external_cc_db_name>/$external_cc_db_name/" manifest/external-db-values-postgresql.yml
    sed -i "s/<external_uaa_db_id>/$external_uaa_db_id/" manifest/external-db-values-postgresql.yml
    sed -i "s/<external_uaa_db_password>/$external_uaa_db_password/" manifest/external-db-values-postgresql.yml
    sed -i "s/<external_uaa_db_name>/$external_uaa_db_name/" manifest/external-db-values-postgresql.yml
    render_ca=$(awk '{print "      "$0}' $external_db_cert_path)
    awk -v ca_values="$render_ca" 'NR==24{print ca_values}1' manifest/external-db-values-postgresql.yml > manifest/external-db-values-postgresql-temp.yml
    awk -v ca_values="$render_ca" 'NR==13{print ca_values}1' manifest/external-db-values-postgresql-temp.yml > manifest/external-db-values-postgresql.yml
    rm manifest/external-db-values-postgresql-temp.yml
  elif [[ ${external_db_kind} = "mysql" ]]; then
    cp support-files/external-db-values-mysql.yml manifest/ -f
    sed -i "s/<external_db_ip>/$external_db_ip/" manifest/external-db-values-mysql.yml
    sed -i "s/<external_db_port>/$external_db_port/" manifest/external-db-values-mysql.yml
    sed -i "s/<external_cc_db_id>/$external_cc_db_id/" manifest/external-db-values-mysql.yml
    sed -i "s/<external_cc_db_password>/$external_cc_db_password/" manifest/external-db-values-mysql.yml
    sed -i "s/<external_cc_db_name>/$external_cc_db_name/" manifest/external-db-values-mysql.yml
    sed -i "s/<external_uaa_db_id>/$external_uaa_db_id/" manifest/external-db-values-mysql.yml
    sed -i "s/<external_uaa_db_password>/$external_uaa_db_password/" manifest/external-db-values-mysql.yml
    sed -i "s/<external_uaa_db_name>/$external_uaa_db_name/" manifest/external-db-values-mysql.yml
    render_ca=$(awk '{print "      "$0}' $external_db_cert_path)
    awk -v ca_values="$render_ca" 'NR==24{print ca_values}1' manifest/external-db-values-mysql.yml > manifest/external-db-values-mysql-temp.yml
    awk -v ca_values="$render_ca" 'NR==13{print ca_values}1' manifest/external-db-values-mysql-temp.yml > manifest/external-db-values-mysql.yml
    rm manifest/external-db-values-mysql-temp.yml
  else
    echo "plz check variables.yml : external_db_kind"
    return
  fi
fi
