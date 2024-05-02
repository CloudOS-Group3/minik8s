#!/bin/bash

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

cd $PROJECT_ROOT

programs=(
    "./cmd/apiserver/apiserver.go"
    "./cmd/shechuler/shceduler/go"
)

for program in "${programs[@]}"; do
    sudo go run "$program"
done