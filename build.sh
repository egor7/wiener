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

im_profile magnitude_cone.png magnitude_cone_profile.png
feh magnitude_cone_profile.png
