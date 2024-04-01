#!/bin/bash

read -p "Please input the version of the image: " version

make docker-build docker-push IMG=registry.cn-chengdu.aliyuncs.com/greatsql/greatsql-controller:$version

make deploy IMG=registry.cn-chengdu.aliyuncs.com/greatsql/greatsql-controller:$version