#!/bin/bash
# last login to linux
# chmod +x last_login.sh 

last_login_info=$(last -1)

user=$(echo "$last_login_info" | awk '{print $1}')
ip=$(echo "$last_login_info" | awk '{print $3}')
datetime=$(echo "$last_login_info" | awk '{print $4, $5, $6, $7}')

echo "User: $user"
echo "IP: $ip"
echo "Date time: $datetime"
