#!/bin/bash

set -euo pipefail



create_func_with_src_package() {


}

create_func_with_deploy() {


}

create_src_package() {
    pkg_name=$1
    pkg_namespace=$2

}

create_env() {
    env_name=$1
    env_namespace=$2
    runtime_image=$3
    builder_image=$4


}


main() {

    # pool mgr tests
    # 1. env ns1, func ns2 with code, route and verify specialized pod in ns1, also verify pkg in ns2

    # 2. env ns3, deploy pkg ns3, func in ns3, route and verify specialized pod in ns3

    # 3. env default, func with code in ns2, route and verify specialized pod in fission-function

    # 4. cleanup all envs, funcs n pkgs.

    # builder mgr tests
    # 1. env with builder image and runtime image in ns1, src pkg and func in ns2, route and verify specialized pod in ns1

    # 2. env with builder image and runtime image in default, src pkg and func in ns2,
    #    route and verify specialized pod in fission-builder

    # new deploy mgr tests
    # 1. env ns1, func ns2 with code, route and verify specialized pod in ns2,

    # 2. env default, func ns2 with code, route and verify specialized pod in fission-function

}

main