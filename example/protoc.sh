#!/bin/bash

protoc \
--proto_path=. \
--proto_path=../../ \
--go_out=. \
--go_opt=paths=source_relative \
--golemporal_out=. \
--golemporal_opt=paths=source_relative \
*/*.proto