#!/bin/bash 
SCHEMA_PATH="./schema"
# SCHEMA_INCLUDE=$SCHEMA_PATH
SCHEMA_INCLUDE="./schema"


OUT_PATH="./types"
SWAGGER_PATH="./swagger/service"

export GOPATH="$(go env GOPATH)"
export PATH="$PATH:$(go env GOPATH)/bin"

mkdir -p $OUT_PATH
mkdir -p $SWAGGER_PATH 

source "./schema/scripts/gen_func.sh"

echo "1. Generate service: api_registration.proto"
# gen_aoi_gateway proto/service.proto
gen_service registration/api_registration.proto
gen_gateway registration/api_registration.proto registration/api_registration_service.yaml
gen_swagger registration/api_registration.proto registration/api_registration_service.yaml
# gen_openapi stamp/stamp.proto stamp/stamp_service.yaml