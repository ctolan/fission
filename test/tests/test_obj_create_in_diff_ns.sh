#!/bin/bash

#test:disabled

set -euo pipefail


create_func_with_src_package() {
    func_name=$1
    func_ns=$2
}

create_func_with_deploy() {
    func_name=$1
    func_ns=$2
    # fission function create --name hello-newdeploy --env python --code ~/Documents/fission-artifacts/hello.py --minscale 1 --maxscale 4 --executortype newdeploy --fns ns2 --envns ns1

}

create_src_package() {
    pkg_name=$1
    pkg_ns=$2
    zip_file=$3

}

create_env() {
    env_name=$1
    env_ns=$2
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


    # internal route tests.
    # 1. env ns1, func ns2 with code, curl http://FISSION_ROUTER/fission-function/ns2/func -> should work

    # 2. env in default, func ns2 with code, curl http://FISSION_ROUTER/fission-function/func -> should work



    # timer trigger tests.
    # 1. env ns1, func ns2 with code, tt ( with a one time cron string for executing imm'ly), verify function is executed
    # this also indirectly tests internal route established at http://FISSION_ROUTER/ns1/func

    # 2. env in default, func ns2 with code, tt ( with a one time cron string for executing imm'ly), verify function is executed
    # this also indirectly tests internal route established at http://FISSION_ROUTER/func



    # kube watch tests.
    # 1. env ns1, func ns2 with code, watch Trigger ( TBD), verify function is executed
    # this also indirectly tests internal route established at http://FISSION_ROUTER/ns1/func

    # 2. env in default, func ns2 with code, watch Trigger ( TBD ), verify function is executed
    # this also indirectly tests internal route established at http://FISSION_ROUTER/func



    # msq trigger tests.
    # 1. integrate after mqtrigger is checked into master.

    log "WIP"

}

main