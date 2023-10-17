#!/bin/bash

umount ./mount
rm -rf upper work mount
cp -r upper_base upper
mkdir work mount
mount -t overlay overlay -o lowerdir=./lower,upperdir=./upper,workdir=./work ./mount
