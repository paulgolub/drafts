#!/bin/bash
# last entrance to linux

last_login=$(last -1 -R | awk '{print $4, $5, $6, $7}')
echo "Last entrance: $last_login"
