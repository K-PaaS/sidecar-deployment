#!/bin/bash

#VARIABLES
INGRESS_HOST_DOMAIN="{host domain}"  		# Host Domain (e.g. xx.xxx.xxx.xx.nip.io)

CP_PORTAL_VERSION=v1.5.2

# SCRIPT START

# 1. CP PORTAL DOWNLOAD
wget --content-disposition https://nextcloud.k-paas.org/index.php/s/2LeyyQTaCySmKzH/download
tar -xvf cp-portal-deployment-$CP_PORTAL_VERSION.tar.gz

# 2. CP PORTAL TEMPLATE COPY
cp cp-portal-template-$CP_PORTAL_VERSION/script/cp-portal-vars.sh cp-portal-deployment/script/cp-portal-vars.sh 
cp cp-portal-template-$CP_PORTAL_VERSION/script/deploy-cp-portal.sh cp-portal-deployment/script/deploy-cp-portal.sh
cp cp-portal-template-$CP_PORTAL_VERSION/values_orig/cp-portal.yaml cp-portal-deployment/values_orig/cp-portal.yaml

# 3. HOST_DOMAIN INPUT

sed -i -e 's/INGRESS_HOST_DOMAIN/'${INGRESS_HOST_DOMAIN}'/g' cp-portal-deployment/script/cp-portal-vars.sh

# 4. deploy
cd cp-portal-deployment/script
chmod +x deploy-cp-portal.sh
./deploy-cp-portal.sh
cd ../..


