#!/usr/bin/env sh

if [ "$(id -u)" -ne 0 ]; then
  echo "please run as root"
  exit
fi

pin=$1
gpioset 0 $pin
