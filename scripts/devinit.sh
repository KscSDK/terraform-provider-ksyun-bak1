#!/bin/bash

set -e

# Detech current os category
unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     OS_TYPE=Linux;;
    Darwin*)    OS_TYPE=Mac;;
    CYGWIN*)    OS_TYPE=Windows;;
    MINGW*)     OS_TYPE=Windows;;
    *)          OS_TYPE="UNKNOWN:${unameOut}"
esac

echo "OS ${OS_TYPE} is deteched."
echo "Compiling ..."

# Choice file path/name by os category
if [ $OS_TYPE == "Linux" ]; then
	GOOS=linux GOARCH=amd64 go build -o bin/terraform-provider-ksyun
	chmod +x bin/terraform-provider-ksyun
    mkdir -p $HOME/.terraform.d/plugins
    mv bin/terraform-provider-ksyun $HOME/.terraform.d/plugins
elif [ $OS_TYPE == "Mac" ]; then
	GOOS=darwin GOARCH=amd64 go build -o bin/terraform-provider-ksyun
	chmod +x bin/terraform-provider-ksyun
    mkdir -p $HOME/.terraform.d/plugins
    mv bin/terraform-provider-ksyun $HOME/.terraform.d/plugins
elif [ $OS_TYPE == "Windows" ]; then
	GOOS=windows GOARCH=amd64 go build -o bin/terraform-provider-ksyun.exe
	chmod +x bin/terraform-provider-ksyun.exe
    mkdir -p $APPDATA/terraform.d/plugins
    mv bin/terraform-provider-ksyun.exe $APPDATA/terraform.d/plugins
else
    echo "Invalid OS"
    exit 1
fi

echo "Installation of ksyun Terraform Provider is completed."
