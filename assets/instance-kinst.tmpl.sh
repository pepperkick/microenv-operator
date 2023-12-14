#!/bin/bash

set -xe

apt update
apt install -y curl zip git

if ! which aws; then
  curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
  unzip awscliv2.zip
  ./aws/install
fi
aws --version

if ! which jq; then
  curl -Lo ./jq https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64
  chmod +x ./jq
  mv ./jq /usr/bin/jq
fi
jq --version

export PATH="$PATH:/host/bin"

function setupKinst() {
  mkdir -p /home/ec2-user/menv-cluster
  cd /home/ec2-user/menv-cluster

  curl -Lo ./build.sh https://raw.githubusercontent.com/pepperkick/microenv/main/build.sh
  chmod +x ./build.sh

   cat <<EOF > ./build.yaml
output: menv
modules:
  - git: https://github.com/pepperkick/microenv.git
    paths:
    - modules/menv-assets-cluster
    - modules/menv-core-ingress
    - modules/menv-cluster-kinst
    - modules/menv-letsencrypt
    - modules/menv-aws
EOF

  ./build.sh -c "./build.yaml"

  unzip -o menv.zip

  mkdir -p certs
  if [[ "$SCRIPT_CERT_ISSUER" == "cert-manager" ]]; then
    if [[ ! -z "$SCRIPT_CERT_CRT" ]]; then
      echo "$SCRIPT_CERT_CRT" > ./certs/cert-bundle.pem
    fi
    if [[ ! -z "$SCRIPT_CERT_KEY" ]]; then
      echo "$SCRIPT_CERT_KEY" > ./certs/key.pem
    fi
  fi

  chmod +x ./menv.sh

   cat <<EOF > ./config.yaml
{{ .microenvConfig }}
EOF

  ./menv.sh create --after-system-setup
}

setupKinst
