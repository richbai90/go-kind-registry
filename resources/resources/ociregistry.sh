#!/bin/sh
set -o errexit

# set parameter values
if [ -z ${reg_name} ]; then
  reg_name='kind-registry'
fi
if [ -z ${reg_port} ]; then
  reg_port='5015'
fi
if [ -z ${reg_path} ]; then
  reg_path="$HOME/.config/kind/registry"
fi
if [ -z ${reg_version} ]; then
  reg_version="1.24.0"
fi

# create docker volume
if [ "$(docker ps -a --filter volume="${reg_name}" 2>/dev/null || true)" != 'true' ]; then

  mkdir -p $reg_path $reg_path &
  docker volume create --driver local \
      --opt type=none \
      --opt device=$reg_path \
      --opt o=bind \
      $reg_name
fi

if [ "$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)" != 'true' ]; then
  if [ "$(docker container ls --all | grep "${reg_name}" 2>/dev/null || true)" != 'true' ]; then
    docker run \
      -d --restart=always -p "127.0.0.1:${reg_port}:5000" --name "${reg_name}" \
      -v $reg_name:/var/lib/registry \
      registry:2
  else
    docker start "$reg_name"
  fi
fi



# create a cluster with the local registry enabled in containerd
cat <<EOF | kind create cluster --name $reg_name --image kindest/node:v$reg_version --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:${reg_port}"]
    endpoint = ["http://${reg_name}:5000"]
EOF

# connect the registry to the cluster network if not already connected
if [ "$(docker inspect -f='{{json .NetworkSettings.Networks.kind}}' "${reg_name}")" = 'null' ]; then
  docker network connect "kind" "${reg_name}"
fi

# Document the local registry
# https://github.com/kubernetes/enhancements/tree/master/keps/sig-cluster-lifecycle/generic/1755-communicating-a-local-registry
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:${reg_port}"
    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF



