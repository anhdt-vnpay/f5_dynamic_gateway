#1. Generate type
gen_type() {
    protoc \
        -I ${SCHEMA_PATH} \
        -I $GOPATH/src \
        -I ${SCHEMA_INCLUDE} \
        -I $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
        --go_out=plugins=grpc:${OUT_PATH} \
        --go_opt paths=source_relative \
        ${SCHEMA_PATH}/$1
}
#2. Generate service
gen_service() {
    protoc  \
        -I ${SCHEMA_PATH} \
        -I $GOPATH/src \
        -I ${SCHEMA_INCLUDE} \
        -I $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
        --go_out=${OUT_PATH} \
        --go_opt paths=source_relative \
        --go-grpc_out ${OUT_PATH} \
        --go-grpc_opt paths=source_relative \
        ${SCHEMA_PATH}/$1
}
#3. Generate gateway
gen_gateway() {
    protoc  \
        -I ${SCHEMA_PATH} \
        -I ${GOPATH}/src \
        -I ${SCHEMA_INCLUDE} \
        -I $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.1.0/third_party/googleapis \
        --grpc-gateway_out ${OUT_PATH}  \
        --grpc-gateway_opt allow_delete_body=true \
        --grpc-gateway_opt logtostderr=true \
        --grpc-gateway_opt paths=source_relative \
        --grpc-gateway_opt grpc_api_configuration=${SCHEMA_PATH}/$2 \
        ${SCHEMA_PATH}/$1
}
#3. Generate swagger
gen_swagger() {
    protoc \
        -I ${SCHEMA_PATH} \
        -I ${GOPATH}/src \
        -I ${SCHEMA_INCLUDE} \
        -I $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway\@v1.16.0/third_party/googleapis \
        --swagger_out ${SWAGGER_PATH} \
        --swagger_opt allow_delete_body=true \
        --swagger_opt grpc_api_configuration=${SCHEMA_PATH}/$2 \
        ${SCHEMA_PATH}/$1

}
#4. Generate openapi
gen_openapi() {
    protoc \
        -I ${SCHEMA_PATH} \
        -I ${GOPATH}/src \
        -I ${SCHEMA_INCLUDE} \
        -I $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway\@v1.16.0/third_party/googleapis \
        --openapiv2_out ${SWAGGER_PATH} \
        --openapiv2_opt  allow_delete_body=true \
        --openapiv2_opt grpc_api_configuration=${SCHEMA_PATH}/$2 \
        ${SCHEMA_PATH}/$1
}

#5. Generate db 
gen_gorm() {
    protoc \
        -I  ${SCHEMA_PATH} \
        -I ${GOPATH}/src \
        -I ${SCHEMA_INCLUDE} \
        -I ${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
        --gorm_out ${OUT_PATH} \
        $1
}

#6. New Generate function
gen_aoi_gateway() {
    protoc  \
    -I ${SCHEMA_PATH} \
    -I ${GOPATH}/src \
    -I ${SCHEMA_INCLUDE} \
    -I $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.1.0/third_party/googleapis \
    --grpc-gateway_out ${OUT_PATH}  \
    --grpc-gateway_opt allow_delete_body=true \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    ${SCHEMA_PATH}/$1
}