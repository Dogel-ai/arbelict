#!/bin/bash
mkdir ../temp
mv ../build/config.yaml ../temp/config.yaml

rm -r ../build/*
mkdir ../build/temp

mv ../temp/config.yaml ../build/config.yaml
rm -r ../temp