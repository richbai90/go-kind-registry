#!/bin/bash

# set parameter values
if [ -z ${reg_name} ]; then
  reg_name='kr'
fi

if [ -z ${debug} ]; then

    exec &>/dev/null
fi

container stop $reg_name
docker container rm $reg_name
docker volume rm $reg_name
docker volume rm ${reg_name}cfg
kind delete cluster --name $reg_name
