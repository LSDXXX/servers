#!/bin/bash
set -e

readonly service="$1"

protoc \
  --proto_path=api/pb "api/pb/$service.proto" \
  "--go_out=plugins=grpc:libs/api/grpc/$service" --go_opt=paths=source_relative 
