#!/bin/bash
go get
go build
mv gojson.com /tmp/
killall gojson.com
declare -x PORT="7897"
su www-data sh -c "/tmp/gojson.com"