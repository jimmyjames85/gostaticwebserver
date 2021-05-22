#!/bin/bash
GOOS=linux GOARCH=arm GOARM=7 go build -o a.out cmd/gostaticwebserver/main.go
