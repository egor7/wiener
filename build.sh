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

feh s.png 

### im_profile -v magnitude_cone.png magnitude_cone_profile.png
### feh magnitude_cone_profile.png 
### #feh magnitude_cone.png 


### if [ -f wiener ]
### then
###   rm wiener
### fi
### 
### goimports -w .
### go fmt
### go build
### 
### if [ -f wiener ]
### then
###   ./wiener
### fi
### 
### im_profile -v magnitude_cone.png magnitude_cone_profile.png
### feh magnitude_cone_profile.png 
### #feh magnitude_cone.png 


#convert 5.3.02.png -gaussian-blur 0x2 5.3.02.gaus.png
#convert -size 512x512 xc:black -fill white \
#        -draw "point 256,256" -gaussian-blur 0x2 -auto-level \
#        -write gaus2.png -fft -delete 1 \
#        -auto-level -evaluate log 10000 gaus2_spectrum.png
#
#im_profile gaus2.png gaus2_pf.gif
#im_profile gaus2_spectrum.png gaus2_spectrum_pf.gif
#feh gaus2_spectrum_pf.gif
