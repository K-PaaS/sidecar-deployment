#!/bin/bash

echo "------------------"
echo "ytt & kapp install"
echo "------------------"
sudo su -c "wget -O- https://k14s.io/install.sh | bash"
sudo chown $(id -u):$(id -g) /usr/local/bin/ytt
sudo chown $(id -u):$(id -g) /usr/local/bin/kapp

echo "------------------"
echo "bosh cli install"
echo "------------------"
sudo apt update
curl -Lo ./bosh https://github.com/cloudfoundry/bosh-cli/releases/download/v6.1.0/bosh-cli-6.1.0-linux-amd64
chmod +x ./bosh
sudo mv ./bosh /usr/local/bin/bosh
bosh -v

echo "------------------"
echo "cf cli install"
echo "------------------"
wget -q -O - https://packages.cloudfoundry.org/debian/cli.cloudfoundry.org.key | sudo apt-key add -
echo "deb https://packages.cloudfoundry.org/debian stable main" | sudo tee /etc/apt/sources.list.d/cloudfoundry-cli.list
sudo apt-get update
sudo apt-get install cf7-cli -y
cf -v
