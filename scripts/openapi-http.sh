#!/bin/bash
set -e

readonly service="$1"
readonly output_dir="$2"
readonly package="$3"

oapi-codegen -generate types -o "$output_dir/openapi_types.gen.go" -package "handlers" "api/openapi/$service.yml"
oapi-codegen -generate gin -o "$output_dir/openapi_api.gen.go" -package "handlers" "api/openapi/$service.yml"
oapi-codegen -generate types -o "libs/api/http/$service/openapi_types.gen.go" -package "$service" "api/openapi/$service.yml"
oapi-codegen -generate client -o "libs/api/http/$service/openapi_client_gen.go" -package "$service" "api/openapi/$service.yml"