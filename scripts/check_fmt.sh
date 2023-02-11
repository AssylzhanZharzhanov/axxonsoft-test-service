#!/usr/bin/env bash
res=$(gofmt -e -l -d .)
if [[ -n "$res" ]]; then
    echo "Go fmt failed"
    echo "$res"
    exit 1
fi
