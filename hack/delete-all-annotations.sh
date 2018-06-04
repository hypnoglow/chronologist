#!/bin/bash

# Copyright 2018 The Chronologist Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script deletes all annotations created by Chronologist.

set -euo pipefail

GRAFANA_ADDR=${GRAFANA_ADDR:-}
if [ -z "${GRAFANA_ADDR}" ]; then
    echo "GRAFANA_ADDR is empty"
    exit 1
fi

GRAFANA_API_KEY=${GRAFANA_API_KEY:-}
if [ -z "${GRAFANA_API_KEY}" ]; then
    echo "GRAFANA_API_KEY is empty"
    exit 1
fi

while true
do
    declare -a ids
    ids=($(curl -sS -XGET "${GRAFANA_ADDR}/api/annotations?tags=owner%3Dchronologist" \
        -H "Authorization: Bearer ${GRAFANA_API_KEY}" \
        | jq ".[].id"))
    [ ${#ids[@]} -eq 0 ] && break

    for id in "${ids[@]}"; do
        echo "Deleting $id"
        curl -sS -XDELETE "${GRAFANA_ADDR}/api/annotations/${id}" \
            -H "Authorization: Bearer ${GRAFANA_API_KEY}"
        echo
    done
done

