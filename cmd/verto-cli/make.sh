#!/usr/bin/env bash
COMMIT=$(git rev-list --abbrev-commit -1 HEAD)
CGOENABLED=0 go build -ldflags "-w -s -X main.Version=0.0.1-$COMMIT -X main.AppName=verto-cli"
