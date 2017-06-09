#!/bin/bash

GOOS=windows GOARCH=amd64 go build -o bin/eport_yml_win.exe export_yml.go
GOOS=linux go build -o bin/export_yml_linux export_yml.go
GOOS=darwin go build -o bin/export_yml_mac export_yml.go
