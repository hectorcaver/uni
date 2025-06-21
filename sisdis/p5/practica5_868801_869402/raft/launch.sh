#!/bin/bash
pkill srvraft
for i in 0 1 2 3
do
    ./bin/srvraft $i "127.0.0.1:29130" "127.0.0.1:29131" "127.0.0.1:29132" "127.0.0.1:29133" &
done
