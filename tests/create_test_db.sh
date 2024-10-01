#!/usr/bin/env bash
set -xeuo pipefail

CONT_NAME=antibforce_test
VOL_NAME=pg_data_antibforce_test

docker container rm ${CONT_NAME} -f 2>/dev/null || true
docker volume remove ${VOL_NAME} 2>/dev/null || true

docker run -d --name ${CONT_NAME} \
	-e POSTGRES_USER=otus_user \
	-e POSTGRES_PASSWORD=dev_pass \
	-e POSTGRES_DB=antibforce \
	-e PGDATA=/var/lib/postgresql/data/pgdata \
	-v ${VOL_NAME}:/var/lib/postgresql/data \
	-p 5532:5432 \
	postgres:14
