#!/bin/bash
set -o errexit

# set parameter values
if [ -z ${reg_name} ]; then
  reg_name='kr'
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

config_mount=
config_vname=
owner=

if [[ ! $EUID -gt 0 ]]; then
  echo "Running as root. Correct file permissions"
  # Try our best to get the ownership correct when running as root
  # Assume the root of the bundle directory has the correct owner
  if [ ! -z "$(which stat)" ] && stat -c 2>/dev/null; then
    owner=$(stat -c '%U' $BUNDLE_DIR/..):$(stat -c '%U' $BUNDLE_DIR/..)
    # gnu stat is available
  elif
    [ ! -z "$(which stat)" ]
  then
  # darwin stat is available
    owner=$(stat -f '%u' $BUNDLE_DIR/..):$(stat -f '%g' $BUNDLE_DIR/..)
  else
  # stat is unavailable. Default to current user
    owner=$(whoami):$(whoami)
  fi
  # set the bundle directory owner to the best guess owner
  chown -R $owner $BUNDLE_DIR
  chmod -R 0760 $BUNDLE_DIR
fi

# create docker volume if it doesn't exist. If it does skip this part.
if [ "$(docker ps -a --filter volume="${reg_name}" 2>/dev/null || true)" != 'true' ]; then

  # make the directories to hold the volume data
  mkdir -p $reg_path ${reg_path}/../config
  config_mount=$(
    cd ${reg_path}/../config
    pwd
  )
  config_vname=${reg_name}cfg

  # correct file permissions
  if [ ! -z $owner ]; then
    chown -R $owner $config_mount &
    chmod -R 0766 $config_mount
  fi

  # create the docker volumes
  docker volume create --driver local \
    --opt type=none \
    --opt device=$reg_path \
    --opt o=bind \
    $reg_name

  docker volume create --driver local \
    --opt type=none \
    --opt device=$config_mount \
    --opt o=bind \
    $config_vname
fi

# Check if docker container is running
if [ "$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)" != 'true' ]; then
  # If it isn't running check if it exists at all
  if [ "$(docker container ls --all | grep "${reg_name}" 2>/dev/null || true)" != 'true' ]; then
    # If it doesn't exist create it
    docker run \
      -d --restart=always -p "127.0.0.1:${reg_port}:5000" --name "${reg_name}" \
      -v $reg_name:/var/lib/registry \
      -v $config_vname:/root/ \
      registry:2
  else
    # If it does exist and is just offline, start it
    docker start "$reg_name"
  fi
fi

# create a cluster with the local registry enabled in containerd
# TODO: Check if the cluster already exists and skip this part if it does
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
