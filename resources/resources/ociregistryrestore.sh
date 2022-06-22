#! /bin/bash
set -o errexit

function abspath() {
    cd "$@" && pwd 
}

if [ -z ${reg_name} ]; then
    reg_name='kr'
fi
if [ -z ${reg_version} ]; then
    reg_version="1.24.0"
fi
if [ -z ${reg_path} ]; then
    reg_path="$HOME/.config/kind/registry"
fi

[ -z $restore_dir ] && echo "A restore directory must be provided for this operation." >&2 && exit 1

# cleanup any existing volumes of the same name
if [ "$(docker ps -a --filter volume="${reg_name}" 2>/dev/null || true)" != 'true' ]; then
    docker volume rm $reg_name || true
    docker volume rm ${reg_name}cfg || true
    rm -rf $reg_path
    rm -rf ${reg_path}../config
    
fi

mkdir -p $reg_path ${reg_path}/../config
config_mount=$( abspath "${reg_path}/../config" )

config_vname=${reg_name}cfg

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

docker run --rm -v $reg_name:/restore -v $restore_dir:/staging alpine:3 sh -c "cd /restore && tar -zxvf /staging/volumes.tar.gz && mv var/lib/registry/* . 2>/dev/null || rm -rf var && rm -rf root"
docker run --rm -v $config_vname:/restore -v $restore_dir:/staging alpine:3 sh -c "cd /restore && tar -zxvf /staging/volumes.tar.gz && mv root/* . 2>/dev/null || rm -rf root && rm -rf var"
