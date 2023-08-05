#!/bin/bash
if [ $# -le 0 ]; then
    echo "Usage: $0 <pcap file>"
    exit 1
fi
strings $1 | grep -i -E "select|insert|update|delete|create|drop|alter|truncate|vacuum|begin|commit"
