#!/bin/bash

go run github.com/webrpc/webrpc/cmd/webrpc-gen -schema=api.ridl -target=go -pkg=main -server -out=./server-go/api.gen.go
go run github.com/webrpc/webrpc/cmd/webrpc-gen -schema=api.ridl -target=go -pkg=main -client -out=./client-go/chat_client.gen.go
go run github.com/webrpc/webrpc/cmd/webrpc-gen -schema=api.ridl -target=ts -client -out=./client-ts/client.gen.ts