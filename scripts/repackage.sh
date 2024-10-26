#!/bin/bash

build() {
    ./brun.sh $1 linux/386
    ./brun.sh $1 linux/amd64
    ./brun.sh $1 windows/386
    ./brun.sh $1 windows/amd64
}

rm ../build/install*
rm ../build/cli*
rm ../build/gui*
rm ../build/arbelict.zip

build install
build cli
build gui

./archive.sh