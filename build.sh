#!/bin/sh

if [ -f wiener ]
then
  rm wiener
fi

goimports -w .
go fmt
go build

if [ -f wiener ]
then
  ./wiener
fi
