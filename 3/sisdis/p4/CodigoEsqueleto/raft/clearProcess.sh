#!/bin/bash

for x in 1 2 3; do
    pid=$(lsof -t -i :2900$x)
    if [ -n "$pid" ]; then
        kill $pid  # <- sin comillas aquí
        echo "Killed process $pid on port 2900$x"
    else
        echo "No process found on port 2900$x"
    fi
done
