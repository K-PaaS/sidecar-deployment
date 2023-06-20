#!/bin/bash

source portal-app-variable.yml

PORTALAPPNAME=portal-app-1.2.13
PORTALAPPDOWNLOADLINK=https://nextcloud.paas-ta.org/index.php/s/6aanBz8osifGnQZ/download

#########################################
# Portal Component Folder Name
PORTAL_API=portal-api-2.4.3
PORTAL_COMMON_API=portal-common-api-2.2.6
PORTAL_GATEWAY=portal-gateway-2.1.2
PORTAL_LOG_API=portal-log-api-2.3.2
PORTAL_REGISTRATION=portal-registration-2.1.0
PORTAL_STORAGE_API=portal-storage-api-2.2.1
PORTAL_WEB_ADMIN=portal-web-admin-2.3.5
PORTAL_WEB_USER=portal-web-user-2.4.9
PORTAL_SSH=portal-ssh-1.0.0

#########################################
# language list check
PORTAL_WEB_USER_INPUT_LANG='ko,en'
PORTAL_WEB_ADMIN_INPUT_LANG='ko,en'
IFS=',' read -r -a PORTAL_WEB_USER_LANG <<< "$PORTAL_WEB_USER_INPUT_LANG"
IFS=',' read -r -a PORTAL_WEB_ADMIN_LANG <<< "$PORTAL_WEB_ADMIN_INPUT_LANG"
#WEB_USER_LANG_COMP=($(printf "%s\n" "${PORTAL_WEB_USER_LANG[@]}" /| sort -u))
#WEB_ADMIN_LANG_COMP=($(printf "%s\n" "${PORTAL_WEB_ADMIN_LANG[@]}" | sort -u))
PORTAL_WEB_USER_LANGUAGE=()
for lang in ${PORTAL_WEB_USER_LANG[@]}
do
        if [[ ! ${PORTAL_WEB_USER_LANGUAGE[@]} =~ ${lang} ]]; then
                PORTAL_WEB_USER_LANGUAGE+=($lang)
        fi
done
PORTAL_WEB_ADMIN_LANGUAGE=()
for lang in ${PORTAL_WEB_ADMIN_LANG[@]}
do
        if [[ ! ${PORTAL_WEB_ADMIN_LANGUAGE[@]} =~ ${lang} ]]; then
                PORTAL_WEB_ADMIN_LANGUAGE+=($lang)
        fi
done


PORTAL_WEB_USER_STR_CHECK=$(grep -r "portal_web_user_language" $COMMON_VARS_PATH | cut -d ':' -f 2 | cut -d '#' -f 1 | cut -f 1)
PORTAL_WEB_ADMIN_STR_CHECK=$(grep -r "portal_web_admin_language" $COMMON_VARS_PATH | cut -d ':' -f 2 | cut -d '#' -f 1 | cut -f 1)

