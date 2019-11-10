#!/usr/bin/env bash

set -e

dir=$(cd `dirname $0` && cd .. && pwd)
export GOARCH=amd64

if [[ -z "${1}" ]]; then
  echo "USAGE: $0 <release-version>"
  exit 1
fi

version=$1

GOOS=linux go build \
  -ldflags "-X github.com/aemengo/concourse-worker-manager/actions.version=${version}" \
  -o ${dir}/cwm-linux \
  github.com/aemengo/concourse-worker-manager

GOOS=darwin go build \
  -ldflags "-X github.com/aemengo/concourse-worker-manager/actions.version=${version}" \
  -o ${dir}/cwm-darwin \
  github.com/aemengo/concourse-worker-manager

GOOS=windows go build \
  -ldflags "-X github.com/aemengo/concourse-worker-manager/actions.version=${version}" \
  -o ${dir}/cwm-windows.exe \
  github.com/aemengo/concourse-worker-manager

cd ${dir}

tar czf cwm-linux-v${version}.tgz cwm-linux 2>/dev/null
tar czf cwm-darwin-v${version}.tgz cwm-darwin 2>/dev/null
zip cwm-windows-v${version}.zip cwm-windows.exe 1>/dev/null

rm -f ${dir}/cwm-linux ${dir}/cwm-darwin ${dir}/cwm-windows.exe