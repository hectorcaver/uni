#!/bin/bash

# This script is in charge of executing servers according to the flags received

# Usage example:
# ./server_launcher.sh -t 127.0.0.1:4000
# ./server_launcher.sh -u 192.168.1.10:5000

# Check if the number of arguments is valid
if [ "$#" -lt 2 ]; then
    echo "Usage: $0 -t|-uX ip:port"
    exit 1
fi

# Extract flags and endpoint
flag=$1
endpoint=$2

# Check synchronization number for UDP
if [[ "$flag" =~ ^-u([0-9]+)$ ]]; then
    sync_num=${BASH_REMATCH[1]}
    flag="-u"
else
    flag="-t"
    sync_num=0
fi

# Execute based on flag
case "$flag" in
  -t)
    sleep 2

    go run network/server_tcp/tcp_server.go $endpoint &

    pid=$!

    sleep 5

    kill $pid
    pkill "tcp_server"
    ;;

  -u)
    go run network/server_udp/server_udp.go $endpoint &

    pid=$!

    sleep $(( sync_num * 2 ))

    kill $pid
    pkill "server_udp"
    ;;

  *)
    echo "Invalid flag. Use -t for TCP or -u for UDP."
    exit 1
    ;;
esac