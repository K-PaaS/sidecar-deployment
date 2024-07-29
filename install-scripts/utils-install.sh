#!/bin/bash

architecture=""
case $(uname -m) in
    i386)   architecture="386" ;;
    i686)   architecture="386" ;;
    x86_64) architecture="amd64" ;;
    arm)    dpkg --print-architecture | grep -q "arm64" && architecture="arm64" || architecture="arm" ;;
esac
OS=$(uname)

echo "------------------"
echo "ytt & kapp install"
echo "------------------"
wget -O- https://carvel.dev/install.sh > install.sh
sudo bash install.sh
rm install.sh

echo "------------------"
echo "cf cli install"
echo "------------------"
curl -L "https://packages.cloudfoundry.org/stable?release=linux64-binary&version=v8&source=github" | tar -zx
sudo mv cf8 /usr/local/bin
sudo mv cf /usr/local/bin
sudo curl -o /usr/share/bash-completion/completions/cf8 https://raw.githubusercontent.com/cloudfoundry/cli-ci/master/ci/installers/completion/cf8
cf version
rm LICENSE NOTICE

echo "------------------"
echo "yq cli install"
echo "------------------"
curl -L https://github.com/mikefarah/yq/releases/download/v4.30.4/yq_$OS\_$architecture -o yq
chmod +x yq
sudo mv yq /usr/local/bin/
yq --version


echo "------------------"
echo "cert-manager cli install"
echo "------------------"
curl -L -o cmctl.tar.gz https://github.com/cert-manager/cert-manager/releases/download/v1.14.7/cmctl-$OS-$architecture.tar.gz
tar xzf cmctl.tar.gz
rm cmctl.tar.gz LICENSE
sudo mv cmctl /usr/local/bin
