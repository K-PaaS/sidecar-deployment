#!/bin/bash

#VARIABLES
INGRESS_HOST_DOMAIN="{host domain}"		# Host Domain (e.g. xx.xxx.xxx.xx.nip.io)

CP_PORTAL_VERSION=v1.6.0

# SCRIPT START

# 1. CP PORTAL DOWNLOAD
if [ -e cp-portal-deployment-$CP_PORTAL_VERSION.tar.gz ]; then
        echo "cp-portal-deployment file exists - download skip"
else
        echo "cp-portal-deployment file not exists - download zip file "
        ## cp-portal-deployment wget download
        wget --content-disposition https://nextcloud.k-paas.org/index.php/s/ZcFt4cpeXj8d4o4/download
fi
tar -xvf cp-portal-deployment-$CP_PORTAL_VERSION.tar.gz

# 2. CP PORTAL TEMPLATE COPY
cp cp-portal-template-$CP_PORTAL_VERSION/script/cp-portal-vars.sh cp-portal-deployment/script/cp-portal-vars.sh 
cp cp-portal-template-$CP_PORTAL_VERSION/script/deploy-cp-portal.sh cp-portal-deployment/script/deploy-cp-portal.sh
cp cp-portal-template-$CP_PORTAL_VERSION/values_orig/cp-portal.yaml cp-portal-deployment/values_orig/cp-portal.yaml

# 3. HOST_DOMAIN INPUT

sed -i 's/INGRESS_HOST_DOMAIN/'${INGRESS_HOST_DOMAIN}'/g' cp-portal-deployment/script/cp-portal-vars.sh


# 4. deploy
cd cp-portal-deployment/script
chmod +x deploy-cp-portal.sh
./deploy-cp-portal.sh
cd ../..


