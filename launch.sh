#!/bin/sh
while :
do
./spuri.io
[[ $? != 0 ]] || break
done
