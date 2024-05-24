#!/bin/bash

cd /root/minik8s

mkdir -p log

programs=(
    "./cmd/apiserver/apiserver.go:./log/apiserver.log"
    "./cmd/kubelet/kubelet.go:./log/kubelet.log"
    "./cmd/controller/controller.go:./log/controller.log"
)

for program in "${programs[@]}"; do

    IFS=':' read -ra ADDR <<< "$program"
    program_file="${ADDR[0]}"
    log_file="${ADDR[1]}"

    echo "go文件是$program_file," "log文件是$log_file"

    touch "$log_file"

    echo "$PATH"

    sudo go run "$program_file"

done

wait
