#!/bin/sh

echo Downloading dependencies
go get -t ./...

echo Building application
cd cmd/standup && GOOS=linux GOARCH=amd64 go build -race -o ~/slack-standup || exit

echo Copying config
cp /go/src/github.com/vehsamrak/slack-standup/configs/config-dev.yml ~/config.yml

echo Running application
cd ~ && ./slack-standup
