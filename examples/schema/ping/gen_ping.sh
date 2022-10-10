#!/bin/bash 
SCHEMA_PATH="./examples/schema"
# SCHEMA_INCLUDE=$SCHEMA_PATH
SCHEMA_INCLUDE="./examples/schema"


OUT_PATH="./examples/types"
SWAGGER_PATH="./examples/swagger/service"

export GOPATH="$(go env GOPATH)"
export PATH="$PATH:$(go env GOPATH)/bin"

mkdir -p $OUT_PATH
mkdir -p $SWAGGER_PATH 

source "./schema/scripts/gen_func.sh"

echo "1. Generate service: ping.proto"
# gen_aoi_gateway proto/service.proto
gen_service ping/ping.proto
gen_gateway ping/ping.proto ping/ping_service.yaml
gen_swagger ping/ping.proto ping/ping_service.yaml
# gen_openapi stamp/stamp.proto stamp/stamp_service.yaml