#!/bin/bash

set -e
exec > >(tee /var/log/user-data.log|logger -t user-data -s 2>/dev/console) 2>&1

function installDocker() {
  mkdir -p /home/ec2-user/menv
  cd /home/ec2-user/menv
  curl -Lo menv.zip https://github.com/pepperkick/microenv/releases/download/nightly/menv-setup-aws.zip
  unzip -o menv.zip
  chmod +x ./menv.sh

   cat <<EOF > ./config.yaml
{{ .microenvConfig }}
EOF

   echo "{{ .customInstanceSetupScript }}" | base64 -d > ./scripts/features/custom-script.sh

  ./menv.sh create --only-system-setup

  systemctl enable docker
}

# Enable cloud-init run on every boot
sed -i 's/scripts-user$/\[scripts-user, always\]/' /etc/cloud/cloud.cfg

installDocker
systemctl start docker

if ! which docker; then
  echo "ERROR: Failed to install docker"
  exit 1
fi

# Ensure docker is running
if ! docker ps; then
  systemctl start docker
fi

if ! docker ps; then
  echo "ERROR: Failed to start docker"
  exit 1
fi
