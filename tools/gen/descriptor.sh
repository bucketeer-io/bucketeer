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
api_gateway_values_path="./manifests/bucketeer/charts/api/values.yaml"
encoded_descriptor=$(cat ${DESCRIPTOR_PATH}/gateway/${descriptor_file} | base64 | tr -d \\n)
echo $encoded_descriptor > encoded_descriptor.txt
yq eval ".envoy.descriptor = load(\"encoded_descriptor.txt\")" -i ${api_gateway_values_path}
rm encoded_descriptor.txt

# web-gateway
web_gateway_values_path="./manifests/bucketeer/charts/web/values.yaml"
proto_descriptor_dirnames=$(find ${DESCRIPTOR_PATH} -name "$descriptor_file" -not -path "**/gateway/*" -print0 | xargs -0 -n1 dirname | awk -F/ '{print $NF}')
for service_name in $proto_descriptor_dirnames
do
  encoded_descriptor=$(cat ${DESCRIPTOR_PATH}/${service_name}/${descriptor_file} | base64 | tr -d \\n)
  yq eval ".envoy.${service_name}Descriptor = \"${encoded_descriptor}\"" -i ${web_gateway_values_path}
done
