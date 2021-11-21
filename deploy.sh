#!/usr/bin/env bash

set -o errexit

docker-compose up --detach --build

echo "
    Web UI: http://localhost:9000
    SSH to node-abc: ssh alpha@localhost:9001 (password: alpha)
    SSH to node-xyz: ssh alpha@localhost:9002 (password: alpha)
"
