#!/usr/bin/env bash
COMMIT=$(git rev-list --abbrev-commit -1 HEAD)
go build -ldflags "-X main.Version=0.0.1-$COMMIT -X main.AppName=verto-cli"
