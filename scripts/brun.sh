#!/bin/bash

SUPPORTED_PLATFORMS=("windows/386" "windows/amd64" "linux/386" "linux/amd64")

cd ..

usage() {
    echo "Usage: $0 <type> <platform/architecture>"
    [[ -n "$1" ]] && echo "$1"
    exit 2
}

check_error() {
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
}

# Check if there are less than 2 arguments.
[[ $# -lt 2 ]] && usage "Insufficient arguments (expected 2, got $#)"

# Extract OS and architecture from input.
IFS="/" read -r input_os input_arch <<< "$2"

is_supported=false
for platform in "${SUPPORTED_PLATFORMS[@]}"; do
    [[ "$platform" == "$input_os/$input_arch" ]] && is_supported=true && break
done

if ! $is_supported; then
    usage "This platform/architecture combo is currently unsupported."
fi

output_name=$1'-'$input_os'-'$input_arch
[[ "$input_os" == "windows" ]] && output_name+=".exe"

if [ $1 = "install" ]; then
    env GOOS=$input_os GOARCH=$input_arch go build -o build/$output_name internal/installer/install.go
else
    env GOOS=$input_os GOARCH=$input_arch go build -o build/$input_os/$output_name cmd/$1/main.go
fi
check_error

[[ $3 != "r" ]] && exit 0
[[ "$input_os" == "windows" ]] && echo "Package built successfully, but can't be run." && exit 0 
cd build
./$output_name
check_error

# Architectures
# aix/ppc64
# android/386
# android/amd64
# android/arm
# android/arm64
# darwin/amd64
# darwin/arm64
# dragonfly/amd64
# freebsd/386
# freebsd/amd64
# freebsd/arm
# freebsd/arm64
# illumos/amd64
# ios/amd64
# ios/arm64
# js/wasm
# linux/386
# linux/amd64
# linux/arm
# linux/arm64
# linux/mips
# linux/mips64
# linux/mips64le
# linux/mipsle
# linux/ppc64
# linux/ppc64le
# linux/riscv64
# linux/s390x
# netbsd/386
# netbsd/amd64
# netbsd/arm
# netbsd/arm64
# openbsd/386
# openbsd/amd64
# openbsd/arm
# openbsd/arm64
# openbsd/mips64
# plan9/386
# plan9/amd64
# plan9/arm
# solaris/amd64
# windows/386
# windows/amd64
# windows/arm