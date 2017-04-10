#!/bin/bash
keySize=8
valSize=8
totalsize=1024
fileName=/data/appdatas/easyKV/persistentmap
cd ./cmd/server/main.go && go run main.go $keySize $valSize $totalsize $fileName
cd ./cmd/client/main.go && go run main.go
