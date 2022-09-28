#!/bin/bash

set -eu

function usage() {
    cat <<_EOT_
Usage:
  $0 Args

Description:
  script for creating PR.

Environment Variables:
  - DIR
  - DESCRIPTOR_PATH
_EOT_
    exit 1
}

# validations
[[ -z $DIR ]] && usage
[[ -z $DESCRIPTOR_PATH ]] && usage

cd $DIR

make proto-go

descriptor_file="proto_descriptor.pb"

# api-gateway
api_gateway_values_path="./manifests/bucketeer/charts/api-gateway/values.yaml"
encoded_descriptor=$(cat ${DESCRIPTOR_PATH}/gateway/${descriptor_file} | base64 | tr -d \\n | sed -E "s|\/|\\\/|")
sed -i -E "s|(descriptor): .+|\1: \"${encoded_descriptor}\"|" ${api_gateway_values_path}

# web-gateway
web_gateway_values_path="./manifests/bucketeer/charts/web-gateway/values.yaml"
proto_descriptor_dirnames=$(find ${DESCRIPTOR_PATH} -name "$descriptor_file" -not -path "**/gateway/*" -print0 | xargs -0 -n1 dirname | awk -F/ '{print $NF}')
for service_name in $proto_descriptor_dirnames
do
  encoded_descriptor=$(cat ${DESCRIPTOR_PATH}/${service_name}/${descriptor_file} | base64 | tr -d \\n | sed -E "s|\/|\\\/|")
	sed -i -E "s|(${service_name}Descriptor): .+|\1: \"${encoded_descriptor}\"|" ${web_gateway_values_path}
done