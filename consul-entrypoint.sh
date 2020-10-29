#!/bin/bash

echo start mongodb
docker-entrypoint.sh "$@" &

echo start health check
consul-health-check