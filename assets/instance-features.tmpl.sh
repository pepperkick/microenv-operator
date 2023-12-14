#!/bin/bash

set -xe

export PATH="$PATH:/host/bin"

function setupFeatures() {
  cd /home/ec2-user/menv-cluster

  export KUBECONFIG=./kubeconfig.external

  if [[ "$SCRIPT_USAGE_INSTALL_ARGO_WORKFLOW" == "true" ]]; then
    kubectl create namespace argo || true
    kubectl apply -n argo -f https://github.com/argoproj/argo-workflows/releases/download/v3.4.13/install.yaml
    if [[ ! -z "$SCRIPT_USAGE_ARGO_WORKFLOW" ]]; then
      echo "$SCRIPT_USAGE_ARGO_WORKFLOW" | base64 -d | gunzip | kubectl apply -f -
    fi
  fi
}

setupFeatures
