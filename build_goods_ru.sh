#!/bin/sh

GOOS=windows GOARCH=amd64 go build -o bin/eport_goods_win.exe export_yml_goods.go
GOOS=linux go build -o bin/export_goods_linux export_yml_goods.go
GOOS=darwin go build -o bin/export_goods_mac export_yml_goods.go