if [[ ${#PORTAL_WEB_USER_LANGUAGE[@]} -eq 0 ]]; then
        if [[ "${PORTAL_WEB_USER_STR_CHECK}" != *"["* ]] && [[ "${PORTAL_WEB_USER_STR_CHECK}" != *"]"* ]]; then
                PORTAL_WEB_USER_LANGUAGE=("ko" "en")
        else
                echo "Language list dose not exist -> portal_web_user_language check plz"
                return
        fi
fi

if [[ ${#PORTAL_WEB_ADMIN_LANGUAGE[@]} -eq 0 ]]; then
        if [[ "${PORTAL_WEB_ADMIN_STR_CHECK}" != *"["* ]] && [[ "${PORTAL_WEB_ADMIN_STR_CHECK}" != *"]"* ]]; then
                PORTAL_WEB_ADMIN_LANGUAGE=("ko" "en")
        else
                echo "Language list dose not exist -> portal_web_admin_language check plz"
                return
        fi
fi

PORTAL_WEB_USER_USE_LANG=$(echo "${PORTAL_WEB_USER_LANGUAGE[*]}" | sed 's/ /,/g')
PORTAL_WEB_ADMIN_USE_LANG=$(echo "${PORTAL_WEB_ADMIN_LANGUAGE[*]}" | sed 's/ /,/g')


# VARIABLE SETTING
SYSTEM_DOMAIN=$(yq e '.system_domain' $SIDECAR_VALUES_PATH)
APP_DOMAIN=$(yq e '.app_domains[0]' $SIDECAR_VALUES_PATH)
CF_USER_ADMIN_USERNAME="admin"
CF_USER_ADMIN_PASSWORD=$(yq e '.cf_admin_password' $SIDECAR_VALUES_PATH)
UAA_ADMIN_CLIENT_SECRET=$(yq e '.uaa.admin_client_secret' $SIDECAR_VALUES_PATH)

## PORTAL DB
if [[ ${IS_PORTAL_EXTERNAL_DB} = "false" ]]; then
        # Portal - Internal DB use
        PORTAL_DB_IP="mariadb."$NAMESPACE_NAME".svc.cluster.local"
        PORTAL_DB_PORT=$MARIADB_SERVICE_PORT
        PORTAL_DB_USER_PASSWORD=$MARIADB_PASSWORD

elif [[ ${IS_PORTAL_EXTERNAL_DB} = "true" ]]; then
        # Portal - External DB use
        PORTAL_DB_IP=$PORTAL_EXTERNAL_DB_IP
        PORTAL_DB_PORT=$PORTAL_EXTERNAL_DB_PORT
        PORTAL_DB_USER_PASSWORD=$PORTAL_EXTERNAL_DB_PASSWORD

else
        # unknown IS_PORTAL_EXTERNAL_DB value
        echo "plz check IS_PORTAL_EXTERNAL_DB"
        return
fi

## SKIP
MONITORING_API_URL=SKIP

## AP DB
if [[ ${IS_PAAS_TA_EXTERNAL_DB} = "false" ]]; then
        # AP - Internal DB use
        # if [[  ]]; then
                PAASTA_DB_DRIVER=org.postgresql.Driver
                PAASTA_DATABASE=postgresql
        # elif [[  ]]; then
                #PAASTA_DB_DRIVER=com.mysql.jdbc.Driver
                #PAASTA_DATABASE=mysql
        # fi
        PAASTA_DB_IP="cf-db-postgresql.cf-db.svc.cluster.local"
        PAASTA_DB_PORT="5432"

elif [[ ${IS_PAAS_TA_EXTERNAL_DB} = "true" ]]; then
        # AP - External DB use
        if [[ $PAAS_TA_EXTERNAL_DB_KIND = "postgres" ]]; then
                PAASTA_DB_DRIVER=org.postgresql.Driver
                PAASTA_DATABASE=postgresql
        elif [[ $PAAS_TA_EXTERNAL_DB_KIND = "mysql" ]]; then
                PAASTA_DB_DRIVER=com.mysql.jdbc.Driver
                PAASTA_DATABASE=mysql
        else
                echo "plz check IS_PAAS_TA_EXTERNAL_DB & PAAS_TA_EXTERNAL_DB_KIND"
                return
        fi

        PAASTA_DB_IP=$PAAS_TA_EXTERNAL_DB_IP
        PAASTA_DB_PORT=$PAAS_TA_EXTERNAL_DB_PORT
else
        # unknown IS_PAAS_TA_EXTERNAL_DB value
        echo "plz check IS_PAAS_TA_EXTERNAL_DB"
        return
fi

CC_DB_USER_PASSWORD=$(yq e '.capi.database.password' $SIDECAR_VALUES_PATH)
UAA_DB_USER_PASSWORD=$(yq e '.uaa.database.password' $SIDECAR_VALUES_PATH)
MAIL_SMTP_PROPERTIES_AUTHURL=portal-web-user.$APP_DOMAIN


## OBJECT STORAGE
if [[ ${IS_PORTAL_EXTERNAL_STORAGE} = "false" ]]; then
        # Portal - Internal Storage use
        OBJECTSTORAGE_TENANTNAME=$PORTAL_OBJECTSTORAGE_TENANTNAME
        OBJECTSTORAGE_USERNAME=$PORTAL_OBJECTSTORAGE_USERNAME
        OBJECTSTORAGE_PASSWORD=$PORTAL_OBJECTSTORAGE_PASSWORD
        OBJECTSTORAGE_IP="openstack-swift-keystone-docker."$NAMESPACE_NAME".svc.cluster.local"
        OBJECTSTORAGE_PORT=$KEYSTONE_SERVICE_PORT

elif [[ ${IS_PORTAL_EXTERNAL_STORAGE} = "true" ]]; then
        # Portal - External Storage use
        OBJECTSTORAGE_TENANTNAME=$PORTAL_EXTERNAL_STORAGE_TENANTNAME
        OBJECTSTORAGE_USERNAME=$PORTAL_EXTERNAL_STORAGE_USERNAME
        OBJECTSTORAGE_PASSWORD=$PORTAL_EXTERNAL_STORAGE_PASSWORD
        OBJECTSTORAGE_IP=$PORTAL_EXTERNAL_STORAGE_IP
        OBJECTSTORAGE_PORT=$PORTAL_EXTERNAL_STORAGE_PORT
else
        # unknown IS_PORTAL_EXTERNAL_STORAGE value
        echo "plz check IS_PORTAL_EXTERNAL_STORAGE"
        return
fi

UAAC_PORTAL_CLIENT_SECRET=$(yq e '.uaa_client_portal_secret' $SIDECAR_VALUES_PATH)
API_TYPE=$API_TYPE
SSH_ENABLE=$SSH_ENABLE
TAIL_LOG_INTERVAL=$TAIL_LOG_INTERVAL

#########################################

CURRENTDIRCTORY=$(pwd)

mkdir $PORTAL_APP_WORKING_DIRECTORY -p
if [ -d $PORTAL_APP_WORKING_DIRECTORY ]; then
        cd $PORTAL_APP_WORKING_DIRECTORY
else
        echo "plz check PORTAL_APP_WORKING_DIRECTORY"
        cd $CURRENTDIRCTORY
        return
fi

# portal-app download
## portal-app zip downloaded check
if [ -e $PORTALAPPNAME.zip ]; then
        echo "portal-app zip file exists - download skip"
else
        echo "portal-app zip file not exists - download zip file "
        ## portal-app wget download
        wget --content-disposition $PORTALAPPDOWNLOADLINK
fi

if [ ! -e $PORTALAPPNAME.zip ]; then
        echo "plz check portal app download link : "$PORTALAPPDOWNLOADLINK
        cd $CURRENTDIRCTORY
        return
fi


#########################################
# portal-app unzip
if [ -d $PORTALAPPNAME ]; then
        echo "portal-app folder exists - delete folder"
        ## existing portal-app directory delete
        ll | grep "$PORTALAPPNAME" | grep ^d | awk '{print $NF}' | xargs rm -rf
fi
## portal-app unzip
echo "portal-app unzip"
unzip -q $PORTALAPPNAME.zip

#########################################
#config change

if [ -d $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME ]; then
        cd $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME
else
        echo "plz check directory : " $PORTAL_APP_WORKING_DIRECTORY"/"$PORTALAPPNAME
        cd $CURRENTDIRCTORY
        return
fi




## COMMON VARIABLE
# SYSTEM_DOMAIN
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_API/manifest.yml -type f | xargs sed -i -e 's/<SYSTEM_DOMAIN>/'${SYSTEM_DOMAIN}'/g'
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_LOG_API/manifest.yml -type f | xargs sed -i -e 's/<SYSTEM_DOMAIN>/'${SYSTEM_DOMAIN}'/g'

# APP_DOMAIN
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_API/manifest.yml -type f | xargs sed -i -e 's/<APP_DOMAIN>/'${APP_DOMAIN}'/g'
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -type f | xargs sed -i -e 's/<APP_DOMAIN>/'${APP_DOMAIN}'/g'
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_GATEWAY/manifest.yml -type f | xargs sed -i -e 's/<APP_DOMAIN>/'${APP_DOMAIN}'/g'
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_LOG_API/manifest.yml -type f | xargs sed -i -e 's/<APP_DOMAIN>/'${APP_DOMAIN}'/g'
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_REGISTRATION/manifest.yml -type f | xargs sed -i -e 's/<APP_DOMAIN>/'${APP_DOMAIN}'/g'
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_STORAGE_API/manifest.yml -type f | xargs sed -i -e 's/<APP_DOMAIN>/'${APP_DOMAIN}'/g'
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_ADMIN/manifest.yml -type f | xargs sed -i -e 's/<APP_DOMAIN>/'${APP_DOMAIN}'/g'
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_USER/config -type f | xargs sed -i -e 's/<APP_DOMAIN>/'${APP_DOMAIN}'/g'


# CF_USER_ADMIN_USERNAME
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_API/manifest.yml -type f | xargs sed -i -e 's/<CF_USER_ADMIN_USERNAME>/'${CF_USER_ADMIN_USERNAME}'/g'
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_LOG_API/manifest.yml -type f | xargs sed -i -e 's/<CF_USER_ADMIN_USERNAME>/'${CF_USER_ADMIN_USERNAME}'/g'

# CF_USER_ADMIN_PASSWORD
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_API/manifest.yml -type f | xargs sed -i -e 's/<CF_USER_ADMIN_PASSWORD>/'${CF_USER_ADMIN_PASSWORD}'/g'
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_LOG_API/manifest.yml -type f | xargs sed -i -e 's/<CF_USER_ADMIN_PASSWORD>/'${CF_USER_ADMIN_PASSWORD}'/g'


# UAA_CLIENT_ID
## UAA_ADMIN_CLIENT_ID == UAA_CLIENT_ID
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_API/manifest.yml -type f | xargs sed -i -e 's/<UAA_CLIENT_ID>/'${UAA_ADMIN_CLIENT_ID}'/g'
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_LOG_API/manifest.yml -type f | xargs sed -i -e 's/<UAA_CLIENT_ID>/'${UAA_ADMIN_CLIENT_ID}'/g'

# UAA_CLIENT_SECRET
## UAA_ADMIN_CLIENT_SECRET == UAA_CLIENT_SECRET
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_API/manifest.yml -type f | xargs sed -i -e 's/<UAA_CLIENT_SECRET>/'${UAA_ADMIN_CLIENT_SECRET}'/g'
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_LOG_API/manifest.yml -type f | xargs sed -i -e 's/<UAA_CLIENT_SECRET>/'${UAA_ADMIN_CLIENT_SECRET}'/g'

# UAA_ADMIN_CLIENT_ID
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_API/manifest.yml -type f | xargs sed -i -e 's/<UAA_ADMIN_CLIENT_ID>/'${UAA_ADMIN_CLIENT_ID}'/g'
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_LOG_API/manifest.yml -type f | xargs sed -i -e 's/<UAA_ADMIN_CLIENT_ID>/'${UAA_ADMIN_CLIENT_ID}'/g'

# UAA_ADMIN_CLIENT_SECRET
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_API/manifest.yml -type f | xargs sed -i -e 's/<UAA_ADMIN_CLIENT_SECRET>/'${UAA_ADMIN_CLIENT_SECRET}'/g'
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_LOG_API/manifest.yml -type f | xargs sed -i -e 's/<UAA_ADMIN_CLIENT_SECRET>/'${UAA_ADMIN_CLIENT_SECRET}'/g'

# UAA_LOGIN_CLIENT_ID
## UAA_ADMIN_CLIENT_ID == UAA_LOGIN_CLIENT_ID
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_API/manifest.yml -type f | xargs sed -i -e 's/<UAA_LOGIN_CLIENT_ID>/'${UAA_ADMIN_CLIENT_ID}'/g'
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_LOG_API/manifest.yml -type f | xargs sed -i -e 's/<UAA_LOGIN_CLIENT_ID>/'${UAA_ADMIN_CLIENT_ID}'/g'

# UAA_LOGIN_CLIENT_SECRET
## UAA_ADMIN_CLIENT_SECRET == UAA_LOGIN_CLIENT_SECRET
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_API/manifest.yml -type f | xargs sed -i -e 's/<UAA_LOGIN_CLIENT_SECRET>/'${UAA_ADMIN_CLIENT_SECRET}'/g'
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_LOG_API/manifest.yml -type f | xargs sed -i -e 's/<UAA_LOGIN_CLIENT_SECRET>/'${UAA_ADMIN_CLIENT_SECRET}'/g'

# PORTAL_DB_IP
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -type f | xargs sed -i -e 's/<PORTAL_DB_IP>/'${PORTAL_DB_IP}'/g'
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_ADMIN/manifest.yml -type f | xargs sed -i -e 's/<PORTAL_DB_IP>/'${PORTAL_DB_IP}'/g'

# PORTAL_DB_PORT
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -type f | xargs sed -i -e 's/<PORTAL_DB_PORT>/'${PORTAL_DB_PORT}'/g'
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_ADMIN/manifest.yml -type f | xargs sed -i -e 's/<PORTAL_DB_PORT>/'${PORTAL_DB_PORT}'/g'

# PORTAL_DB_USER_PASSWORD
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -type f | xargs sed -i -e 's/<PORTAL_DB_USER_PASSWORD>/'${PORTAL_DB_USER_PASSWORD}'/g'
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_ADMIN/manifest.yml -type f | xargs sed -i -e 's/<PORTAL_DB_USER_PASSWORD>/'${PORTAL_DB_USER_PASSWORD}'/g'


## PORTAL-API
# ABACUS_URL(Deprecated)
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_API/manifest.yml -type f | xargs sed -i -e 's/<ABACUS_URL>/'${ABACUS_URL}'/g'
# MONITORING_API_URL
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_API/manifest.yml -type f | xargs sed -i -e 's/<MONITORING_API_URL>/'${MONITORING_API_URL}'/g'
# API_TYPE
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_API/manifest.yml -type f | xargs sed -i -e 's/<API_TYPE>/'${API_TYPE}'/g'


## PORTAL-COMMON-API
# PAASTA_DB_DRIVER
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -type f | xargs sed -i -e 's/<PAAS-TA_DB_DRIVER>/'${PAASTA_DB_DRIVER}'/g'

# PAASTA_DATABASE
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -type f | xargs sed -i -e 's/<PAAS-TA_DATABASE>/'${PAASTA_DATABASE}'/g'

# PAASTA_DB_IP
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -type f | xargs sed -i -e "s/<PAAS-TA_DB_IP>/$PAASTA_DB_IP/g"

# PAASTA_DB_PORT
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -type f | xargs sed -i -e 's/<PAAS-TA_DB_PORT>/'${PAASTA_DB_PORT}'/g'

# CC_DB_NAME
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -type f | xargs sed -i -e 's/<CC_DB_NAME>/'${CC_DB_NAME}'/g'

# CC_DB_USER_NAME
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -type f | xargs sed -i -e 's/<CC_DB_USER_NAME>/'${CC_DB_USER_NAME}'/g'

# CC_DB_USER_PASSWORD
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -type f | xargs sed -i -e 's/<CC_DB_USER_PASSWORD>/'${CC_DB_USER_PASSWORD}'/g'

# UAA_DB_NAME
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -type f | xargs sed -i -e 's/<UAA_DB_NAME>/'${UAA_DB_NAME}'/g'

# UAA_DB_USER_NAME
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -type f | xargs sed -i -e 's/<UAA_DB_USER_NAME>/'${UAA_DB_USER_NAME}'/g'

# UAA_DB_USER_PASSWORD
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -type f | xargs sed -i -e 's/<UAA_DB_USER_PASSWORD>/'${UAA_DB_USER_PASSWORD}'/g'

# MAIL_SMTP_HOST
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -type f | xargs sed -i -e 's/<MAIL_SMTP_HOST>/'${MAIL_SMTP_HOST}'/g'

# MAIL_SMTP_PORT
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -type f | xargs sed -i -e 's/<MAIL_SMTP_PORT>/'${MAIL_SMTP_PORT}'/g'

# MAIL_SMTP_USERNAME
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -type f | xargs sed -i -e 's/<MAIL_SMTP_USERNAME>/'${MAIL_SMTP_USERNAME}'/g'

# MAIL_SMTP_PASSWORD
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -type f | xargs sed -i -e 's/<MAIL_SMTP_PASSWORD>/'${MAIL_SMTP_PASSWORD}'/g'

# MAIL_SMTP_USEREMAIL
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -type f | xargs sed -i -e 's/<MAIL_SMTP_USEREMAIL>/'${MAIL_SMTP_USEREMAIL}'/g'

# MAIL_SMTP_PROPERTIES_AUTHURL
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -type f | xargs sed -i -e 's/<MAIL_SMTP_PROPERTIES_AUTHURL>/'${MAIL_SMTP_PROPERTIES_AUTHURL}'/g'

# PORTAL_USE_LANGUAGE
COMMON_API_DIRECTORY=$PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API
APP_CONFIG=$COMMON_API_DIRECTORY/manifest.yml
SEARCH_FILTER=$(unzip -q -l ${COMMON_API_DIRECTORY}/paas-ta-portal-common-api.jar | grep "template/" | cut -d "/" -f4 | uniq)

ORIGIN_LANG=()
for lang in $SEARCH_FILTER
do
        if [[ -n "${lang}" ]]; then
                ORIGIN_LANG+=(${lang})
        fi
done

for element in ${PORTAL_WEB_USER_LANGUAGE[@]}
do
        if [[ ! ${ORIGIN_LANG[@]} =~ ${element} ]]; then
                echo "\"${element}\" is unsupported language -> portal_web_user_language check plz"
                return
        fi
done

find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -type f | xargs sed -i -e 's/<PORTAL_USE_LANGUAGE>/'${PORTAL_WEB_USER_USE_LANG}'/g'



## PORTAL-LOG-API
# LOGGING_INFLUXDB_IP
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_LOG_API/manifest.yml -type f | xargs sed -i -e 's/<LOGGING_INFLUXDB_IP>/'${influxdb_ip}'/g'

# LOGGING_INFLUXDB_PORT
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_LOG_API/manifest.yml -type f | xargs sed -i -e 's/<LOGGING_INFLUXDB_PORT>/'${influxdb_http_port}'/g'

# LOGGING_INFLUXDB_USERNAME
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_LOG_API/manifest.yml -type f | xargs sed -i -e 's/<LOGGING_INFLUXDB_USERNAME>/'${influxdb_username}'/g'

# LOGGING_INFLUXDB_PASSWORD
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_LOG_API/manifest.yml -type f | xargs sed -i -e 's/<LOGGING_INFLUXDB_PASSWORD>/'${influxdb_password}'/g'

# LOGGING_INFLUXDB_DATABASE
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_LOG_API/manifest.yml -type f | xargs sed -i -e 's/<LOGGING_INFLUXDB_DATABASE>/'${influxdb_database}'/g'

# LOGGING_INFLUXDB_MEASUREMENT
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_LOG_API/manifest.yml -type f | xargs sed -i -e 's/<LOGGING_INFLUXDB_MEASUREMENT>/'${influxdb_measurement}'/g'

# LOGGING_INFLUXDB_LIMIT
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_LOG_API/manifest.yml -type f | xargs sed -i -e 's/<LOGGING_INFLUXDB_LIMIT>/'${influxdb_limit}'/g'

# LOGGING_INFLUXDB_HTTPS_ENABLED
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_LOG_API/manifest.yml -type f | xargs sed -i -e 's/<LOGGING_INFLUXDB_HTTPS_ENABLED>/'${influxdb_https_enabled}'/g'

# LOGGING_INFLUXDB_URL
if [[ ${influxdb_https_enabled} = false ]]; then
        find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_LOG_API/manifest.yml -type f | xargs sed -i -e 's/influxdb_url: https/influxdb_url: http/g'
fi



## PORTAL-STORAGE-API
# OBJECTSTORAGE_TENANTNAME
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_STORAGE_API/manifest.yml -type f | xargs sed -i -e 's/<OBJECTSTORAGE_TENANTNAME>/'${OBJECTSTORAGE_TENANTNAME}'/g'

# OBJECTSTORAGE_USERNAME
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_STORAGE_API/manifest.yml -type f | xargs sed -i -e 's/<OBJECTSTORAGE_USERNAME>/'${OBJECTSTORAGE_USERNAME}'/g'

# OBJECTSTORAGE_PASSWORD
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_STORAGE_API/manifest.yml -type f | xargs sed -i -e 's/<OBJECTSTORAGE_PASSWORD>/'${OBJECTSTORAGE_PASSWORD}'/g'

# OBJECTSTORAGE_IP
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_STORAGE_API/manifest.yml -type f | xargs sed -i -e 's/<OBJECTSTORAGE_IP>/'${OBJECTSTORAGE_IP}'/g'

# OBJECTSTORAGE_PORT
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_STORAGE_API/manifest.yml -type f | xargs sed -i -e 's/<OBJECTSTORAGE_PORT>/'${OBJECTSTORAGE_PORT}'/g'



## PORTAL-WEBADMIN
# PORTAL_USE_LANGUAGE
WEB_ADMIN_DIRECTORY=$PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_ADMIN
APP_CONFIG=$WEB_ADMIN_DIRECTORY/manifest.yml
SEARCH_FILTER=$(unzip -q -l $WEB_ADMIN_DIRECTORY/paas-ta-portal-webadmin.war | grep "message_" | cut -d "_" -f2)
BEFORE_LANG_LIST=$(grep "languageList" $APP_CONFIG | sed -e 's/^ *//g')

ORIGIN_LANG=()
for lang in $SEARCH_FILTER
do
        ORIGIN_LANG+=(`basename -s ".properties" "${lang}"`)
done

for element in ${PORTAL_WEB_ADMIN_LANGUAGE[@]}
do
        if [[ ! ${ORIGIN_LANG[@]} =~ ${element} ]]; then
                echo "\"${element}\" is unsupported language -> portal_web_admin_language check plz"
                return
        fi
done

find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_ADMIN/manifest.yml -type f | xargs sed -i -e 's/<PORTAL_USE_LANGUAGE>/'${PORTAL_WEB_ADMIN_USE_LANG}'/g'

AFTER_LANG_LIST=$(grep "languageList" $APP_CONFIG | sed -e 's/^ *//g')

echo "====================================================="
echo "BEFORE :: $BEFORE_LANG_LIST"
echo "====================================================="
echo "AFTER  :: $AFTER_LANG_LIST"
echo "====================================================="



## PORTAL-WEBUSER
# UAAC_PORTAL_CLIENT_ID
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_USER/config -type f | xargs sed -i -e 's/<UAAC_PORTAL_CLIENT_ID>/'${UAAC_PORTAL_CLIENT_ID}'/g'

# UAAC_PORTAL_CLIENT_SECRET
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_USER/config -type f | xargs sed -i -e 's/<UAAC_PORTAL_CLIENT_SECRET>/'${UAAC_PORTAL_CLIENT_SECRET}'/g'

# USER_APP_SIZE_MB
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_USER/config -type f | xargs sed -i -e 's/<USER_APP_SIZE_MB>/'${USER_APP_SIZE_MB}'/g'

# MONITORING_ENABLE
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_USER/config -type f | xargs sed -i -e 's/<MONITORING_ENABLE>/'${MONITORING_ENABLE}'/g'

# SSH_ENABLE
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_USER/config -type f | xargs sed -i -e 's/<SSH_ENABLE>/'${SSH_ENABLE}'/g'

# TAIL_LOG_INTERVAL
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_USER/config -type f | xargs sed -i -e 's/<TAIL_LOG_INTERVAL>/'${TAIL_LOG_INTERVAL}'/g'

# PAASTA_DEPLOY_TYPE
find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_USER/config -type f | xargs sed -i -e 's/<PAASTA_DEPLOY_TYPE>/'${API_TYPE}'/g'

# PORTAL_USE_LANGUAGE
PORTAL_WEB_USER_USE_LANG_LIST=$(echo "[\"${PORTAL_WEB_USER_LANGUAGE[*]}\"]" | sed 's/ /\",\"/g')

find $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_USER/config -type f | xargs sed -i -e 's/<PORTAL_USE_LANGUAGE>/'${PORTAL_WEB_USER_USE_LANG_LIST}'/g'

# PORTAL WEBUSER MAIN
BEFORE_CONFIG=$PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_USER/paas-ta-portal-webuser/assets/resources/env/config.json
AFTER_CONFIG=$PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_USER/config/config.json
MAIN_JS=$PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_USER/paas-ta-portal-webuser/main.*.js
LANG_DIR_PATH=$PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_USER/paas-ta-portal-webuser/assets/i18n

BEFORE_LANG=()
for file in `ls $LANG_DIR_PATH/*`
do
        BEFORE_LANG+=(`basename -s ".json" "${file}"`)
done

AFTER_LANG=$(grep "languageList" ${AFTER_CONFIG} | tr -d '\[' | tr -d '\]' | tr -d '"' | tr -d ' ' | tr -d '\r' | cut -d ":" -f2)
IFS=',' read -r -a AFTER_LANG_LIST <<< "$AFTER_LANG"

for element in ${AFTER_LANG_LIST[@]}
do
        if [[ ! ${BEFORE_LANG[@]} =~ ${element} ]]; then
                echo "\"${element}\" is unsupported language -> plz check portal_web_user_language"
                return
        fi
done

BEFORE_FILTER=$(cat ${BEFORE_CONFIG} | tr -d '{'  |tr -d '\r\n' | tr -d '"' | sed -e 's/: /:\"/g' | sed -e 's/,  /\",/g' | sed -e 's/^ *//g' -e 's/ *$//g' | sed -e 's/}/"/g' | sed -e 's/"false"/!1/g' | sed -e 's/"true"/!0/g'| sed -e 's/\//\\\//g' | sed -e 's/, /\",\"/g' | sed -e 's/\"\[/\\[\"/g' | sed -e 's/\]\"/\"\\]/g')
AFTER_FILTER=$(cat ${AFTER_CONFIG} | tr -d '{'  |tr -d '\r\n' | tr -d '"' | sed -e 's/: /:\"/g' | sed -e 's/,/\",\"/g' | sed -e 's/\"  //g' | sed -e 's/^ *//g' -e 's/ *$//g' | sed -e 's/}/"/g' | sed -e 's/"false"/!1/g' | sed -e 's/"true"/!0/g'| sed -e 's/\//\\\//g' | sed -e 's/\"\[/\\[\"/g' | sed -e 's/\]\"/\"\\]/g' | sed -e 's/\" /\"/g')

echo "====================================================="
echo "BEFORE :: $BEFORE_FILTER"
echo "====================================================="
echo "AFTER  :: $AFTER_FILTER"
echo "====================================================="

CHANGE_CONFIG="'s/${BEFORE_FILTER}/${AFTER_FILTER}/g' ${MAIN_JS}"

echo $CHANGE_CONFIG | xargs sed -i

cp $AFTER_CONFIG $BEFORE_CONFIG

## nginx.conf (nginx.conf change relative path portal-webuser 2.4.8)
#sed -i '/root/ c\    root /workspace;' $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_USER/paas-ta-portal-webuser/nginx.conf


#########################################
# Portal App push


cf login -a https://api.${SYSTEM_DOMAIN} --skip-ssl-validation -u ${CF_USER_ADMIN_USERNAME} -p ${CF_USER_ADMIN_PASSWORD} << EOF
EOF

# Create Portal Org, Space
cf create-quota ${PORTAL_QUOTA_NAME} -m 20G -i -1 -s -1 -r -1 --reserved-route-ports -1 --allow-paid-service-plans
cf create-org ${PORTAL_ORG_NAME} -q ${PORTAL_QUOTA_NAME}
cf create-space ${PORTAL_SPACE_NAME} -o ${PORTAL_ORG_NAME}

cf target -o ${PORTAL_ORG_NAME} -s ${PORTAL_SPACE_NAME}


# Portal APP push
cf push -i $PORTAL_API_INSTANCE -f $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_API/manifest.yml -b paketo-buildpacks/java
cf push -i $PORTAL_COMMON_API_INSTANCE -f $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_COMMON_API/manifest.yml -b paketo-buildpacks/java
cf push -i $PORTAL_GATEWAY_INSTANCE -f $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_GATEWAY/manifest.yml -b paketo-buildpacks/java
cf push -i $PORTAL_REGISTRATION_INSTANCE -f $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_REGISTRATION/manifest.yml -b paketo-buildpacks/java
cf push -i $PORTAL_STORAGE_API_INSTANCE -f $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_STORAGE_API/manifest.yml -b paketo-buildpacks/java
cf push -i $PORTAL_WEB_ADMIN_INSTANCE -f $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_ADMIN/manifest.yml -b paketo-buildpacks/java
cf push -i $PORTAL_WEB_USER_INSTANCE -f $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_WEB_USER/manifest.yml -b paketo-buildpacks/nginx

if [[ ${use_logging_service} = true ]]; then
        cf push -f $PORTAL_APP_WORKING_DIRECTORY/$PORTALAPPNAME/$PORTAL_LOG_API/manifest.yml -b paketo-buildpacks/java
fi

cf apps

cd $CURRENTDIRCTORY
