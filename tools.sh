#!/usr/bin/env bash

set -o errexit
set -o nounset
# set -o xtrace

if set +o | grep -F 'set +o pipefail' > /dev/null; then
  # shellcheck disable=SC3040
  set -o pipefail
fi

if set +o | grep -F 'set +o posix' > /dev/null; then
  # shellcheck disable=SC3040
  set -o posix
fi

# -----------------------------------------------------------------------------

__usage="
Usage: $(basename $0) [OPTIONS]

Options:
  node              ...
  php               ...
"

# -----------------------------------------------------------------------------

function func_node_cli() {
  docker run -it --rm -p "6060:6060" -v $(pwd)/example/node:/workdir -w /workdir \
    --entrypoint=ash node:16.16.0-alpine3.15
}

function func_php_cli() {
  CON=$(docker build -q -f deployment/container/php/Dockerfile deployment/container/php/Dockerfile) 
  docker run -it --rm -p "5050:5050" -v $(pwd)/example/php:/workdir -w /workdir \
    --entrypoint=ash ${CON}
}

function func_redis_server() {
  docker run --name redis-server -d --rm -p "6379:6379" redis:7.0.4-alpine
}

function func_redis_cli() {
  docker run -it --rm --name redis-cli --entrypoint=ash redis:7.0.4-alpine
}

# -----------------------------------------------------------------------------

if [ -z "$*" ]
then
  echo "$__usage"
else
    if [ $1 == "--help" ] || [ $1 == "-h" ]
    then
        echo "$__usage"
    fi

    if [ $1 == "node" ]
    then
      func_node_cli
    fi

    if [ $1 == "php" ]
    then
      func_php_cli
    fi

    if [ $1 == "redis-server" ]
    then
      func_redis_server
    fi

    if [ $1 == "redis-cli" ]
    then
      func_redis_cli
    fi
fi
