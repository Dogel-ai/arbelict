#!/bin/bash

build() {
    ./brun.sh $1 linux/386
    ./brun.sh $1 linux/amd64
    ./brun.sh $1 windows/386
    ./brun.sh $1 windows/amd64
}

rm ../build/install*
rm -r ../build/windows/
rm -r ../build/linux/
rm ../build/arbelict.zip

build install
build cli
build gui

./archive.sh