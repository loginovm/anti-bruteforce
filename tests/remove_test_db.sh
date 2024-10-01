#!/usr/bin/env bash

CONT_NAME=antibforce_test
VOL_NAME=pg_data_antibforce_test

docker container rm ${CONT_NAME} -f 2>/dev/null || true
docker volume remove ${VOL_NAME} 2>/dev/null || true