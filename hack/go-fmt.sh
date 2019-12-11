#!/usr/bin/env bash

go fmt ./...

if [[ -n ${CI} ]]; then
    git diff --exit-code
fi
